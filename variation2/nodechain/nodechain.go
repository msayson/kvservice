package nodechain

import (
	"errors"
	"fmt"
	"github.com/msayson/kvservice/api"
	"github.com/msayson/kvservice/util/rpc_util"
	"net/rpc"
	"sync"
)

// First two back-end nodes
type NodeChain struct {
	HeadIpPort string
	NextIpPort string
	lock       *sync.RWMutex // read/write mutex for safe concurrent access
}

func New() *NodeChain {
	var chain NodeChain
	chain.HeadIpPort = "" // Initialize node ip:port values
	chain.NextIpPort = ""
	chain.lock = &sync.RWMutex{} // Initialize read/write mutex
	return &chain
}

// Retrieves key-value from the network
func (chain *NodeChain) Get(args *api.GetArgs, reply *api.ValReply) error {
	rpcClient, err := chain.connectToFirstLiveNode()
	if err != nil {
		return err
	}
	defer rpcClient.Close()
	return rpcClient.Call("KeyValService.Get", args, reply)
}

// Sets key-value in the network
func (chain *NodeChain) Set(args *api.SetArgs, reply *api.ValReply) error {
	rpcClient, err := chain.connectToFirstLiveNode()
	if err != nil {
		return err
	}
	defer rpcClient.Close()
	return rpcClient.Call("KeyValService.Set", args, reply)
}

// Test-sets key-value in the network
func (chain *NodeChain) TestSet(args *api.TestSetArgs, reply *api.ValReply) error {
	rpcClient, err := chain.connectToFirstLiveNode()
	if err != nil {
		return err
	}
	defer rpcClient.Close()
	return rpcClient.Call("KeyValService.TestSet", args, reply)
}

// Adds a new back-end node to the network
// Returns "success" if the node has been added to the end of the chain, or
//   the ip:port of the next node if there are more nodes to visit
func (chain *NodeChain) Join(args *api.JoinArgs, reply *api.ValReply) error {
	if args.IpPort == "" {
		return errors.New("Join: expected an ip:port, received empty string")
	}

	chain.lock.Lock()
	if chain.HeadIpPort == "" {
		chain.HeadIpPort = args.IpPort
		reply.Val = "success"
	} else if chain.NextIpPort == "" {
		chain.NextIpPort = args.IpPort
		reply.Val = chain.HeadIpPort
	} else {
		// Forward node to head instead of next so that each node
		// is aware of the next two in the chain
		reply.Val = chain.HeadIpPort
	}
	chain.lock.Unlock()
	return nil
}

// GetNextNodes RPC call: returns ip:port addresses of next nodes in chain
func (chain *NodeChain) GetNextNodes(reply *api.GetNextNodesReply) error {
	chain.lock.RLock()
	reply.HeadIpPort = chain.HeadIpPort
	reply.NextIpPort = chain.NextIpPort
	chain.lock.RUnlock()
	return nil
}

// Print contents of chain for debugging purposes
func (chain *NodeChain) Print(prefix string) {
	fmt.Printf("%sNodeChain{%s, %s}\n", prefix, chain.HeadIpPort, chain.NextIpPort)
}

// Connect to the first live node in the chain,
// removing unresponsive nodes as they are encountered
func (chain *NodeChain) connectToFirstLiveNode() (*rpc.Client, error) {
	var rpcClient *rpc.Client
	chain.lock.Lock()
	defer chain.lock.Unlock()

	if chain.HeadIpPort == "" {
		return rpcClient, storeUnavailableError()
	}
	rpcClient, err := rpc_util.Connect(chain.HeadIpPort)
	if err != nil && chain.NextIpPort != "" {
		chain.HeadIpPort = chain.NextIpPort
		chain.NextIpPort = ""
		rpcClient, err = rpc_util.Connect(chain.HeadIpPort)
	}
	go chain.updateEndOfChain()
	return rpcClient, err
}

// Connect to the last node in the chain that is live,
// removing unresponsive tail nodes as they are encountered
func (chain *NodeChain) connectToLastInChain() (*rpc.Client, error) {
	chain.lock.Lock()
	defer chain.lock.Unlock()
	rpcClient, err := rpc_util.Connect(chain.NextIpPort)
	if err == nil {
		return rpcClient, err
	}
	chain.NextIpPort = "" // Next is unresponsive, remove from chain
	rpcClient, err = rpc_util.Connect(chain.HeadIpPort)
	if err != nil {
		chain.HeadIpPort = "" // Head is unresponsive, remove from chain
	}
	return rpcClient, err
}

// If chain is not full, contact the last known node
// to obtain the addresses to any subsequent nodes
func (chain *NodeChain) updateEndOfChain() {
	connLastLive, err := chain.connectToLastInChain()
	if err != nil {
		fmt.Printf("Error updating chain: %s\n", err.Error())
		return
	}
	defer connLastLive.Close()
	reply := api.GetNextNodesReply{}
	err = connLastLive.Call("KeyValService.GetNextNodes", 0, &reply)
	if err != nil {
		fmt.Printf("Error calling KeyValService.GetNextNodes: %s\n", err.Error())
	} else {
		chain.appendToLocalChain(reply.HeadIpPort)
		chain.appendToLocalChain(reply.NextIpPort)
	}
}

// Add ip:port to end of current node's local chain
func (chain *NodeChain) appendToLocalChain(ipPort string) {
	if ipPort != "" {
		chain.Join(&api.JoinArgs{ipPort}, &api.ValReply{})
	}
}

func storeUnavailableError() error {
	return errors.New("Key-value store is unavailable")
}
