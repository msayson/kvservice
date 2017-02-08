// A simple key-value service that clients can interact with using RPC calls
//
// Supported operations:
// - get(key)
// - set(key,val)
// - testset(key,testval,newval)
//
// Usage: go run kvservice.go [ip:port]
//
// - [ip:port] : the IP address and TCP port to use to listen for connections

package main

import (
	"fmt"
	"github.com/msayson/kvservice/api"
	"github.com/msayson/kvservice/kvstore"
	"github.com/msayson/kvservice/util/rpc_util"
	"net/rpc"
	"os"
)

type KeyValService int

var store *kvstore.KVStore

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

func parseRuntimeParams() string {
	usage := fmt.Sprintf("Usage: %s ip:port\n", os.Args[0])
	if len(os.Args) != 2 {
		fmt.Printf(usage)
		os.Exit(1)
	}
	return os.Args[1]
}

func main() {
	ip_port := parseRuntimeParams()

	// Setup key-value store and register service.
	store = kvstore.New()
	kvservice := new(KeyValService)
	rpc.Register(kvservice)

	// Serve RPC connections to clients
	rpc_util.ServeRpc(ip_port)
}
