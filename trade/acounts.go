package trade

import (
	"tradeRobot/utils"
)

// 创建交易账户，设置交易需要传递的参数
type Account struct {
	utils.RestfulApiRequest
}

const API_KEY = "ed5d8197-26db-45be-b1ce-719f13847b6c"

func SetAccounts(apiKeys []string) []*Account {

	tempAccount := &Account{
		utils.RestfulApiRequest{
			ResPostDataLimit:&utils.PostDataLimit{},
			ResPostDataMarket:&utils.PostDataMarket{},
		},
	}
	//fmt.Println("tempAccount==nil:",tempAccount==nil)
	accounts := make([]*Account,0)

	for _,v := range apiKeys {
		tempAccount.API_KEY = v
		tempAccount.SECRET_KEY = SECRET_KEY

		accounts = append(accounts,tempAccount)
	}
	/*
	fmt.Println("accounts:",accounts)
	fmt.Println("API_KEY:",accounts[0].API_KEY)
	fmt.Println("API_KEY:",accounts[0].SECRET_KEY)
	*/
	return accounts
}
