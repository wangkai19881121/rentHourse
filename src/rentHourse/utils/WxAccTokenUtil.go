package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

const accessTokenFetchUrl string = "https://api.weixin.qq.com/cgi-bin/token"

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpriesIn   int32  `json:"expires_in"`
	Errcode     string `json:"errcode"`
	Errmsg      string `json:"errmsg"`
}

type WxAccessToken struct {
	Id          int64     `orm:"column(id);pk"`
	AccessToken string    `json:"access_token"`
	ExpiresIn   int32     `json:"expires_in"`
	EndTime     time.Time `json:"end_time"`
	Status      string    `json:"status"`
}

func GetAccessToken() (*WxAccessToken, error) {
	o := orm.NewOrm()
	wsAccessToken := WxAccessToken{}
	err := o.QueryTable("WxAccessToken").Filter("status", "01").One(&wsAccessToken)
	if err == orm.ErrNoRows {
		logs.Informational("accesstoken查询未找到")
		fetchAccessToken, fetcherr := FetchAccessToken()
		if fetcherr == nil {
			o.Begin()
			wsAccessToken.AccessToken = (*fetchAccessToken).AccessToken
			wsAccessToken.ExpiresIn = (*fetchAccessToken).ExpriesIn
			wsAccessToken.EndTime = time.Now().Add(time.Hour * 2)
			wsAccessToken.Status = "01"
			id, err := o.Insert(&wsAccessToken)
			if err == nil {
				wsAccessToken.Id = id
				o.Commit()
			} else {
				o.Rollback()
			}
			return &wsAccessToken, err
		} else {
			return nil, nil
		}

	} else {
		logs.Informational("accesstoken查询成功")
		now := time.Now()
		endtime := wsAccessToken.EndTime.Add(-time.Minute * 10)
		logs.Informational(endtime)
		if endtime.Before(now) {
			fetchAccessToken, err := FetchAccessToken()
			o.Begin()
			wsAccessToken.Status = "02"
			o.Update(&wsAccessToken, "status")
			wxInsertAccessToken := WxAccessToken{}
			wxInsertAccessToken.AccessToken = (*fetchAccessToken).AccessToken
			wxInsertAccessToken.ExpiresIn = (*fetchAccessToken).ExpriesIn
			wxInsertAccessToken.EndTime = time.Now().Add(time.Hour * 2)
			wxInsertAccessToken.Status = "01"
			id, err := o.Insert(&wxInsertAccessToken)
			if err == nil {
				wsAccessToken.Id = id
				o.Commit()
			} else {
				o.Rollback()
			}
			return &wsAccessToken, nil
		} else {
			return &wsAccessToken, nil
		}

	}
}

/**
*强制获取AccessToken
**/
func FetchAccessToken() (*AccessTokenResponse, error) {
	appID := beego.AppConfig.String("wechat::appID")
	appSecret := beego.AppConfig.String("wechat::appsecret")
	accessTokenRequest := strings.Join([]string{accessTokenFetchUrl, "?grant_type=client_credential&appid=", appID, "&secret=", appSecret}, "")
	resp, err := http.Get(accessTokenRequest)

	if err != nil || resp.StatusCode != http.StatusOK {
		logs.Error("获取微信获取accessToken错误 %+v\n", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("获取微信获取accessToken错误 %+v\n", err)
		return nil, err
	}
	strbody := string(body)
	logs.Informational("微信accessTokenHttpget请求返回", strbody)
	if strings.Contains(strbody, "access_token") {
		atr := AccessTokenResponse{}
		err = json.Unmarshal(body, &atr)
		if err != nil {
			logs.Error("发送get请求获取 atoken 返回数据json解析错误 %+v\n", err)
			return nil, err
		}
		logs.Info(atr.AccessToken)
		return &atr, nil

	} else {
		ater := AccessTokenResponse{}
		err = json.Unmarshal(body, &ater)
		logs.Error("获取微信获取accessToken微信返回的错误信息 %+v\n", ater)
		if err != nil {
			return nil, err
		}
		return &ater, err
	}

}
func init() {
	// 需要在init中注册定义的model
	orm.RegisterModel(new(WxAccessToken))
	logs.Informational("orm注册WxAccessToken\n")
}
