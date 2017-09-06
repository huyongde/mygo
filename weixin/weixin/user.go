package weixin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	. "../log"
)

type OpenidList struct {
	Openid []string `json:"openid"`
}
type UserList struct {
	Total       int        `json:"total"`
	Count       int        `json:"count"`
	Data        OpenidList `json:"data"`
	Next_openid string     `json:"next_openid"`
}
type UserMark struct {
	Openid string `json:"openid"`
	Remark string `json:"remark"`
}
type RespBody struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func CommonGet(url string) (body []byte) {
	Trace.Println("url: ", url)
	resp, err := http.Get(url)
	if err != nil {
		Fatal.Println("http get error", err)
		return []byte("")
	}
	if resp.Status != "200 OK" {
		Fatal.Println("http get not ok status: ", resp.Status)
		return []byte("")
	}
	body, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		Fatal.Println("read body failed err: ", err, string(body))
		return []byte("")
	}
	return

}
func GetUserList() {
	url := "https://api.weixin.qq.com/cgi-bin/user/get?access_token=" + AccessToken
	body := CommonGet(url)
	fmt.Println(string(body))
	var user_list UserList
	err := json.Unmarshal(body, &user_list)
	if err != nil {
		Fatal.Println("json decode failed err: ", err, string(body))
	}
	fmt.Printf("%+v \n", user_list, err)
	for i := 0; i < user_list.Count; i++ {
		GetUserBaseInfo(user_list.Data.Openid[i])
		SetUserMark(user_list.Data.Openid[i], string(i))
	}
}

func GetUserBaseInfo(openid string) {
	lang := "zh_CN"
	url := "https://api.weixin.qq.com/cgi-bin/user/info?access_token=" + AccessToken + "&openid=" + openid + "&lang=" + lang
	body := CommonGet(url)
	fmt.Println(string(body))

}
func SetUserMark(openid, remark string) {
	url := "https://api.weixin.qq.com/cgi-bin/user/info/updateremark?access_token=" + AccessToken
	user_mark := UserMark{openid, remark}
	data, _ := json.Marshal(user_mark)
	resp, _ := http.Post(url, "application/json", bytes.NewBuffer(data))
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	var resp_body RespBody
	json.Unmarshal(body, &resp_body)
	if resp_body.Errcode != 0 {
		fmt.Println("remark failed err: ", string(body))
	} else {
		fmt.Println("remark success")
	}
	defer resp.Body.Close()

}
