// A simple key-value service that clients can interact with using RPC calls
//
// Supported operations:
// - get(key)
// - set(key,val)
// - testset(key,testval,newval)
//
// Usage: go run kvservice.go [ip:port] [backend ip:port]
//
// - [ip:port] : the IP address and TCP port to use to listen for client connections
// - [backend ip:port] : the IP address and TCP port to use to listen for backend connections

package main

import (
	"fmt"
	"github.com/msayson/kvservice/api"
	"github.com/msayson/kvservice/util/rpc_util"
	"github.com/msayson/kvservice/variation2/nodechain"
	"net/rpc"
	"os"
)

type KeyValService int

// Network of back-end nodes which store key-values
var nodeChain nodechain.NodeChain

// Get RPC call: retrieves a key-value from the network
func (kvs *KeyValService) Get(args *api.GetArgs, reply *api.ValReply) error {
	return nodeChain.Get(args, reply)
}

// Set RPC call: sets a key-value in the network
func (kvs *KeyValService) Set(args *api.SetArgs, reply *api.ValReply) error {
	return nodeChain.Set(args, reply)
}

// TestSet RPC call: test-sets a key-value in the network
func (kvs *KeyValService) TestSet(args *api.TestSetArgs, reply *api.ValReply) error {
	return nodeChain.TestSet(args, reply)
}

// Join RPC call: add a new back-end node to the network
func (kvs *KeyValService) Join(args *api.JoinArgs, reply *api.ValReply) error {
	return nodeChain.Join(args, reply)
}

func main() {
	client_ip_port, backend_ip_port := parseRuntimeParams()

	// Setup key-value service.
	kvservice := new(KeyValService)
	rpc.Register(kvservice)

	// Listen for backend node connections in a concurrent goroutine
	nodeChain = nodechain.NodeChain{"", ""}
	go rpc_util.ServeRpc(backend_ip_port)

	// Listen for client connections
	rpc_util.ServeRpc(client_ip_port)
}

// Returns ip:port addresses to listen on for clients and backends
func parseRuntimeParams() (string, string) {
	usage := fmt.Sprintf("Usage: %s [ip:port] [backend ip:port]\n", os.Args[0])
	if len(os.Args) != 3 {
		fmt.Printf(usage)
		os.Exit(1)
	}
	return os.Args[1], os.Args[2]
}
