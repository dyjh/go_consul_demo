package main

import (
	"fmt"
	"net/rpc/jsonrpc"
)

func client() {
	client, err := jsonrpc.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer client.Close()

	var request = "Greetings"
	var reply string

	err = client.Call("HelloService.SayHello", request, &reply)
	if err != nil {
		fmt.Println("Error calling remote procedure:", err)
		return
	}

	fmt.Println(reply)
}
