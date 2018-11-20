package controllers

import (
	"github.com/astaxie/beego"
	"context"
)

type RobotStopController struct {
	beego.Controller
}

func (c *RobotStopController) Stop() {

	result := make(map[string]interface{})

	result["code"] = 0
	result["message"] = "操作成功"
	defer func() {
		c.Data["json"] = result
		c.ServeJSON()
	}()

	id,_ := c.GetInt("robotId")

	if len(Robots) == 0 {
		result["code"] = 1001
		result["message"] = "操作失败"
		result["error"] = "未启动器人"
		return
	}

	robot,ok := Robots[id]
	if ok {
		cancel := robot.cancel
		cancel()
	}else{
		result["code"] = 1001
		result["message"] = "操作失败"
		result["error"] = "无效机器人编号"
		return
	}

	// 重新生成 ctx，不让重新启动不了了
	parentCtx := context.WithValue(context.Background(), "symbol", robot.symbol)
	ctx, cancel := context.WithCancel(parentCtx)
	robot.ctx = ctx
	robot.cancel = cancel
	robot.Stutas = "stop"

	Robots[id] = robot
}

