package main

import (
	"github.com/astaxie/beego/logs"
	"tradeRobot/initialize"
	"tradeRobot/dataAgent"
	"tradeRobot/models"
	"tradeRobot/trade"
)

func main() {

	models.OKexSymbolsArray = [][2]string {{"eth","usdt"}} // ,{"eos","usdt"}
	models.HuobiSymbolsArray = [][2]string{} // {"btc","usdt"},{"ltc","usdt"}

	err := initialize.InitRobot()
	if err != nil {
		logs.Error("init robot failed err:",err)
		panic("init robot failed")
	}

	logs.Debug("init conf succ")
	initialize.AppConfig.Wg.Add(1)
	dataAgent.AgentServerRun()
	if err != nil {
		logs.Error("dataAgent.AgentServerRun failed err:",err)
		panic(err)
	}

	trade.ZTTradeServerRun()

	initialize.AppConfig.Wg.Wait()
}
