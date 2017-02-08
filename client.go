// A command-line client for the key-value service
//
// Usage: go run client.go [server ip:port]
//
// - [server ip:port] : the IP address and TCP port of the server to connect to

package main

import (
	"./kvserviceapi"
	"./userinput"
	"bufio"
	"fmt"
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
	fmt.Println("   get(id)                 - returns value for id")
	fmt.Println("   set(id,val)             - sets value for id")
	fmt.Println("   testset(id,prevVal,val) - if id has prevVal as its value, set new value")
	fmt.Println("   exit                    - shut down client")
	reader := bufio.NewReader(os.Stdin)
	for {
		processUserCommand(reader)
	}

	os.Exit(0)
}

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

func runUserCommand(cmd userinput.LegalCommand) {
	if cmd.Command == userinput.GET {
		val := kvserviceapi.Get(kvserver, cmd.Args[0])
		fmt.Printf("get(%s) -> %s\n", cmd.Args[0], val)
	} else if cmd.Command == userinput.SET {
		val := kvserviceapi.Set(kvserver, cmd.Args[0], cmd.Args[1])
		fmt.Printf("set(%s,%s) -> %s\n", cmd.Args[0], cmd.Args[1], val)
	} else if cmd.Command == userinput.TESTSET {
		val := kvserviceapi.TestSet(kvserver, cmd.Args[0], cmd.Args[1], cmd.Args[2])
		fmt.Printf("testset(%s,%s,%s) -> %s\n", cmd.Args[0], cmd.Args[1], cmd.Args[2], val)
	}
}

// If error is non-nil, print error and shut down
func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
