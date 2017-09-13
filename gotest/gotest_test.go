package gotest

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
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
	c := http.Client{}
	c.Timeout = 3 * time.Second
	resp, err := c.Get("http://www.google.com")
	//resp, err := http.Get("http://www.google.com") // 30秒后超时
	fmt.Println(resp, err)

	t.Log("测试客户端的超时设置")
}
func Test_ClientTimeout2(t *testing.T) {
	c := make(chan struct{})
	timer := time.AfterFunc(1*time.Second, func() { // 1秒后关闭c
		fmt.Println("Close chan c")
		close(c)
	})
	req, err := http.NewRequest("GET", "https://www.google.com.hk/?gws_rd=cr,ssl", nil)
	if err != nil {
		fmt.Println("new request error")
		os.Exit(1)
	}
	req.Cancel = c // req 的Cancel设置为chan c
	fmt.Println("start send request")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("send request error")
		fmt.Println(err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	fmt.Println("start read body")
	for {
		timer.Reset(50 * time.Millisecond)
		_, err = io.CopyN(ioutil.Discard, resp.Body, 256)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

	}
	t.Log("测试客户端的超时设置2")
}

func Test_ClientTimeout3(t *testing.T) {
	ctx, cancel := context
}
