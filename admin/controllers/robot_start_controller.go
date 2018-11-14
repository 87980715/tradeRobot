package controllers

import (
	"github.com/astaxie/beego"
	"tradeRobot/robot/server"
	"context"
	"strings"
)

type RobotStartController struct {
	beego.Controller
}
// 需要的参数 机器人编号
func  (c *RobotStartController) Start() {

	result := make(map[string]interface{})

	result["code"] = 0
	result["message"] = "操作成功"
	defer func() {
		c.Data["json"] = result
		c.ServeJSON()
	}()

	if len(Robots) == 0 {
		result["code"] = 1001
		result["message"] = "操作失败"
		result["error"] = "未绑定机器人"
		return
	}

	id,_ := c.GetInt("robotId")
	robot , ok:= Robots[id]
	if !ok {
		result["code"] = 1001
		result["message"] = "操作失败"
		result["error"] = "无效的机器编号"
		return
	}

	ctx := robot.ctx
	symbol := strings.Split(robot.Symbol,"_")

	go func(symbol []string, c context.Context) {
		server.RobotRun(symbol, ctx)
	}(symbol, ctx)

	robot.Stutas = "start"
	Robots[id] = robot
}

