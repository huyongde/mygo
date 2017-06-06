package main

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net/http"
)

const c1 = "const 1"
const c2 = "const 2"
const c3 = "const 3"

func main() {
	fmt.Println("test")
	fmt.Println("hello World")
	fmt.Println("Hello GO")
	fmt.Println("my favourite number is ", rand.Intn(100))
	const a = 100
	var b int32 = 1000
	fmt.Println(b)
	fmt.Println(a)
	fmt.Println(math.Pi)
	fmt.Println("test")
	fmt.Println(a)
	const c = 10201

	var s = "string test"

	fmt.Println(s)

	resp, err := http.Get("http://www.qq.com")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	fmt.Println(1111111)
	fmt.Println(resp)
	fmt.Println(errors.New("error test test"))

	// var buf bytes.Buffer
	// logger := log.New(io.WriterAt("./log.txt"), "logger test: ", log.Llongfile)
	// logger.Print("hello I am a log file")
	// fmt.Println(&buf)
	// http.HandleFunc("/hello", HelloServer)
	// err2 := http.ListenAndServe(":12345", nil)
	// if err2 != nil {
	// 	fmt.Println("error", err)
	// } else {
	// 	fmt.Println("success listen on 12345")

	// 	fmt.Println(http.MethodGet)

	// }

}

/**
 *
 */
// func HelloServer(w http.ResponseWriter, req *http.Request) {
// 	io.WriteString(w, "hello golang from http server \n")
// }
