package main

import (
	"net"
	"fmt"
	"bufio"
	"os"
	"flag"
)


func main()  {

	remoteAddr := flag.String("addr", "127.0.0.1", "The address of the server to connect to." +
		"Default is localhost")

	remotePort := flag.String("port", "9999", "Port to listen for incoming connections.")

	flag.Parse()

	conn, err :=net.Dial("tcp", *remoteAddr + ":" + string(*remotePort))
	//conn, err :=net.Dial("tcp", "127.0.0.1:9999")

	if err != nil{
		fmt.Printf("Error establishing a connection\n")
		return
	}

	sendMessage(conn)
}

func sendMessage(conn net.Conn){

	var message string
	var msgSent []byte
	msgSent = make([]byte, 255)

	fmt.Printf("Enter the Message to send to the client: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	message = scanner.Text()

	msgSent = []byte(message+"\n")
	n, err:= conn.Write(msgSent)

	fmt.Printf("Size of text written is: %d\n", n)

	if err != nil{
		fmt.Printf("Error while writing message to socket")
	}

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
