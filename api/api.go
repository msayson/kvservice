package api

import (
	"errors"
	"fmt"
	"github.com/msayson/kvservice/util/rpc_util"
	"net/rpc"
)

// Struct for Get() RPC call arguments
type GetArgs struct {
	Key string // Will look up value associated with Key
}

// Struct for Set() RPC call arguments
type SetArgs struct {
	Key string // Will set value for Key
	Val string
}

// Struct for TestSet() RPC call arguments
// Semantics: if val(Key) == TestVal, will set val(Key) = NewVal
type TestSetArgs struct {
	Key     string // Key to test/set value for
	TestVal string
	NewVal  string
}

// Struct for Join() RPC call arguments
type JoinArgs struct {
	IpPort string // ip:port of node requesting to join network
}

// Struct for RPC call replies
type ValReply struct {
	Val string
}

// Initiate a Get() RPC call
func Get(kvserver *rpc.Client, key string) (string, error) {
	reply := ValReply{}
	err := kvserver.Call("KeyValService.Get", GetArgs{key}, &reply)
	if err != nil {
		err = errors.New(fmt.Sprintf("KeyValService.Get RPC call failed: %s", err.Error()))
	}
	return reply.Val, err
}

// Initiate a Set() RPC call
func Set(kvserver *rpc.Client, key, value string) (string, error) {
	reply := ValReply{}
	err := kvserver.Call("KeyValService.Set", SetArgs{key, value}, &reply)
	if err != nil {
		err = errors.New(fmt.Sprintf("KeyValService.Set RPC call failed: %s", err.Error()))
	}
	return reply.Val, err
}

// Initiate a TestSet() RPC call
func TestSet(kvserver *rpc.Client, key, testValue, newValue string) (string, error) {
	reply := ValReply{}
	err := kvserver.Call("KeyValService.TestSet", TestSetArgs{key, testValue, newValue}, &reply)
	if err != nil {
		err = errors.New(fmt.Sprintf("KeyValService.TestSet RPC call failed: %s", err.Error()))
	}
	return reply.Val, err
}

// Initiate a Join() RPC call
func JoinNetwork(kvserver *rpc.Client, ipPort string) (string, error) {
	reply := ValReply{}
	joinArgs := JoinArgs{ipPort}
	err := kvserver.Call("KeyValService.Join", joinArgs, &reply)
	if err != nil {
		reply.Val = ""
		fmt.Sprintf("KeyValService.Join RPC call failed: %v", err)
	}
	return reply.Val, err
}

// Initialiate a Join() RPC call using a known node's ip:port
func JoinNetworkByIpPort(targetIpPort, ipPort string) (string, error) {
	rpcClient, err := rpc_util.Connect(targetIpPort)
	if err != nil {
		return "", err
	}
	defer rpcClient.Close()
	replyVal, err := JoinNetwork(rpcClient, ipPort)
	return replyVal, err
}
