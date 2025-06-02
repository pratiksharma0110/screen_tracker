package internal

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"screen_tracker/pkg/utils"
	"syscall"
)

var SOCKET_TYPE string = "unix"
var SOCKET_PATH string = "/var/tmp/screen_tracker.sock"

var clientConn net.Conn

func CreateSockt(SOCKET_TYPE string, SOCKET_PATH string) net.Conn {
	fmt.Printf("Socket running on %s\n", SOCKET_PATH)
	if _, err := os.Stat(SOCKET_PATH); err == nil {
		err = os.Remove(SOCKET_PATH)
		if err != nil {
			log.Fatalf("Failed to remove existing socket file: %v", err)
		}
	}

	ln, err := net.Listen(SOCKET_TYPE, SOCKET_PATH)
	utils.Check(err)
	defer ln.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Remove(SOCKET_PATH)
		os.Exit(1)
	}()

	conn, err := ln.Accept()
	if err != nil {
		log.Println("Failed to accept connection:", err)

	}
	clientConn = conn
	return clientConn //only handles one connection
	//TODO: need to handle broken pipe error when client disconnects, implement that later
}
