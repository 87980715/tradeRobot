package routers

import (
	"tradeRobot/admin/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/zgorder/finished", &controllers.MainController{})
    beego.Router("/zgorder/pending", &controllers.MainController{})
	beego.Router("/huobiorder/finished", &controllers.MainController{})
	beego.Router("/robot/start", &controllers.RobotStartController{},"*:Start")
	beego.Router("/robot/stop", &controllers.RobotStopController{},"*:Stop")
}
