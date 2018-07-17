package main

import (
	"flag"
	"net"
	"fmt"
	"strings"
)

func main()  {
	server := flag.String("serverAddr", "", "The address of the server to connect to." +
		"Default is localhost")

	port := flag.String("port", "9999", "Port to listen for incoming connections.")

	flag.Parse()

	address := *server+":"+*port

	AcceptAndProcess(address, &port)
}

func AcceptAndProcess(address string, port **string){
	adapter, err := net.Listen("tcp", address)

	for{

		fmt.Println("Listening on Port: " + **port)

		if err != nil{
			fmt.Printf("Error while listening to the on port: %s", **port)
			break
		}

		conn, err := adapter.Accept()
		if err != nil {
			println(err.Error())
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn){
	clientAddr := conn.RemoteAddr().String()

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
