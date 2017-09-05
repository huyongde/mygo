package weixin

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	. "../log"
	. "../redis"

	"github.com/go-redis/redis"
)

var Appid, AppSecret, AccessToken string

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

func init() {
	InitAppInfo()
}
func GenNewsXml() NewsXml {
	var news_xml NewsXml
	item := ItemXml{"Nginx Test", "Nginx 配置中的全局变量", "http://mmbiz.qpic.cn/mmbiz/pNv7yyl5J3jM6trKf8dL6BogFslUMGOQB2u7xGM3t8UDVAlBggJ6TeBGD1sWy7yk32EXqSSjpQwlMkq8rfe9ag/640?wx_fmt=jpeg&tp=webp&wxfrom=5", "http://mp.weixin.qq.com/s/_cxGX1JZnzwDS2kaaEBjHg"}
	var article_xml ArticleXml
	article_xml.Item = []ItemXml{item}
	news_xml.Articles = article_xml
	news_xml.ArticleCount = 1
	news_xml.MsgType = "news"

	return news_xml

}

func InitAppInfo() {
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
func GetIpList() {
	url := "https://api.weixin.qq.com/cgi-bin/getcallbackip?access_token=" + AccessToken
	resp, err := http.Get(url)
	if err != nil {
		Fatal.Fatalf("http get ip list error: %s ", err)
	}
	if resp.Status != "200 OK" {
		Fatal.Fatalf("http get ip list not ok : %s", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Fatal.Fatalf("ip list read body err: %s ", err)
	} else {
		Trace.Printf("ip list response body: %s", string(body))
		fmt.Println(string(body))
	}
	type IpList struct {
		Ip_list []string `json:"ip_list"`
	}
	var data IpList
	err = json.Unmarshal(body, &data)
	fmt.Println(data.Ip_list)

}
