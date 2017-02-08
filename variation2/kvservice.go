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
	"errors"
	"fmt"
	"github.com/msayson/kvservice/api"
	"github.com/msayson/kvservice/util/rpc_util"
	"net/rpc"
	"os"
)

type KeyValService int

// First two back-end nodes
var nodeChain struct {
	HeadIpPort string
	NextIpPort string
}

// Get RPC Call
func (kvs *KeyValService) Get(args *api.GetArgs, reply *api.ValReply) error {
	rpcClient, err := connectToNode()
	if err != nil {
		return err
	}
	defer rpcClient.Close()
	return rpcClient.Call("KeyValService.Get", args, reply)
}

// Set RPC Call
func (kvs *KeyValService) Set(args *api.SetArgs, reply *api.ValReply) error {
	rpcClient, err := connectToNode()
	if err != nil {
		return err
	}
	defer rpcClient.Close()
	return rpcClient.Call("KeyValService.Set", args, reply)
}

// TestSet RPC Call
func (kvs *KeyValService) TestSet(args *api.TestSetArgs, reply *api.ValReply) error {
	rpcClient, err := connectToNode()
	if err != nil {
		return err
	}
	defer rpcClient.Close()
	return rpcClient.Call("KeyValService.TestSet", args, reply)
}

func connectToNode() (*rpc.Client, error) {
	var rpcClient *rpc.Client
	if nodeChain.HeadIpPort == "" {
		if nodeChain.NextIpPort == "" {
			return rpcClient, storeUnavailableError()
		}
		nodeChain.HeadIpPort = nodeChain.NextIpPort
		nodeChain.NextIpPort = ""
	}
	rpcClient, err := rpc_util.Connect(nodeChain.HeadIpPort)
	if err != nil {
		nodeChain.HeadIpPort = ""
		err = storeUnavailableError()
	}
	return rpcClient, err
}

// Join RPC Call - a new back-end node is joining the server chain
// Returns:
// - "success" if the node has been added to the end of the chain
// - the ip:port of the next node if there are more nodes to visit
func (kvs *KeyValService) Join(args *api.JoinArgs, reply *api.ValReply) error {
	if nodeChain.HeadIpPort == "" {
		nodeChain.HeadIpPort = args.IpPort
		reply.Val = "success"
	} else if nodeChain.NextIpPort == "" {
		nodeChain.NextIpPort = args.IpPort
		reply.Val = nodeChain.HeadIpPort
	} else {
		// Pass to head instead of next so that each node
		// is aware of the next two in the chain
		reply.Val = nodeChain.HeadIpPort
	}
	return nil
}

func main() {
	client_ip_port, backend_ip_port := parseRuntimeParams()

	// Setup key-value service.
	kvservice := new(KeyValService)
	rpc.Register(kvservice)

	// Listen for backend node connections in a concurrent goroutine
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

func storeUnavailableError() error {
	return errors.New("Key-value store is unavailable")
}
