package tests

import (
	"testing"
	"server"
	"time"
	"client"
)

func TestServer(t *testing.T)  {
	serverTest1(t)
}

/*
Function which tests for the values
after the serverBuild has started
 */
func serverTest1(t *testing.T)  {
	//Start the serverBuild
	go server.StartServer()

	time.Sleep(time.Second * 30)

	go client.Start()
}