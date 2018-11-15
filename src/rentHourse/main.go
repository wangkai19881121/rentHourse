package main

import (
	"encoding/json"
	_ "rentHourse/routers"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	config := make(map[string]interface{})
	config["filename"] = beego.AppConfig.String("logs::log_path")
	config["log_level"] = beego.AppConfig.String("logs::log_level")
	config["daily"] = true
	config["rotate"] = true
	configStr, err := json.Marshal(config)
	if err != nil {
		logs.Error("marshal failed,err:", err)
		return
	}
	logs.SetLogger(logs.AdapterConsole, `{color":true}`)
	logs.SetLogger(logs.AdapterFile, string(configStr))
	logs.EnableFuncCallDepth(true)

	beego.Run()

}

func init() {
	dburl := beego.AppConfig.String("db::dburl")
	dbmaxIdle := beego.AppConfig.String("db::dbmaxIdle")
	dbmaxConn := beego.AppConfig.String("db::dbmaxConn")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	maxidle, err := strconv.Atoi(dbmaxIdle)
	maxConn, err := strconv.Atoi(dbmaxConn)
	if err != nil {
		logs.Error("strconv failed,err:", err)
		return
	}
	orm.Debug = true
	orm.RegisterDataBase("default", "mysql", dburl, maxidle, maxConn)
}
