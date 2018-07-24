package client

import (
	"net"
	"fmt"
	"bufio"
	"os"
	"flag"
	"strconv"
)

//global variable declaration
type peer struct{
	address string
	port int
	networkAddr string
	backupPeer string
	masterNode string
}

var peerNode peer


func Start()  {

	//Start the clientBuild
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

	//get the IP address of the system,
	//the clientBuild is running on
	ip, err := ExternalIP()
	if err != nil {
		fmt.Println(err)
	}

	peerNode = peer{networkAddr: ip, masterNode:address}

	//Connect to serverBuild
	establishConnection()

}

/*
Establish the connection to the serverBuild and
handle the incoming messages from the serverBuild
*/
func establishConnection() {

	//Dial the connection to the serverBuild
	conn, err :=net.Dial("tcp", peerNode.masterNode)
	if err != nil{
		fmt.Printf("Error establishing a connection\n")
		return
	}

	//handle the connection to the serverBuild
	sendMessage(conn)
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

	//Take the input from the clientBuild
	fmt.Printf("Enter the Message to send to the clientBuild: ")
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
