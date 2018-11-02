package controllers

import (
	"github.com/astaxie/beego"
	"strconv"
)

type RobotStopController struct {
	beego.Controller
}

func (c *RobotStopController) Stop() {

	result := make(map[string]interface{})

	result["code"] = 0
	result["message"] = "stop success"
	defer func() {
		c.Data["json"] = result
		c.ServeJSON()
	}()

	robotId := c.Input().Get("robotId")
	id,_ := strconv.Atoi(robotId )

	robot,ok := Robots[id]
	if ok {
		cancel := robot.cancel
		cancel()
	}else{
		result["code"] = 1001
		result["message"] = "stop failed"
	}

	// 从机器人列表中减去相应的robot
	delete(Robots,id)
}

