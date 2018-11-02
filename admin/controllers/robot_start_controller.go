package controllers

import (
	"github.com/astaxie/beego"
	"tradeRobot/robot/server"
	"context"
	"strings"
	"fmt"
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

	id,_ := c.GetInt("robotId")
	fmt.Println("id:",id)
	robot , ok:= Robots[id]
	if !ok {
		result["code"] = 0
		result["message"] = "操作失败"
		result["error"] = "无效的机器编号"
	}
	ctx := robot.ctx
	fmt.Println("ctx:",ctx)
	symbol := strings.Split(robot.Symbol,"_")

	go func(symbol []string, c context.Context) {
		server.RobotRun(symbol, ctx)
	}(symbol, ctx)
}

// 新增 删除
// 启动 暂停