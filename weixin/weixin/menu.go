package weixin

import "fmt"

func GetMenu() {
	url := "https://api.weixin.qq.com/cgi-bin/menu/get?access_token=" + AccessToken
	body := CommonGet(url)
	fmt.Println(string(body))
}
func SetMenu() {
	//	url := "https://api.weixin.qq.com/cgi-bin/menu/create?access_token=" + AccessToken

}
