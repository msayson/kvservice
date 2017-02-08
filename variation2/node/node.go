// A backend node for storing key-values.
//
// Supported operations:
// - get(key)
// - set(key,val)
// - testset(key,testval,newval)
//
// Usage: go run node.go [ip:port] [frontend ip:port]
//
// - [ip:port] : the IP address and TCP port to use to listen for connections
// - [frontend ip:port] : the IP address and TCP port of the frontend server

package main

import (
	"fmt"
	"github.com/msayson/kvservice/api"
	"github.com/msayson/kvservice/kvstore"
	"github.com/msayson/kvservice/util/rpc_util"
	"log"
	"net/rpc"
	"os"
)

type KeyValService int

var store *kvstore.KVStore

// First two back-end nodes
var nodeChain struct {
	HeadIpPort string
	NextIpPort string
}

// Get RPC Call
func (kvs *KeyValService) Get(args *api.GetArgs, reply *api.ValReply) error {
	reply.Val = store.Get(args.Key)
	return nil
}

// Set RPC Call
func (kvs *KeyValService) Set(args *api.SetArgs, reply *api.ValReply) error {
	reply.Val = store.Set(args.Key, args.Val)
	return nil
}

// TestSet RPC Call
func (kvs *KeyValService) TestSet(args *api.TestSetArgs, reply *api.ValReply) error {
	reply.Val = store.TestSet(args.Key, args.TestVal, args.NewVal)
	return nil
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
	ip_port, frontend_ip_port := parseRuntimeParams()

	// Setup key-value store and register service.
	store = kvstore.New()
	kvservice := new(KeyValService)
	rpc.Register(kvservice)

	// Contact front-end server to join the network
	joinNetwork(ip_port, frontend_ip_port)

	// Listen for client connections
	rpc_util.ServeRpc(ip_port)
}

// Contact the front-end server to join the network
func joinNetwork(ip_port, frontend_ip_port string) {
	nextNodeIpPort := frontend_ip_port
	for {
		joinResult, err := api.JoinNetworkByIpPort(nextNodeIpPort, ip_port)
		checkUnrecoverable(err, "Error joining network:")

		if joinResult == "success" {
			fmt.Println("Successfully joined network")
			return
		}
		nextNodeIpPort = joinResult
	}
}

// Returns ip:port to listen on, and ip:port of front-end server
func parseRuntimeParams() (string, string) {
	usage := fmt.Sprintf("Usage: %s [ip:port] [frontend ip:port]\n", os.Args[0])
	if len(os.Args) != 3 {
		fmt.Printf(usage)
		os.Exit(1)
	}
	return os.Args[1], os.Args[2]
}

func checkUnrecoverable(err error, msgIfFail string) {
	if err != nil {
		log.Fatal(msgIfFail, err)
	}
}
