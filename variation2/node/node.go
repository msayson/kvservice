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
	"github.com/msayson/kvservice/variation2/nodechain"
	"log"
	"net/rpc"
	"os"
)

type KeyValService int

var store *kvstore.KVStore

// Network of subsequent back-end nodes
var nodeChain *nodechain.NodeChain

// Get RPC call: retrieves a key-value from the network
func (kvs *KeyValService) Get(args *api.GetArgs, reply *api.ValReply) error {
	reply.Val = store.Get(args.Key)
	return nil
}

// Set RPC call: sets a key-value in the network
func (kvs *KeyValService) Set(args *api.SetArgs, reply *api.ValReply) error {
	reply.Val = store.Set(args.Key, args.Val)
	go nodeChain.Set(args, reply) // Propagate change to subsequent nodes
	return nil
}

// TestSet RPC call: test-sets a key-value in the network
func (kvs *KeyValService) TestSet(args *api.TestSetArgs, reply *api.ValReply) error {
	reply.Val = store.TestSet(args.Key, args.TestVal, args.NewVal)
	go nodeChain.TestSet(args, reply) // Propagate change to subsequent nodes
	return nil
}

// Join RPC call: add a new back-end node to the network
func (kvs *KeyValService) Join(args *api.JoinArgs, reply *api.ValReply) error {
	return nodeChain.Join(args, reply)
}

func main() {
	ip_port, frontend_ip_port := parseRuntimeParams()

	// Setup key-value store and register service.
	store = kvstore.New()
	kvservice := new(KeyValService)
	rpc.Register(kvservice)

	// Contact front-end server to join the network
	nodeChain = nodechain.New()
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
		log.Fatal(usage)
	}
	return os.Args[1], os.Args[2]
}

func checkUnrecoverable(err error, msgIfFail string) {
	if err != nil {
		log.Fatal(msgIfFail, err)
	}
}
