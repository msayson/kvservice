// A command-line client for the key-value service
//
// Usage: go run client.go [server ip:port]
//
// - [server ip:port] : the IP address and TCP port of the server to connect to

package main

import (
	"bufio"
	"fmt"
	"github.com/msayson/kvservice/api"
	"github.com/msayson/kvservice/util/rpc_util"
	"github.com/msayson/kvservice/util/userinput"
	"net/rpc"
	"os"
)

// The RPC object for the key-value server
var kvserver *rpc.Client

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run client.go [server ip:port]")
		os.Exit(1)
	}

	// Establish connection with key-value server
	var err error
	kvserver, err = rpc.Dial("tcp", os.Args[1])
	checkError(err)

	fmt.Printf("Enter commands below.\nSupported commands:\n")
	fmt.Println("   get(id)                    - returns value for id")
	fmt.Println("   set(id,val)                - sets value for id")
	fmt.Println("   testset(id,testVal,newVal) - if id has testVal as its value, set to newVal")
	fmt.Println("   exit                       - shuts down client")
	reader := bufio.NewReader(os.Stdin)
	for {
		processUserCommand(reader)
	}

	os.Exit(0)
}

// Convert next user input to a key-value request,
// and print result to console
func processUserCommand(reader *bufio.Reader) {
	fmt.Print("> ")
	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Unexpected error reading user input: %s\n", err.Error())
		return
	}
	fullCmd, err := userinput.ParseCommand(text)
	if err != nil {
		fmt.Printf("Unexpected error parsing command: %s\n", err.Error())
		return
	}
	if fullCmd.Command == userinput.EXIT {
		fmt.Println("Received exit signal, shutting down...")
		kvserver.Close()
		os.Exit(0)
	}
	runUserCommand(fullCmd)
}

// Send key-value request and print result to console
func runUserCommand(cmd userinput.LegalCommand) {
	if cmd.Command == userinput.GET {
		val, err := api.Get(kvserver, cmd.Args[0])
		processKVResult("get(%s) -> %s\n", err, cmd.Args[0], val)
	} else if cmd.Command == userinput.SET {
		val, err := api.Set(kvserver, cmd.Args[0], cmd.Args[1])
		processKVResult("set(%s,%s) -> %s\n", err, cmd.Args[0], cmd.Args[1], val)
	} else if cmd.Command == userinput.TESTSET {
		val, err := api.TestSet(kvserver, cmd.Args[0], cmd.Args[1], cmd.Args[2])
		processKVResult("testset(%s,%s,%s) -> %s\n", err, cmd.Args[0], cmd.Args[1], cmd.Args[2], val)
	}
}

// Print server response to console, and if received error response,
// try to reconnect to server
func processKVResult(msgPattern string, err error, a ...interface{}) {
	if err != nil {
		fmt.Println(err)
		reconnectToKVServer()
	} else {
		fmt.Printf(msgPattern, a...)
	}
}

// Close current rpc connection to server and try to reconnect
func reconnectToKVServer() {
	var err error
	fmt.Println("Reconnecting to server...")
	kvserver.Close()
	kvserver, err = rpc_util.Connect(os.Args[1])
	checkError(err)
	fmt.Println("Connection successful.")
}

// If error is non-nil, print error and shut down
func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
