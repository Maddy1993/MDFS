package client

import (
	"bufio"
	"encoding/gob"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"unsafe"
	"utils"
)

type client struct {
	address       string
	port          int
	myPrimaryPeer string
	backupPeer    string
	masterNode    string
}

// global data
var clientNode client
var encode *gob.Encoder
var decode *gob.Decoder

func Start() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter file_name: \n")
	file_name, err := reader.ReadString('\n')
	utils.Check(err)
	file_name = strings.TrimSuffix(file_name, "\n")
	//fmt.Printf(file_name)

	initializeClient(file_name)

}

// init the client
func initializeClient(file_name string) {
	//parse the command line arguments
	remoteAddr := flag.String("addr", "127.0.0.1", "The address of the Master to connect to."+
		"Default is localhost")

	remotePort := flag.String("port", "9999", "Port of the Master daemon.")

	flag.Parse()

	//form the network address for the node
	address := *remoteAddr + ":" + *remotePort
	clientNode = client{masterNode: address}

	establishConnection(file_name)
}

func sendPacket(conn net.Conn, packet utils.Packet) {
	//gob.Register(os.FileInfo)
	encode = gob.NewEncoder(conn)
	err := encode.Encode(packet)
	utils.Check(err)
}

func establishConnection(file_name string) {
	conn, err := net.Dial("tcp", clientNode.masterNode)
	utils.Check(err)
	networkAddr := conn.LocalAddr().String()
	addr := strings.Split(networkAddr, ":")
	clientNode.address = addr[0]
	clientNode.port, err = strconv.Atoi(addr[1])
	utils.Check(err)
	total_size := unsafe.Sizeof(utils.STORE) + unsafe.Sizeof(string(file_name))
	packet := utils.CreatePacket(utils.STORE, string(file_name), total_size)
	//fileInfo, err := os.Stat(file_name)
	//utils.Check(err)
	//packet.PfileInfo = fileInfo
	sendPacket(conn, packet)
	decode = gob.NewDecoder(conn)

	var response utils.ClientResponse
	err = decode.Decode(&response)
	utils.Check(err)

	if response.Ptype == utils.RESPONSE {
		fmt.Printf("Primary %s, Secondary %s", response.PrimaryNetAddr, response.BackupNetAddr)
		clientNode.myPrimaryPeer = response.PrimaryNetAddr
		clientNode.backupPeer = response.BackupNetAddr
	}
	conn.Close()
	sendDataToPeer(file_name)
}

func sendDataToPeer(file_name string) {
	file_data, err := ioutil.ReadFile(file_name)
	utils.Check(err)
	conn, err := net.Dial("tcp", clientNode.myPrimaryPeer)
	utils.Check(err)
	total_size := unsafe.Sizeof(utils.STORE) + unsafe.Sizeof(string(file_data))
	packet := utils.CreatePacket(utils.STORE, string(file_data), total_size)
	sendPacket(conn, packet)
}
