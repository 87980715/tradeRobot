package main

import "tradeRobot/robot/utils"

func main() {

	account := utils.SetAccountsZT("ZBC7CAodYDdHYxSI1h8xLqyqTxnV3dt92dpr")
	account.API_KEY = "ZBC7CAodYDdHYxSI1h8xLqyqTxnV3dt92dpr"
	account.SECRET_KEY = "8XUkytGxGH2ctvTHhrylBzho1lEz6YrS7AM8"
	account.ZTQueyMd5Sign()

	account.ZTGetUserAssets()


}