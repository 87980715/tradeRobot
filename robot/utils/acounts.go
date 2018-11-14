package utils

import (
	"tradeRobot/robot/models"
)

// 创建交易账户，设置交易需要传递的参数


var ZTAccount *ZTRestfulApiRequest
var OKexAccount *OKexRestfulApiRequest
var HuobiAccount *HuobiRestfulApiRequest

func SetAccountsZT(apiKey string) *ZTRestfulApiRequest {
	account := &ZTRestfulApiRequest{
		PostDataLimit:  &ZTPostDataLimit{},
		PostDataMarket: &ZTPostDataMarket{},
		PostDataQueryPending : &ZTPostDataQureyPending{},
		PostDataCancel:&ZTPostDataCancel{},
		PostDataOrderFinished:&ZTPostDataOrderFinished{},
	}
	account.API_KEY = apiKey
	account.SECRET_KEY = models.ZT_SECRET_KEY
	return account
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

func SetAccountHuobi() *HuobiRestfulApiRequest {
	tempAccount := &HuobiRestfulApiRequest{
		PostDataLimit : &HuobiPostDataLimit{},
		PostDataMarket: &HuobiPostDataMarket{},
		PostPataCancle :&HuobiPostDataCancel{},
		GetDataPending: &HuobiGetDataPending{},
		GetTradesDeal: &HuobiGetTradesDeal {},
	}
	return tempAccount
}