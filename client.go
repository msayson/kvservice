// Command-line client for key-value service
// Usage: go run client.go [server ip:port]
package main

import (
	"./userinput"
	"bufio"
	"fmt"
	"net/rpc"
	"os"
)

// The RPC object for the key-value server
var kvserver *rpc.Client

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
}

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
	reader := bufio.NewReader(os.Stdin)
	for {
		processUserCommand(reader)
	}

	os.Exit(0)
}

// General-purpose error handler.
func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
