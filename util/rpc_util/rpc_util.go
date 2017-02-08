package rpc_util

import (
	"log"
	"net"
	"net/rpc"
)

func ServeRpc(ip_port string) {
	// Listen for client connections
	listener := initializeTcpListener(ip_port)
	for {
		conn, _ := listener.Accept()
		go rpc.ServeConn(conn)
	}
}

// Initialize a TCP listener on the given ip:port
func initializeTcpListener(ip_port string) net.Listener {
	listener, err := net.Listen("tcp", ip_port)
	if err != nil {
		log.Fatal("Error initializing listener:", err)
	}
	return listener
}
