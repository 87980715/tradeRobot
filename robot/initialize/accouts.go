package initialize

import (
	"tradeRobot/robot/utils"
	"tradeRobot/robot/models"
)

func initAccounts() {
	utils.ZTAccount = utils.SetAccountsZT(models.ZT_API_KEY)
	utils.OKexAccount = utils.SetAccountOKex(models.OKex_API_KEY)
	utils.HuobiAccount = utils.SetAccountHuobi()
}