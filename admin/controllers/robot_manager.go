package controllers

import (
	"github.com/astaxie/beego"
	"context"
	"strings"
	"strconv"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"tradeRobot/robot/models"
	"tradeRobot/robot/initialize"
)

type RobotManagerController struct {
	beego.Controller
}

var (
	Robots  = make(map[int]Robot)
	RobotId = 100
)

type Robot struct {
	Symbol  string `json:"symbol"`
	cancel  context.CancelFunc
	RobotId int    `json:"robot_id"`
	ctx     context.Context
}

func (c *RobotManagerController) Add() {
	result := make(map[string]interface{})

	result["code"] = 0
	result["message"] = "操作成功"
	defer func() {
		c.Data["json"] = result
		c.ServeJSON()
	}()

	models.Huobi_AccessKeyId = c.Input().Get("HuobiAccessKeyId")
	models.Huobi_Secretkey = c.Input().Get("HuobiSecretkey")

	models.ZG_API_KEY = c.Input().Get("ZGApiKey")
	models.ZG_SECRET_KEY = c.Input().Get("ZGSecret_key")

	s := c.Input().Get("symbol")

	// 从参数获取
	symbol := strings.Split(s, "-")

	// 策略
	models.TradeInspectTime = 500 //单位毫秒
	models.TradePriceAdjust = 0.001
	models.TradeAmountMultiple = 1

	id, err := initialize.HuobiUserId()
	if err != nil {
		result["code"] = 1001
		result["message"] = "操作失败"
		result["error"] = "invalid HuobiAccessKeyId or HuobiSecretkey"
		return
	}
	models.HuobiUserID = strconv.Itoa(id)

	if !initialize.VerfiZGKey() {
		result["code"] = 1001
		result["message"] = "操作失败"
		result["error"] = "invalid ZGAccessKeyId or ZGSecretkey"
		return
	}

	robot := Robot{}
	robot.Symbol = strings.ToUpper(symbol[0] + "_" + "CNZ")

	parentCtx := context.WithValue(context.Background(), "symbol", robot.Symbol)
	ctx, cancel := context.WithCancel(parentCtx)
	robot.ctx = ctx
	robot.cancel = cancel

	RobotId ++
	robot.RobotId = RobotId
	Robots[RobotId] = robot

	result["robotId"] = robot.RobotId
	result["symbol"] = robot.Symbol
}

func (c *RobotManagerController) Delete() {

	result := make(map[string]interface{})

	result["code"] = 0
	result["message"] = "操作成功"
	defer func() {
		c.Data["json"] = result
		c.ServeJSON()
	}()

	tempId := c.Input().Get("robotId")
	id,_ := strconv.Atoi(tempId )

	_ , ok:= Robots[id]
	if !ok {
		result["code"] = 0
		result["message"] = "操作失败"
		result["error"] = "无效的机器编号"
	}

	// 从机器人列表中减去相应的robot
	delete(Robots,id)
}

func (c *RobotManagerController) RobotsList() {

	result := make(map[string]interface{})

	result["code"] = 0
	result["message"] = "操作成功"
	defer func() {
		c.Data["json"] = result
		c.ServeJSON()
	}()

	bytes, err := json.Marshal(Robots)
	if err != nil {
		logs.Error("json marshal robot failed err:", err)
	}

	result["Robots"] = string(bytes)
}