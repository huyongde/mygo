package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	. "./log"

	"./weixin"
)

type TextXml struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   uint
	MsgType      string
	Content      string
	/*
		MsgId        string
		MediaId      string
		PicUrl       string
	*/
}

func logInit(traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {
	Trace = log.New(traceHandle, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(infoHandle, "Info: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(warningHandle, "Warning: ", log.Ldate|log.Ltime|log.Lshortfile)
	Fatal = log.New(errorHandle, "Fatal: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func dealweixin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form) //这些是服务器端的打印信息
	echostr := strings.Join(r.Form["echostr"], "")
	fmt.Println(echostr)
	fmt.Fprintf(w, echostr) //输出到客户端的信息
	req_body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Fatal.Fatal("Read xml data err", err)
	}
	var xml_o TextXml
	err = xml.Unmarshal(req_body, &xml_o)
	if err != nil {
		Fatal.Fatal("xml decode error", err)
	}
	xml_str := string(req_body)
	Trace.Printf("receive body : %s", xml_str)
	/*
		xml_res := xml_o
		xml_res.FromUserName = xml_o.ToUserName
		xml_res.ToUserName = xml_o.FromUserName
		xml_res.MsgType = "text"
		xml_res.Content = "持续开发中 by golang "
		res_body, err2 := xml.Marshal(xml_res)
	*/
	news_xml := weixin.GenNewsXml()
	news_xml.FromUserName = xml_o.ToUserName
	news_xml.ToUserName = xml_o.FromUserName
	news_xml.CreateTime = xml_o.CreateTime

	res_body, err2 := xml.Marshal(news_xml)
	if err2 != nil {
		Fatal.Fatal(err2)
	}
	w.Write(res_body)
	Trace.Printf("response body: %s \n", string(res_body))
}
func main() {

	fmt.Println(weixin.AccessToken)
	weixin.GetIpList()
	os.Exit(1)
	Info.Println("ListenAndServe on 80 ")
	http.HandleFunc("/weixin", dealweixin) //设置访问的路由
	err := http.ListenAndServe(":80", nil) //设置监听的端口
	if err != nil {
		Fatal.Fatal("ListenAndServe: ", err)
	}
}
