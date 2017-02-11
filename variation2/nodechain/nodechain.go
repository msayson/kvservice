package nodechain

import (
	"errors"
	"fmt"
	"github.com/msayson/kvservice/api"
	"github.com/msayson/kvservice/util/rpc_util"
	"net/rpc"
)

// First two back-end nodes
type NodeChain struct {
	HeadIpPort string
	NextIpPort string
}

// Get RPC call: retrieves a key-value from the network
func (chain *NodeChain) Get(args *api.GetArgs, reply *api.ValReply) error {
	rpcClient, err := chain.connectToNode()
	if err != nil {
		return err
	}
	defer rpcClient.Close()
	return rpcClient.Call("KeyValService.Get", args, reply)
}

// Set RPC call: sets a key-value in the network
func (chain *NodeChain) Set(args *api.SetArgs, reply *api.ValReply) error {
	rpcClient, err := chain.connectToNode()
	if err != nil {
		return err
	}
	defer rpcClient.Close()
	return rpcClient.Call("KeyValService.Set", args, reply)
}

// TestSet RPC call: test-sets a key-value in the network
func (chain *NodeChain) TestSet(args *api.TestSetArgs, reply *api.ValReply) error {
	rpcClient, err := chain.connectToNode()
	if err != nil {
		return err
	}
	defer rpcClient.Close()
	return rpcClient.Call("KeyValService.TestSet", args, reply)
}

// Join RPC call: add a new back-end node to the network
// Returns:
// - "success" if the node has been added to the end of the chain
// - the ip:port of the next node if there are more nodes to visit
func (chain *NodeChain) Join(args *api.JoinArgs, reply *api.ValReply) error {
	if chain.HeadIpPort == "" {
		chain.HeadIpPort = args.IpPort
		reply.Val = "success"
	} else if chain.NextIpPort == "" {
		chain.NextIpPort = args.IpPort
		reply.Val = chain.HeadIpPort
	} else {
		// Direct node to head instead of next so that each node
		// is aware of the next two in the chain
		reply.Val = chain.HeadIpPort
	}
	return nil
}

// Print contents of chain for debugging purposes
func (chain *NodeChain) Print(prefix string) {
	fmt.Printf("%sNodeChain{%s, %s}\n", prefix, chain.HeadIpPort, chain.NextIpPort)
}

func (chain *NodeChain) connectToNode() (*rpc.Client, error) {
	var rpcClient *rpc.Client
	if chain.HeadIpPort == "" {
		if chain.NextIpPort == "" {
			return rpcClient, storeUnavailableError()
		}
		chain.HeadIpPort = chain.NextIpPort
		chain.NextIpPort = ""
	}
	rpcClient, err := rpc_util.Connect(chain.HeadIpPort)
	if err != nil {
		chain.HeadIpPort = ""
		err = storeUnavailableError()
	}
	return rpcClient, err
}

func storeUnavailableError() error {
	return errors.New("Key-value store is unavailable")
}
