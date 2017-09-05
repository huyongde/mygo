package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis"
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
type NewsXml struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   uint
	MsgType      string
	ArticleCount uint
	Articles     ArticleXml
}
type ArticleXml struct {
	Item []ItemXml `xml:"item"`
}
type ItemXml struct {
	Title       string
	Description string
	PicUrl      string
	Url         string
}

var Trace, Info, Warning, Fatal *log.Logger
var RedisClient = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})
var Appid, AppSecret, AccessToken string

func logInit(traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {
	Trace = log.New(traceHandle, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(infoHandle, "Info: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(warningHandle, "Warning: ", log.Ldate|log.Ltime|log.Lshortfile)
	Fatal = log.New(errorHandle, "Fatal: ", log.Ldate|log.Ltime|log.Lshortfile)
}
func genNewsXml() NewsXml {
	var news_xml NewsXml
	item := ItemXml{"Nginx Test", "Nginx 配置中的全局变量", "http://mmbiz.qpic.cn/mmbiz/pNv7yyl5J3jM6trKf8dL6BogFslUMGOQB2u7xGM3t8UDVAlBggJ6TeBGD1sWy7yk32EXqSSjpQwlMkq8rfe9ag/640?wx_fmt=jpeg&tp=webp&wxfrom=5", "http://mp.weixin.qq.com/s/_cxGX1JZnzwDS2kaaEBjHg"}
	var article_xml ArticleXml
	article_xml.Item = []ItemXml{item}
	news_xml.Articles = article_xml
	news_xml.ArticleCount = 1
	news_xml.MsgType = "news"

	return news_xml

}

func weixin(w http.ResponseWriter, r *http.Request) {
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
	news_xml := genNewsXml()
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
func initAppInfo() {
	key := "appid"
	var errRedis error
	Appid, errRedis = RedisClient.Get(key).Result()
	if errRedis != nil {
		Fatal.Fatalf("get key :%s error: %s ", key, errRedis)
	}
	key = "appsecret"
	AppSecret, errRedis = RedisClient.Get(key).Result()
	if errRedis != nil {
		Fatal.Fatalf("get key :%s error: %s ", key, errRedis)
	}
	getAccessToken()
}
func getAccessToken() {
	var err error
	AccessToken, err = RedisClient.Get("accesstoken").Result()
	if err == redis.Nil {
		Trace.Printf("start get access token from weixin ")
		url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + Appid + "&secret=" + AppSecret
		resp, errHttp := http.Get(url)
		if errHttp != nil {
			Fatal.Fatalf("http get access token err: %s ", errHttp)
		}

		status := resp.Status
		if status != "200 OK" {
			Fatal.Fatalf("http get access status not ok : %s", status)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			Fatal.Fatalf("read body err: %s ", err)
		} else {
			Trace.Printf("response body: %s", string(body))
		}
		type ResToken struct {
			Access_token string `json:"access_token"`
			Expire_inf   uint   `json:"expire_in"`
		}
		var data ResToken
		err = json.Unmarshal(body, &data)
		if err != nil {
			Fatal.Fatalf("decode json body err: %s ", err)
		}
		AccessToken = data.Access_token
		err = RedisClient.SetNX("accesstoken", AccessToken, 7200*time.Second).Err()
		if err != nil {
			Fatal.Fatalf("set redis access token err: %s ", err)
		}
		Trace.Printf("http get access token end, accesstoken: %s", AccessToken)

	} else if err != nil {
		Fatal.Fatalf("get accesstoken from redis error : %s ", err)
	}
}
func main() {
	logFile, logErr := os.OpenFile("log.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if logErr != nil {
		fmt.Println("open file err", logErr)
		os.Exit(1)
	}
	logInit(logFile, logFile, logFile, logFile)
	//log.SetOutput(logFile)

	_, errRedis := RedisClient.Ping().Result()
	if errRedis != nil {
		Fatal.Fatalf("redis error: %s ", errRedis)
	}
	initAppInfo()
	fmt.Println(AccessToken)
	Info.Println("ListenAndServe on 80 ")
	http.HandleFunc("/weixin", weixin)     //设置访问的路由
	err := http.ListenAndServe(":80", nil) //设置监听的端口
	if err != nil {
		Fatal.Fatal("ListenAndServe: ", err)
	}
}
