package peer

import (
	"encoding/gob"
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"unsafe"
	"utils"
)

////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////
//global variable declaration
type peer struct {
	address       string
	port          int
	networkAddr   string
	myPrimaryPeer string
	backupPeer    string
	masterNode    string
}

//Global variables
var peerNode peer
var enc *gob.Encoder
var dec *gob.Decoder
var index map[string]bool

////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////

/*
Driver Program
*/
func Start() {

	//Start the peerBuild
	initializePeer()

}

/*
Function which initializes the  peer struct with initial
values which are parsed from the command line.

Returns: nil
*/
func initializePeer() {

	//parse the command line arguments
	remoteAddr := flag.String("addr", "127.0.0.1", "The address of the serverBuild to connect to."+
		"Default is localhost")

	remotePort := flag.String("port", "9999", "Port to listen for incoming connections.")

	flag.Parse()

	//form the network address for the node
	address := *remoteAddr + ":" + *remotePort

	//initialize the global variable
	//representing master node
	_, err := strconv.Atoi(*remotePort)
	if err != nil {
		fmt.Printf("Conversion Error: %s", err.Error())
	}

	peerNode = peer{masterNode: address}

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
	conn, err := net.Dial("tcp", peerNode.masterNode)
	if err != nil {
		fmt.Printf("Error establishing a connection\n")
		return
	}

	//get the local address of the system
	peerNode.networkAddr = conn.LocalAddr().String()
	a := strings.Split(peerNode.networkAddr, ":")
	peerNode.address = a[0]
	peerNode.port, _ = strconv.Atoi(a[1])

	//create packet to send to the master
	p := utils.CreatePacket(utils.PEER, "", unsafe.Sizeof(utils.PEER))

	//send message to the master
	//initialize the encoder and decoder
	//to read the packets
	enc = gob.NewEncoder(conn) // Will write to network.
	dec = gob.NewDecoder(conn) // Will read from network.

	//Encode and send data over network
	err = enc.Encode(p)
	if err != nil {
		print("Error while encoding peer packet: ", err.Error())
	}

	//Receive and decode data on the network
	var resp utils.Response
	err = dec.Decode(&resp)
	if err != nil {
		print("Error while decoding peer packet: ", err.Error())
	}

	//process the response packet
	if resp.Ptype == utils.RESPONSE {
		if resp.Backup {
			//update the primary peer
			//address to the current instance
			peerNode.myPrimaryPeer = resp.NetAddress

			//forward the same update to the
			//primary peer
			defer fmt.Println("Primary Peer Updated")
			go updatePrimary()
		} else {
			peerNode.myPrimaryPeer = ""
		}
	}

	conn.Close()
	//handle the connection to the serverBuild
	//sendMessage(conn)
}

/*
Function which listens to incoming requests from
clients for file access
*/
func listenAndAccept() {
	//listen on the designates network address
	adapter, err := net.Listen("tcp", peerNode.networkAddr)
	if err != nil {
		fmt.Printf("Error while listening to the on port: %d", peerNode.port)
		return
	}

	//until a SIGNAL interrupt is passed or an exception is
	//raised, keep on accepting peerBuild connections and add it
	//to the peer map.
	for {

		//debug information
		fmt.Printf("\nListening on Port: %d\n", peerNode.port)

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
Function which handles the incoming requests
to the peer
Params: net.Conn
Returns: Nil
*/
func handleConnection(conn net.Conn) {
	// Will write to network.
	enc = gob.NewEncoder(conn)
	// Will read from network.
	dec = gob.NewDecoder(conn)

	//read and decode the packet sent
	var recv utils.Packet
	err := dec.Decode(&recv)
	if err != nil {
		fmt.Println("Error Decoding the incoming packet of the peer: ", err.Error())
	}

	//check the packet type
	switch recv.Ptype {
	case utils.UPDATE:
		defer conn.Close()
		updateBackupPeer(enc, recv.Pcontent)
		fmt.Println("Backup Peer updated")
	case utils.STORE:
		defer conn.Close()
		//storeAndIndexFile(enc, recv.Pcontent)
		fmt.Println("Store request handled")
	}
}

/*
Function which updates the backup peer in the
current instance and sends the confirmation
to the peer
Params:
	encB: Encoder of the network connection
		  established
    content: the content received in the packet
			  sent by sender
Returns: Nil
*/
func updateBackupPeer(encB *gob.Encoder, content string) {
	//update the backup peer in the
	//current instance
	peerNode.backupPeer = content

	//send the confirmation to the backup
	//peer
	pkt := utils.CreatePacket(utils.RESPONSE, "", 0)
	err := encB.Encode(pkt)
	if err != nil {
		fmt.Println("Error while encoding response packet to the"+
			"backup peer: ", err.Error())
	}
}

/*
Function which updates the primary
peer about its new backup peer
*/
func updatePrimary() {

	//Dial the connection to the serverBuild
	conn, err := net.Dial("tcp", peerNode.myPrimaryPeer)
	if err != nil {
		fmt.Printf("Error establishing a connection to the primary peer\n")
		return
	}

	//send the packet the primary peer
	enc_1 := gob.NewEncoder(conn)
	dec_1 := gob.NewDecoder(conn)

	for {
		pkt := utils.CreatePacket(utils.UPDATE, peerNode.networkAddr, 0)
		err = enc_1.Encode(pkt)
		if err != nil {
			fmt.Println("Error encoding the update packet: ", err.Error())
		}

		//receive and decode the packet from the
		//primary peer for all ok status
		err = dec_1.Decode(&pkt)
		if err != nil {
			fmt.Println("Error decoding the response packet: ", err.Error())
		}

		//If the received packet type
		//is not a response, consider the
		//primary peer update has failed
		if pkt.Ptype != utils.RESPONSE {
			fmt.Println("Primary Peer status update unsuccessful")
		} else {
			conn.Close()
			break
		}
	}

}
