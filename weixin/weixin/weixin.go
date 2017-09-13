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

type ComStruct struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   uint
	MsgType      string
	MsgId        int64
}
type TextXml struct {
	ComStruct
	Content string
}
type ImgXml struct {
	ComStruct
	PicUrl  string
	MediaId string
}
type VoiceXml struct {
	ComStruct
	MediaId     string
	Format      string
	Recognition string
}
type VideoXml struct {
	ComStruct
	MediaId      string
	ThumbMediaId string
}
type ShortVideoXml struct {
	VideoXml
}
type LocationXml struct {
	ComStruct
	Location_X string
	Location_Y string
	Scale      int
	Label      string
}
type LinkXml struct {
	ComStruct
	Title       string
	Description string
	Url         string
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
type NewsXml struct {
	ComStruct
	ArticleCount uint
	Articles     ArticleXml
}
type EventXml struct {
	ComStruct
	Event string
}

func init() {
	InitAppInfo()
}
func GenNewsXml(msgType string) NewsXml {
	var news_xml NewsXml
	item := ItemXml{msgType + " 消息收到",
		"请点击查看公众号历史消息",
		"",
		"https://mp.weixin.qq.com/mp/profile_ext?action=home&__biz=MzAxNTU2NzAxNQ==&scene=123#wechat_redirect"}
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
		Fatal.Printf("get key :%s error: %s ", key, errRedis)
	}
	key = "appsecret"
	AppSecret, errRedis = RedisClient.Get(key).Result()
	if errRedis != nil {
		Fatal.Printf("get key :%s error: %s ", key, errRedis)
	}
	getAccessToken()
}
func getAccessToken() {
	var err error
	AccessToken, err = RedisClient.Get("accesstoken").Result()
	if err == redis.Nil || AccessToken == "" {
		Trace.Printf("start get access token from weixin ")
		url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + Appid + "&secret=" + AppSecret
		resp, errHttp := http.Get(url)
		if errHttp != nil {
			Fatal.Printf("http get access token err: %s ", errHttp)
		}

		status := resp.Status
		if status != "200 OK" {
			Fatal.Printf("http get access status not ok : %s", status)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			Fatal.Printf("read body err: %s ", err)
		} else {
			Trace.Printf("response body: %s", string(body))
		}
		type ResToken struct {
			Access_token string `json:"access_token"`
			Expire_in    uint   `json:"expire_in"`
		}
		var data ResToken
		err = json.Unmarshal(body, &data)
		if err != nil {
			Fatal.Printf("decode json body err: %s ", err)
		}
		AccessToken = data.Access_token
		err = RedisClient.SetNX("accesstoken", AccessToken, 7200*time.Second).Err()
		if err != nil {
			Fatal.Printf("set redis access token err: %s ", err)
		}
		Trace.Printf("http get access token end, accesstoken: %s", AccessToken)

	} else if err != nil {
		Fatal.Printf("get accesstoken from redis error : %s ", err)
	}
}
func GetIpList() {
	url := "https://api.weixin.qq.com/cgi-bin/getcallbackip?access_token=" + AccessToken
	resp, err := http.Get(url)
	if err != nil {
		Fatal.Printf("http get ip list error: %s ", err)
	}
	if resp.Status != "200 OK" {
		Fatal.Printf("http get ip list not ok : %s", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Fatal.Printf("ip list read body err: %s ", err)
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
