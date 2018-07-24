package server

import (
	"flag"
	"net"
	"fmt"
	"strings"
	"strconv"
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
	//raised, keep on accepting clientBuild connections and add it
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

		//start a go routine to handle
		//the incoming connections
		go handleConnection(conn)
	}
}


/*
Function which handles the incoming
clientBuild requests to the serverBuild.
It performs any necessary action and/or invokes
other functions to complete the tasks

Returns: nil
 */
func handleConnection(conn net.Conn){
	//get the address of the tcp-clientBuild
	clientAddr := conn.RemoteAddr().String()

	//add the clientBuild to the peer list
	networkAddr := strings.Split(clientAddr, ":")
	clientPort, err := strconv.Atoi(networkAddr[1])
	if err != nil{
		fmt.Printf("COnversion Error: %s", err.Error())
	}

	//masterNode.peers[networkAddr[0]] = clientPort
	masterNode.peers[clientAddr] = clientPort

	//debug information
	fmt.Println("Connected Client: " + clientAddr)

	info := make([]byte, 255)
	n, err := conn.Read(info)

	fmt.Printf("size of text read %d\n", n)

	if err != nil{
		fmt.Println("Error Reading message from Client")
	} else {
		fmt.Println("Message Received from Client: " + string(info))
	}

	temp := strings.Split(string(info), " ")

	if temp[0] == "Hello"{
		conn.Write([]byte("Hello " + string(clientAddr)))
	}

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