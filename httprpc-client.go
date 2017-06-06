package main

import "fmt"
import "net/rpc"
import "os"

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], "ip:port")
		os.Exit(1)
	}

	serverAddress := os.Args[1]

	client, err := rpc.DialHTTP("tcp", serverAddress)

	if err != nil {
		fmt.Println(err)
	}

	args := Args{18, 9}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%d * %d = %d \n", args.A, args.B, reply)

	var quo Quotient

	err = client.Call("Arith.Divide", args, &quo)

	if err != nil {
		fmt.Println("divide err", err)
	}

	fmt.Printf("Arith %d/%d=%d, remainder %d \n", args.A, args.B, quo.Quo, quo.Rem)
}
