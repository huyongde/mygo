package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	. "./log"

	"./weixin"
)

func logInit(traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {
	Trace = log.New(traceHandle, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(infoHandle, "Info: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(warningHandle, "Warning: ", log.Ldate|log.Ltime|log.Lshortfile)
	Fatal = log.New(errorHandle, "Fatal: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func helloworld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
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
	var xml_o weixin.ComStruct
	err = xml.Unmarshal(req_body, &xml_o)
	if err != nil {
		Fatal.Fatal("xml decode error", err)
	}
	var msgType = xml_o.MsgType
	var res_xml interface{}
	switch msgType {
	case "text":
		res_xml = new(weixin.TextXml)
	case "image":
		res_xml = new(weixin.ImgXml)
	case "voice":
		res_xml = new(weixin.VoiceXml)
	case "video":
		res_xml = new(weixin.VideoXml)
	case "shortvideo":
		res_xml = new(weixin.ShortVideoXml)
	case "location":
		res_xml = new(weixin.LocationXml)
	case "link":
		res_xml = new(weixin.LinkXml)
	case "event":
		res_xml = new(weixin.EventXml)
	default:
		fmt.Println(string(req_body))
	}
	err = xml.Unmarshal(req_body, res_xml)
	fmt.Printf("%s : %+v \n ", msgType, res_xml)

	xml_str := string(req_body)
	Trace.Printf("receive body : %s", xml_str)

	var res_body []byte
	var err2 error
	if msgType == "event" {
		news_xml := new(weixin.TextXml)
		news_xml.Content = "欢迎关注，共同学习进步 ..."
		news_xml.MsgType = "text"
		news_xml.FromUserName = xml_o.ToUserName
		news_xml.ToUserName = xml_o.FromUserName
		news_xml.CreateTime = xml_o.CreateTime
		res_body, err2 = xml.Marshal(news_xml)
	} else {
		news_xml := weixin.GenNewsXml(msgType)
		//news_xml := new(weixin.TextXml)
		//news_xml.Content = msgType + "消息已收到， 请查看历史消息"
		//news_xml.MsgType = "text"
		news_xml.FromUserName = xml_o.ToUserName
		news_xml.ToUserName = xml_o.FromUserName
		news_xml.CreateTime = xml_o.CreateTime
		res_body, err2 = xml.Marshal(news_xml)
	}

	if err2 != nil {
		Fatal.Fatal(err2)
	}
	w.Write(res_body)
	Trace.Printf("response body: %s \n", string(res_body))
}
func main() {

	fmt.Println(weixin.AccessToken)
	Info.Println("ListenAndServe on 80 ")
	http.HandleFunc("/weixin", dealweixin) //设置访问的路由
	http.HandleFunc("/", helloworld)       //设置访问的路由
	err := http.ListenAndServe(":80", nil) //设置监听的端口
	if err != nil {
		Fatal.Fatal("ListenAndServe: ", err)
	}
}
