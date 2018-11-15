package routers

import (
	"rentHourse/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/wx_connect", &controllers.WxConnectController{})
	beego.Router("/wx_createMenu", &controllers.WxMenuController{})
}
