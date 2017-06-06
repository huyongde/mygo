package main

import "fmt"
import "os"
import "net/rpc"

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage as ", os.Args[0], "ip:port")
	}
	service := os.Args[1]

	client, err := rpc.Dial("tcp", service)

	if err != nil {
		fmt.Println("arith error ", err)
	}

	args := Args{17, 8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		fmt.Println("err 1 :", err)
	}

	fmt.Printf("%d * %d = %d \n", args.A, args.B, reply)
	var quo Quotient

	err = client.Call("Arith.Divide", args, &quo)
	if err != nil {
		fmt.Println("err 2 ", err)
	}

	fmt.Printf("%d/%d = %d, remainer : %d \n", args.A, args.B, quo.Quo, quo.Rem)
}
