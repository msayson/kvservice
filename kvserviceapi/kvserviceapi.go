package kvserviceapi

import (
	"fmt"
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

// Struct for RPC call replies
type ValReply struct {
	Val string
}

// Initiate a Get() RPC call
func Get(kvserver *rpc.Client, key string) string {
	reply := ValReply{}
	setArgs := GetArgs{key}
	err := kvserver.Call("KeyValService.Get", setArgs, &reply)
	if err != nil {
		fmt.Sprintf("KeyValService.Get RPC call failed: %v", err)
	}
	return reply.Val
}

// Initiate a Set() RPC call
func Set(kvserver *rpc.Client, key, value string) string {
	reply := ValReply{}
	setArgs := SetArgs{key, value}
	err := kvserver.Call("KeyValService.Set", setArgs, &reply)
	if err != nil {
		fmt.Sprintf("KeyValService.Set RPC call failed: %v", err)
	}
	return reply.Val
}

// Initiate a TestSet() RPC call
func TestSet(kvserver *rpc.Client, key, testValue, newValue string) string {
	reply := ValReply{}
	testSetArgs := TestSetArgs{key, testValue, newValue}
	err := kvserver.Call("KeyValService.TestSet", testSetArgs, &reply)
	if err != nil {
		fmt.Sprintf("KeyValService.TestSet RPC call failed: %v", err)
	}
	return reply.Val
}
