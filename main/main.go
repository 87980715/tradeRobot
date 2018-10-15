package main

import (
	"github.com/astaxie/beego/logs"
	"tradeRobot/initialize"
	"tradeRobot/dataAgent"
	"tradeRobot/trade"
)

func main() {

	err := initialize.InitRobot()
	if err != nil {
		logs.Error("init robot failed err:",err)
		panic("init robot failed")
	}

	logs.Debug("init conf succ, config")

	initialize.AppConfig.Wg.Add(1)

	var coinNames = map[int]string{1: "ltc",2:"true"}//,3:"eos"}
	var toCoinNames = map[int]string{1: "eth", 2: "btc",3:"usdt"}// 3: "usdt"}

	api_keys := []string{trade.API_KEY}

	accounts := trade.SetAccounts(api_keys)

	dataAgent.AgentServerRun(coinNames,toCoinNames)

	trade.TradeServerRun(coinNames,toCoinNames,accounts)

	initialize.AppConfig.Wg.Wait()

}
