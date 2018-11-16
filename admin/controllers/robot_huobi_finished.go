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

func (c *HuobiFinishedController) HuobiFinished() {
	result := make(map[string]interface{})

	result["code"] = 0
	result["message"] = "操作成功"

	defer func() {
		c.Data["json"] = result
		c.ServeJSON()
	}()

	s := strings.ToUpper(c.GetString("symbol"))
	tempSymbol := strings.Split(s, "-")
	symbol := tempSymbol[0] + tempSymbol[1]

	db,err := LoadRobotDB()
	if err != nil  {
		result["code"] = 1001
		result["message"] = "操作失败"
		return
	}

	var tradeResult []models.HuobiTradeResults

	db.Model(&models.HuobiTradeResults{}).Where(&models.HuobiTradeResults{Symbol: symbol}).Find(&tradeResult)
	result["results"] = tradeResult
}
