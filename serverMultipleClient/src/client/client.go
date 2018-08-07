package client

import (
		"encoding/gob"
	"fmt"
	"io/ioutil"
	"net"
		"strconv"
	"strings"
	"unsafe"
	"utils"
	"bufio"
	"os"
	"path/filepath"
	)


////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////

//Structures
type client struct {
	address       string
	port          int
	masterNode string
	//backupPeer    string
	//masterNode    string
}

// global variables
var (
	clientNode client
	encode *gob.Encoder
	decode *gob.Decoder
	dirPath string
)

////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////



//Entry point to the client process. This function
//takes the master server IP address and port as the
//input and acts as an interface to the client process
//by initializing the client process.
//Params:
//	@remoteAddr: string
//		Takes the master server IP address in the
//		string format
//	@remotePort: string
//		Takes the master server port in the string
//		format
//Returns: nil
func Start(remoteAddr string, remotePort string) {

	//initialize the client instance
	initializeClient(remoteAddr, remotePort)

	//instantiate the command line interface
	//to the user
	initializeCLI()

}

// init the client
func initializeClient(remoteAddr string, remotePort string) {

	//form the network address for the node
	address := remoteAddr + ":" + remotePort
	clientNode = client{masterNode: address}

	//set the directory path from which the
	//client can read the files
	//r := bufio.NewReader(os.Stdin)
	//dirPath, _ := r.ReadString('\n')
	dirPath = "C:\\Users\\mohan\\Desktop\\Courses\\Projects\\MDFS\\serverMultipleClient\\clientFiles"
}

//Function which initializes the Command-line
//interface to the user, making the client features
//available to the user in terms of commmands.
func initializeCLI()  {

	//cli declarations
	clientIpAddr, err := utils.ExternalIP()
	utils.ValidateError(err)

	cliMessage := "client@" + clientIpAddr + ">>"
	reader := bufio.NewReader(os.Stdin)

	//in an infinite loop
	for{
		fmt.Print(cliMessage)

		//read the input command.
		command, err := reader.ReadString('\n')
		utils.ValidateError(err)

		//process and validate the input command
		processAndValidate(command)
	}
}

//Function which processes the input command
//and validates it against the valid options.
func processAndValidate(command string){
	//Step-1: Remove unexpected suffixes
	command = strings.TrimSuffix(command, "\n")

	//Step-2: Split the string into tokens
	tokens := strings.Split(command, " ")

	switch tokens[0] {
	case "send":
		conn := establishConnection()

		//create a struct for the file
		f, err := os.Stat(filepath.Join(dirPath,tokens[1]))
		utils.ValidateError(err)

		fileV := utils.File{
			Name:f.Name(),
			Size:f.Size(),
		}

		println("Received primary details")
		println(fileV.Name)
		defer conn.Close()
		sendFile(conn, fileV)
		break
	case "receive":
	}
}

//Function which establishes a connection
//to the master server on demand
func establishConnection() (conn net.Conn){
	//dial a TCP connection to the master node/server
	conn, err := net.Dial("tcp", clientNode.masterNode)
	utils.ValidateError(err)

	//initialize the client instance with its
	//system local address and port
	if clientNode.address == "" && clientNode.port == 0 {
		//get the local address from the connection
		networkAddr := conn.LocalAddr().String()
		addr := strings.Split(networkAddr, ":")

		//initialize the processed values
		clientNode.address = addr[0]
		clientNode.port, err = strconv.Atoi(addr[1])
		utils.ValidateError(err)
	}

	return
}

//Function which sends the file to the
//server established on the conn instance passed
//as a parameter
//Params:
//	@conn: net.Conn
//		Instance which holds the TCP connection
//		to the master server
//	@fileV: file
//		An instance of the file structure
//		which holds the file name, the address
//		of the primary and backup peers where the
//		file is present at, in a string format.
//Returns: nil
func sendFile(conn net.Conn, fileV utils.File)  {
	var err error
	//create the packet to send to the server
	totalSize := unsafe.Sizeof(utils.STORE) + unsafe.Sizeof(string(fileV.Name))
	packet := utils.CreatePacket(utils.STORE, string(fileV.Name), totalSize)
	packet.PfileInfo = fileV

	println("Reached here")
	//send the packet
	gob.Register(utils.Packet{})
	gob.Register(utils.File{})
	encode = gob.NewEncoder(conn)
	err = encode.Encode(packet)
	utils.ValidateError(err)

	//Receive the confirmation packet
	//from the master and decode the peer
	//details
	var response utils.ClientResponse
	decode = gob.NewDecoder(conn)
	err = decode.Decode(&response)
	println("Reached here too")
	utils.ValidateError(err)


	if response.Ptype == utils.RESPONSE {
		fmt.Printf("Primary %s, Secondary %s\n", response.PrimaryNetAddr, response.BackupNetAddr)
		fileV.PrimaryPeer = response.PrimaryNetAddr
		fileV.BackupPeer = response.BackupNetAddr
	}

	//once the primary and backup peer
	//credentials have been established,
	//contact the primary and send the file.
	sendData(fileV)
}

//Function which establishes connection
//with the primary peer address received
//from the master node and sends the file
func sendData(fileV utils.File) {

	//read the file contents
	fileData, err := ioutil.ReadFile(filepath.Join(dirPath,fileV.Name))
	utils.ValidateError(err)

	//establish connection with the
	//primary peer
	conn, err := net.Dial("tcp", fileV.PrimaryPeer)
	utils.ValidateError(err)

	//create a network encoder and deocder
	encode = gob.NewEncoder(conn)
	decode = gob.NewDecoder(conn)

	//send a STORE request to the primary peer
	totalSize := unsafe.Sizeof(utils.STORE) + unsafe.Sizeof(string(fileData))
	packet := utils.CreatePacket(utils.STORE, string(fileData), totalSize)

	//register the interface with the gob
	gob.Register(utils.Packet{})
	packet.PfileInfo = fileV
	utils.ValidateError(err)

	//once the packet is ready,
	//send a STORE request to the primary
	err = encode.Encode(packet)

	//Expect a RESPONSE from the primary
	//confirming the client that the necessary
	//setups are done and file can now be sent.
	var resp utils.Packet
	err = decode.Decode(&resp)

	//validate the packet received. If
	//it is of the type response, send
	//the data on the same established connection
	if resp.Ptype == utils.RESPONSE {
		totalSize := unsafe.Sizeof(utils.DATA_END) + unsafe.Sizeof(string(fileData))
		packet := utils.CreatePacket(utils.DATA_END, string(fileData), totalSize)
		err = encode.Encode(packet)
		utils.ValidateError(err)
	}
}
