package controllers

import (
	"github.com/astaxie/beego"
	"tradeRobot/robot/models"
	"strings"
	"tradeRobot/robot/utils"
)

type ZGFinishedController struct {
	beego.Controller
}

type DBConfig struct {
	User     string `default:"root"`
	Password string `default:""`
	Path     string `required:"true"`
	Port     uint   `default:"3306"`
	DbName   string `required:"true"`
	Charset  string `default:"utf8"`
}

func (c *ZGFinishedController) ZGFinished() {

	result := make(map[string]interface{})

	result["code"] = 0
	result["message"] = "操作成功"

	defer func() {
		c.Data["json"] = result
		c.ServeJSON()
	}()

	s := strings.ToUpper(c.GetString("symbol"))
	tempSymbol := strings.Split(s, "-")
	symbol := tempSymbol[0] + "_" + "CNT"

	db := utils.RobotDB

	var tradeResult []models.ZGTradeResults
	db.Model(&models.ZGTradeResults{}).Where(&models.ZGTradeResults{Symbol: symbol}).Find(&tradeResult)

	result["results"] = tradeResult
}


