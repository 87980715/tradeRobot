package main

import (
	"tradeRobot/robot/utils"
	"tradeRobot/robot/models"
)

func main() {
	acount := utils.SetAccountHuobi()
	acount.GetDataPending.Symbol = "mteth"
	acount.GetDataPending.Account_id = "4821321"
	acount.GetDataPending.Size = models.Huobi_PendingOrdersSize
	acount.HuobiCancelPendingOrders()
}