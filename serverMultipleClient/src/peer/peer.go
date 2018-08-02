package peer

import (
	"net"
	"fmt"
	"bufio"
	"os"
	"flag"
	"strconv"
	"utils"
	"unsafe"
	"encoding/gob"
	"strings"
)

//global variable declaration
type peer struct{
	address string
	port int
	networkAddr string
	backupPeer string
	masterNode string
}

//Global variables
var peerNode peer
var backupNetAddr string
var enc *gob.Encoder
var dec *gob.Decoder

func Start()  {

	//Start the peerBuild
	initializePeer()

}

/*
Function which initializes the peer struct with initial
values which are parsed from the command line.

Returns: nil
 */
func initializePeer(){

	//parse the command line arguments
	remoteAddr := flag.String("addr", "127.0.0.1", "The address of the serverBuild to connect to." +
		"Default is localhost")

	remotePort := flag.String("port", "9999", "Port to listen for incoming connections.")

	flag.Parse()

	//form the network address for the node
	address := *remoteAddr+":"+*remotePort

	//initialize the global variable
	//representing master node
	_, err := strconv.Atoi(*remotePort)
	if err != nil{
		fmt.Printf("Conversion Error: %s", err.Error())
	}

	peerNode = peer{masterNode:address}

	//Connect to serverBuild
	//establishConnection(enc, dec)
	establishConnection()

	listenAndAccept()
	//listenAndAccept(enc, dec)
}

/*
Establish the connection to the serverBuild and
handle the incoming messages from the serverBuild
*/
func establishConnection() {

	//Dial the connection to the serverBuild
	//conn, err :=net.Dial("tcp", peerNode.masterNode)
	conn, err :=net.Dial("tcp", peerNode.masterNode)
	if err != nil{
		fmt.Printf("Error establishing a connection\n")
		return
	}

	//get the local address of the system
	peerNode.networkAddr = conn.LocalAddr().String()
	a := strings.Split(peerNode.networkAddr, ":")
	peerNode.address = a[0]
	peerNode.port,_ = strconv.Atoi(a[1])

	//create packet to send to the master
	p := utils.CreatePacket(utils.PEER,"", unsafe.Sizeof(utils.PEER))

	//send message to the master
	//initialize the encoder and decoder
	//to read the packets
	enc = gob.NewEncoder(conn) // Will write to network.
	dec = gob.NewDecoder(conn) // Will read from network.

	//Encode and send data over network
	err = enc.Encode(p)
	if err!= nil{
		print("Error while encoding peer packet: ", err.Error())
	}

	//Receive and decode data on the network
	var resp utils.Response
	err = dec.Decode(&resp)
	if err!= nil{
		print("Error while decoding peer packet: ", err.Error())
	}

	//process the response packet
	if resp.Ptype == utils.RESPONSE{
		if resp.Backup{
			backupNetAddr = resp.NetAddress
		}
		backupNetAddr = ""
	}

	conn.Close()
	//handle the connection to the serverBuild
	//sendMessage(conn)
}

/*
Function which is responsible for communication
with the serverBuild

Params: net.Conn
Returns: nil
 */
func sendMessage(conn net.Conn){

	var message string
	var msgSent []byte
	msgSent = make([]byte, 255)

	//Take the input from the peerBuild
	fmt.Printf("Enter the Message to send to the peerBuild: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	message = scanner.Text()

	//Send the message
	msgSent = []byte(message+"\n")
	n, err:= conn.Write(msgSent)
	if err != nil{
		fmt.Printf("Error while writing message to socket")
	}

	//Read the message
	var msg []byte
	msg = make([]byte, 255)
	n, err1 := conn.Read(msg)
	if err1 != nil{
		fmt.Printf("Error while reading mmessage from the remote %s", err1.Error())
	} else {
		fmt.Printf("text size is:  %d\n", n)
		fmt.Println("Message from Server: " + string(msg))
	}
}

/*
Function which listens to incoming requests from
clients for file access
*/
func listenAndAccept(){
	//listen on the designates network address
	adapter, err := net.Listen("tcp", peerNode.networkAddr)
	if err != nil{
		fmt.Printf("Error while listening to the on port: %d", peerNode.port)
		return
	}

	//until a SIGNAL interrupt is passed or an exception is
	//raised, keep on accepting peerBuild connections and add it
	//to the peer map.
	for{

		//debug information
		fmt.Printf("\nListening on Port: %d\n", peerNode.port)

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
		//go handleConnection(conn)
	}
}