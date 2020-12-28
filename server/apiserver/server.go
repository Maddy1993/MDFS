package apiserver

import (
	"context"
	"encoding/gob"
	"fmt"
	"net"
	"strconv"
	"sync"
)

type master struct {
	address     string
	port        int
	networkAddr string
	peers       map[string]int
	backupPeers map[string]string
}

var (
	masterNode   master
	enc          *gob.Encoder
	dec          *gob.Decoder
	mutex        = &sync.Mutex{}
	previousPeer string
)

//NewMasterServer returns a new master instance with server options.
func NewMasterServer(port string, servAddr string) (*master, error) {
	//form the network address for the node
	address := servAddr + ":" + port

	//initialize the global variable
	//representing master node
	p, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}

	return &master{address: servAddr, port: p, networkAddr: address,
		peers: make(map[string]int), backupPeers: make(map[string]string)}, nil
}

func (fs *master) Run(ctx context.Context) {

	//heart beat signal handler
	go heartBeatHandler()

	//listen on the designates network address
	adapter, err := net.Listen("tcp", fs.networkAddr)
	if err != nil {
		fmt.Printf("Error while listening to the on port: %d", fs.port)
		return
	}

	//until a SIGNAL interrupt is passed or an exception is
	//raised, keep on accepting peerBuild connections and add it
	//to the peer map.
	for {

		//debug information
		fmt.Printf("\nListening on Port: %d\n", fs.port)

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
