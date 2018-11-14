package main

import (
	"tradeRobot/robot/utils"
)

func main() {
	account := utils.SetAccountHuobi()

	account.PostDataLimit.Symbol ="mteth"
	account.PostDataLimit.Amount = "50000"
	account.PostDataLimit.Price = "0.00001850"
	account.PostDataLimit.Type = "buy-limit"
	account.PostDataLimit.Account_id = "4821321"
	// 执行交易
	account.HuobiLimitTrade()

}

