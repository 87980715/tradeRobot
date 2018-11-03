package main

import (
	"tradeRobot/robot/utils"
	)

func main() {

	acount := utils.SetAccountHuobi()

	acount.GetTradesDeal.Symbol = "ethusdt"

	acount.HuobiTradesDeal()
}