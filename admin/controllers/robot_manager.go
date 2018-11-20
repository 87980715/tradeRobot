package controllers

import (
	"github.com/astaxie/beego"
	"context"
	"strings"
	"strconv"
	"tradeRobot/robot/models"
	"tradeRobot/robot/initialize"
)

// 新增 删除
// 启动 暂停

type RobotManagerController struct {
	beego.Controller
}

var (
	Robots  = make(map[int]Robot)
	RobotId = 100
)

type Robot struct {
	symbol     string
	cancel     context.CancelFunc
	RobotId    int    `json:"robot_id"`
	ctx        context.Context
	Stutas     string `json:"stutas"`
	ShowSymbol string `json:"symbol"`
}
// 差总得交易对集合 map
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

	models.ZT_API_KEY = c.Input().Get("ZTApiKey")
	models.ZT_SECRET_KEY = c.Input().Get("ZTSecretKey")

	s := c.Input().Get("symbol")
	symbol := strings.Split(s, "-")

	// 策略参数
	i := c.Input().Get("inspectTime")
	models.TradeInspectTime, _ = strconv.ParseInt(i, 10, 64) //单位毫秒
	p := c.Input().Get("priceAdjust")
	models.TradePriceAdjust, _ = strconv.ParseFloat(p, 64)
	a := c.Input().Get("amountMultiple")
	models.TradeAmountMultiple, _ = strconv.ParseFloat(a, 64)
	// Usdt第一个值
	r := c.Input().Get("usdtPrice")
	usdtPrice, _ := strconv.ParseFloat(r, 64)
	models.UsdtPrice["Huobi"] = usdtPrice

	// 初始化账户
	initialize.InitAccounts()

	id, err := initialize.HuobiUserId()
	if err != nil {
		result["code"] = 1001
		result["message"] = "操作失败"
		result["error"] = "invalid HuobiAccessKeyId or HuobiSecretkey"
		return
	}
	models.HuobiUserID = strconv.Itoa(id)

	id, err = initialize.ZGUserId()
	if err != nil {
		result["code"] = 1001
		result["message"] = "操作失败"
		result["error"] = "invalid AccessKeyId or Secretkey"
		return
	}
	models.ZGUserID = strconv.Itoa(id)

	// 新建机器人
	robot := Robot{}
	robot.symbol = strings.ToUpper(symbol[0] + "_" + "CNT" + "_" + symbol[1])
	robot.ShowSymbol = strings.ToUpper(symbol[0] + "_" + "CNT")
	parentCtx := context.WithValue(context.Background(), "symbol", robot.symbol)
	ctx, cancel := context.WithCancel(parentCtx)
	robot.ctx = ctx
	robot.cancel = cancel
	robot.Stutas = "stop"

	RobotId ++
	robot.RobotId = RobotId
	Robots[RobotId] = robot

	result["robotId"] = robot.RobotId
	result["symbol"] = robot.ShowSymbol
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
	id, _ := strconv.Atoi(tempId)

	_, ok := Robots[id]
	if !ok {
		result["code"] = 0
		result["message"] = "操作失败"
		result["error"] = "无效的机器编号"
	}
	// 从机器人列表中减去相应的robot
	delete(Robots, id)
}

func (c *RobotManagerController) RobotsList() {

	result := make(map[string]interface{})

	result["code"] = 0
	result["message"] = "操作成功"
	defer func() {
		c.Data["json"] = result
		c.ServeJSON()
	}()

	result["Robots"] = Robots
}
