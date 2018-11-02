package routers

import (
	"tradeRobot/admin/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/zgorder/finished", &controllers.ZGFinishedController{},"*:ZGFinished")
    beego.Router("/zgorder/pending", &controllers.MainController{})
	beego.Router("/huobiorder/finished", &controllers.HuobiFinishedController{},"*:HuobiFinished")
	beego.Router("/robot/start", &controllers.RobotStartController{},"*:Start")
	beego.Router("/robot/stop", &controllers.RobotStopController{},"*:Stop")
	beego.Router("/robot/add", &controllers.RobotManagerController{},"*:Add")
	beego.Router("/robot/delete", &controllers.RobotManagerController{},"*:Delete")
	beego.Router("/robot/list", &controllers.RobotManagerController{},"*:RobotsList")
}
