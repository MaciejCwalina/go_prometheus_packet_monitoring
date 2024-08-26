package main

import (
	"log"
	"net"
)

type HttpServer struct {
	serverType string
	port       string
}

func (hs HttpServer) StartHttpServer() net.Listener {
	listener, err := net.Listen(hs.serverType, hs.port)
	if err != nil {
		log.Fatal(err.Error())
	}

	return listener
}

func (hs HttpServer) HandleIncomingConnections(listener net.Listener) net.Conn {
	conn, err := listener.Accept()
	if err != nil {
		log.Fatal("Failed to accept connection")
	}

	return conn
}

func (hs HttpServer) ReadAllBytesFromClient(validDataChannel chan []byte, conn net.Conn) {
	go func() {
		for {
			byteArr := make([]byte, 4096)
			readLength, _ := conn.Read(byteArr)
			validDataChannel <- byteArr[:readLength]
		}
	}()
}
