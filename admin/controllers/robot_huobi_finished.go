package controllers

import (
	"github.com/astaxie/beego"
	"strings"
	"tradeRobot/robot/models"
	"tradeRobot/robot/utils"
)

type HuobiFinishedController struct {
	beego.Controller
}

// 参数-火币交易对
func (c *HuobiFinishedController) HuobiFinished() {
	result := make(map[string]interface{})

	result["code"] = 0
	result["message"] = "操作成功"

	defer func() {
		c.Data["json"] = result
		c.ServeJSON()
	}()

	s := strings.ToUpper(c.GetString("HuobiSymbol"))
	tempSymbol := strings.Split(s, "-")
	symbol := tempSymbol[0] + tempSymbol[1]

	db:=utils.RobotDB

	var tradeResult []models.HuobiTradeResults

	db.Model(&models.HuobiTradeResults{}).Where(&models.HuobiTradeResults{Symbol: symbol}).Find(&tradeResult)
	result["results"] = tradeResult
}
