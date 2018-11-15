package controllers

import (
	"crypto/sha1"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type WxConnectController struct {
	beego.Controller
}

func (c *WxConnectController) Get() {
	signature, timestamp, nonce, echostr := c.GetString("signature"), c.GetString("timestamp"), c.GetString("nonce"), c.GetString("echostr")
	signatureResult := makeSignature(timestamp, nonce)
	//将加密后获得的字符串与signature对比，如果一致，说明该请求来源于微信
	if signature != signatureResult {
		logs.Info("微信加密签名不正确\n")
		c.Ctx.WriteString("false")
	} else {
		//如果请求来自于微信，则原样返回echostr参数内容 以上完成后，接入验证就会生效，开发者配置提交就会成功。
		logs.Info("微信签名验证成功\n")
		c.Ctx.WriteString(echostr)
	}
}

/**
*参数字典排序,加密获得signature
**/
func makeSignature(timestamp, nonce string) string {
	token := beego.AppConfig.String("wechat::token")
	tmpArray := []string{token, timestamp, nonce}
	sort.Strings(tmpArray) //按照字典升序排列
	s := sha1.New()
	io.WriteString(s, strings.Join(tmpArray, ""))
	return fmt.Sprintf("%x", s.Sum(nil))
}
