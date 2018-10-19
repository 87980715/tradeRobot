package utils

import (
	"tradeRobot/models"
)

// 创建交易账户，设置交易需要传递的参数
type ZTAccount struct {
	ZTRestfulApiRequest
}

var ZTAccounts []*ZTAccount
var OKexAccount *OKexRestfulApiRequest

func SetAccountsZT(apiKeys []string) []*ZTAccount {
	tempAccount := &ZTAccount{
		ZTRestfulApiRequest{
			PostDataLimit:  &ZTPostDataLimit{},
			PostDataMarket: &ZTPostDataMarket{},
		},
	}
	accounts := make([]*ZTAccount, 0)
	for _, v := range apiKeys {
		tempAccount.API_KEY = v
		tempAccount.SECRET_KEY = models.ZG_SECRET_KEY

		accounts = append(accounts, tempAccount)
	}
	return accounts
}

func SetAccountOKex(apiKey string) *OKexRestfulApiRequest {
	tempAccount := &OKexRestfulApiRequest{
			PostDataLimit:  &OKexPostDataLimit{},
			PostDataMarket: &OKexPostDataMarket{},
			PostPataCancle: &OKexPostDataCancel{},
	}
	tempAccount.API_KEY = apiKey
	tempAccount.SECRET_KEY = models.OKex_SECRET_KEY

	return tempAccount
}
