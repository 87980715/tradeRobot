package controllers

import (
	"github.com/astaxie/beego"
	"strings"
	"github.com/astaxie/beego/logs"
	"encoding/json"
	"tradeRobot/robot/models"
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

	db, err := LoadRobotDB()

	if err != nil {
		logs.Error("load robotDB failed err:", err)
		result["code"] = 1001
		result["message"] = "操作失败"
	}

	defer db.Close()
	var tradeResult []models.HuobiTradeResults

	db.Model(&models.HuobiTradeResults{}).Where(&models.HuobiTradeResults{Symbol: symbol}).Find(&tradeResult)

	bytes, _ := json.Marshal(tradeResult)

	result["results"] = string(bytes)
}
