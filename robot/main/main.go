package main

import (
	"github.com/astaxie/beego/logs"
	"tradeRobot/robot/models"
	"tradeRobot/robot/initialize"
	"tradeRobot/robot/trade"
	"tradeRobot/robot/dataAgent"
)

func main() {
	symbol := []string{"eth","usdt"}

	models.TradeInspectTime = 500  //单位毫秒
	models.TradePriceAdjust = 0.001
	models.TradeAmountMultiple = 1

	err := initialize.InitRobot()
	if err != nil {
		logs.Error("init robot failed err:",err)
		panic("init robot failed")
	}

	logs.Debug("init conf succ")

	initialize.AppConfig.Wg.Add(1)

	dataAgent.AgentServerRun(symbol)

	trade.TradeServerRun(symbol)

	initialize.AppConfig.Wg.Wait()

}