package rpc_util

import (
	"log"
	"net"
	"net/rpc"
	"time"
)

var maxConnectTries = 5

// Returns an rpc connection, or an error
// if unable to connect after a max number of tries
func Connect(ip_port string) (*rpc.Client, error) {
	var rpcClient *rpc.Client
	var err error
	for i := 0; i < maxConnectTries; i++ {
		rpcClient, err = rpc.Dial("tcp", ip_port)
		if err == nil {
			break
		}
		if i < maxConnectTries-1 {
			time.Sleep(time.Duration(500) * time.Millisecond)
		}
	}
	return rpcClient, err
}

// Serve RPC calls to incoming clients
func ServeRpc(ip_port string) {
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
