package initialize

import (
	"tradeRobot/utils"
	"tradeRobot/models"
)

func initAccounts() {
	utils.ZTAccounts = utils.SetAccountsZT(models.ZTApiKeys)
	utils.OKexAccount = utils.SetAccountOKex(models.OKex_API_KEY)
}