package controllers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"rentHourse/utils"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

const weMenuCreatUrl string = "https://api.weixin.qq.com/cgi-bin/menu/create"

type WxMenuController struct {
	beego.Controller
}

func (c *WxMenuController) Get() {
	wxAccessToken, error := utils.GetAccessToken()
	if error != nil {
		c.Ctx.WriteString("微信获取accessToken失败")
	}
	menuStr := beego.AppConfig.String("wechat::menu")
	postRequest, err := http.NewRequest("POST", strings.Join([]string{weMenuCreatUrl, "?access_token=", (*wxAccessToken).AccessToken}, ""), bytes.NewReader([]byte(menuStr)))
	if err != nil {
		logs.Error("微信菜单建立请求失败\n", err)
		c.Ctx.WriteString("微信菜单建立请求失败")
	}
	postRequest.Header.Set("Content-Type", "application/json; encoding=utf-8")
	client := &http.Client{}
	response, err := client.Do(postRequest)
	defer response.Body.Close()
	body, ioerr := ioutil.ReadAll(response.Body)
	if ioerr != nil {
		logs.Error("微信菜单建立解析响应错误", err)
	}
	strbody := string(body)
	logs.Informational("微信菜单建立接口返回", strbody)
	if err != nil {
		logs.Error("微信发送菜单建立请求失败", err)
		c.Ctx.WriteString("微信发送菜单建立请求失败")
	} else {
		logs.Informational("微信发送菜单建立成功")
		c.Ctx.WriteString("微信发送菜单建立成功")
	}

}
