package client

import (
	"encoding/gob"
	"flag"
	"net"
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

// error check
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// init the client
func initializeClient() {
	//parse the command line arguments
	remoteAddr := flag.String("addr", "127.0.0.1", "The address of the Master to connect to."+
		"Default is localhost")

	remotePort := flag.String("port", "9999", "Port of the Master daemon.")

	flag.Parse()

	//form the network address for the node
	address := *remoteAddr + ":" + *remotePort
	clientNode = client{masterNode: address}
}

func sendPacket(conn net.Conn, packet utils.Packet) {
	encode := gob.NewEncoder(conn)
	err := encode.Encode(packet)
	check(err)
}

func establishConnection() {
	conn, err := net.Dial("tcp", clientNode.masterNode)
	check(err)
	networkAddr := conn.LocalAddr().String()
	addr := strings.Split(networkAddr, ":")
	clientNode.address = addr[0]
	clientNode.port, err = strconv.Atoi(addr[1])
	check(err)

	packet := utils.CreatePacket(utils.STORE, "", unsafe.Sizeof(utils.STORE))
	sendPacket(conn, packet)
	decode := gob.NewDecoder(conn)

	var response utils.Response
	err = decode.Decode(conn)
	check(err)
}
