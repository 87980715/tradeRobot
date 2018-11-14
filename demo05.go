package main

import (
	"tradeRobot/robot/utils"
	"github.com/astaxie/beego/logs"
	"tradeRobot/robot/models"
	"fmt"
)

func main() {

	db,err := utils.LoadRobotDB()
	if err != nil {
		logs.Error("load robotDB failed err:",err)
	}

	var tradeResult []models.ZGTradeResults
	db.Model(&models.ZGTradeResults{}).Where(&models.ZGTradeResults{Symbol:"ETH_CNT"}).Find(&tradeResult)
	db.Close()

	fmt.Println(tradeResult)
}
