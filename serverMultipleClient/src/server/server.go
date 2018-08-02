package server

import (
	"flag"
	"net"
	"fmt"
	"strings"
	"strconv"
		"encoding/gob"
	"utils"
	"sync"
)

//global variable declaration
type master struct{
	address string
	port int
	networkAddr string
	peers map[string]int
	backupPeers map[string]string
}

var masterNode master
var enc *gob.Encoder
var dec *gob.Decoder
var mutex = &sync.Mutex{}

func main()  {

	//Start the serverBuild
	StartServer()
}


/*
Function which initializes the master struct with initial
values which are parsed from the command line.

Returns: nil
 */
func initializeMaster()  {

	//parse the command line arguments
	server := flag.String("serverAddr", "", "The address of the serverBuild to connect to." +
		"Default is localhost")

	port := flag.String("port", "9999", "Port to listen for incoming connections.")

	flag.Parse()

	//form the network address for the node
	address := *server+":"+*port

	//initialize the global variable
	//representing master node
	p, err := strconv.Atoi(*port)
	if err != nil{
		fmt.Printf("Conversion Error: %s", err.Error())
	}

	masterNode = master{address:*server, port:p, networkAddr:address,
						peers:make(map[string]int), backupPeers:make(map[string]string)}

}


/*
Function which listens on the dedicated port specified
for the master node and accepts the clients requests only
to pass it to a go routine to handle the requests

Returns: nil
 */
func acceptAndProcess(node master){

	//listen on the designates network address
	adapter, err := net.Listen("tcp", node.networkAddr)
	if err != nil{
		fmt.Printf("Error while listening to the on port: %d", node.port)
		return
	}

	//until a SIGNAL interrupt is passed or an exception is
	//raised, keep on accepting peerBuild connections and add it
	//to the peer map.
	for{

		//debug information
		fmt.Printf("\nListening on Port: %d\n", node.port)

		//accept incoming connections
		conn, err := adapter.Accept()
		if err != nil {
			println(err.Error())
			continue
		}

		// Will write to network.
		enc = gob.NewEncoder(conn)
		// Will read from network.
		dec = gob.NewDecoder(conn)


		//start a go routine to handle
		//the incoming connections
		go handleConnection(conn)
	}
}


/*
Function which handles the incoming
peerBuild requests to the serverBuild.
It performs any necessary action and/or invokes
other functions to complete the tasks

Returns: nil
 */
func handleConnection(conn net.Conn){

	//Receive and Decode the packet on the
	//network.
	var recv utils.Packet
	err := dec.Decode(&recv)
	if err!= nil{
		print("Error while decoding peer packet: ", err.Error())
	}

	//parse the packet
	if recv.Ptype == utils.PEER{

		//debug
		n := len(masterNode.peers)

		//get the address of the tcp-peerBuild
		clientAddr := conn.RemoteAddr().String()

		//add the peerBuild to the peer list
		networkAddr := strings.Split(clientAddr, ":")
		clientPort, err := strconv.Atoi(networkAddr[1])
		if err != nil{
			fmt.Printf("Conversion Error: %s", err.Error())
		}


		mutex.Lock()
		//masterNode.peers[networkAddr[0]] = clientPort
		masterNode.peers[clientAddr] = clientPort
		mutex.Unlock()

		//send response packet
		p := utils.Response{Ptype: utils.RESPONSE,Backup:false,NetAddress:""}
		err = enc.Encode(p)
		if err!= nil{
			print("Error while encoding peer packet: ", err.Error())
		}

		if n!= len(masterNode.peers){
			fmt.Println("Peer added")
		}
	}

	//close the connection
	conn.Close()
}

/*
Function which returns the address of master struct
for testing functions
*/
func StructAddr() *master {
	return &masterNode
}

/*
Function which starts the serverBuild and
passes the initialized values for listening
on the designated port.
*/
func StartServer()  {

	//initialize the structure to define the
	//master node
	initializeMaster()

	//once the master node is initialized,
	//listen on the dedicated port and accept
	//connections
	acceptAndProcess(masterNode)

}