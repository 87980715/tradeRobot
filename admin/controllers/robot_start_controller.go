package controllers

import (
	"github.com/astaxie/beego"
	"tradeRobot/robot/models"
	"tradeRobot/robot/initialize"
	"tradeRobot/robot/server"
	"strconv"
	"strings"
	"context"
	"encoding/json"
	"github.com/astaxie/beego/logs"
)

var (
	Robots = make(map[int]*Robot)
	RobotId  = 100
)


type Robot struct {
	Symbol string `json:"symbol"`
	cancel context.CancelFunc
	RobotId   int `json:"robot_id"`
}
type RobotStartController struct {
	beego.Controller
}

func (c *RobotStartController) Start() {

	result := make(map[string]interface{})

	result["code"] = 0
	result["message"] = "start success"
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
		result["message"] = "invalid HuobiAccessKeyId or HuobiSecretkey"
		return
	}
	models.HuobiUserID = strconv.Itoa(id)

	//fmt.Println("models.HuobiUserID:", models.HuobiUserID)

	if !initialize.VerfiZGKey() {
		result["code"] = 1001
		result["message"] = "invalid ZGAccessKeyId or ZGSecretkey"
		return
	}

	parentCtx := context.WithValue(context.Background(), "symbol", symbol[0]+symbol[1])

	ctx, cancel := context.WithCancel(parentCtx)

	robot := &Robot{}

	robot.cancel = cancel
	robot.Symbol = strings.ToUpper(symbol[0]+ "_" + "CNZ" )
	RobotId ++
	robot.RobotId = RobotId

	Robots[RobotId] = robot

	go func(symbol []string, c context.Context) {
		server.RobotRun(symbol, ctx)
	}(symbol, ctx)

	bytes,err := json.Marshal(robot)
	if err != nil {
		logs.Error("json marshal robot failed err:",err)
	}

	result["robot"] = string(bytes)
}
