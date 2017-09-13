package gotest

import (
	"fmt"
	"net/http"
	"testing"
)

func Test_one(t *testing.T) {
	t.Log("测试通过")
}
func Test_two(t *testing.T) {
	t.Log("测试也通过")
}

/*
func Benchmark_one(b *testing.B) {
	for i := 0; i < 1000; i++ {
		f1()
	}
}
*/
func Test_ClientTimeout(t *testing.T) {
	/*
		c := http.Client{}
		c.Timeout = 3 * time.Second
		resp, err := c.Get("http://www.google.com")
	*/
	resp, err := http.Get("http://www.google.com")
	fmt.Println(resp, err)

	t.Log("测试客户端的超时设置")
}
