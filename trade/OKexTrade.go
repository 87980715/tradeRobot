package trade

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"tradeRobot/utils"
)

var OrdersChan = make(chan string, 100)

// 限价交易
func TradeLimitOKex() {

Loop:
	for {
		data := <-OrdersChan
		limitData := &utils.OKexPostDataLimit{}
		err := json.Unmarshal([]byte(data), limitData)
		if err != nil {
			logs.Error("json.Unmarshal depth failed err:", err)
			continue Loop
		}
		//账户进行限价交易
		utils.OKexAccount.PostDataLimit.Symbol = limitData.Symbol
		utils.OKexAccount.PostDataLimit.Type = limitData.Type
		utils.OKexAccount.PostDataLimit.Amount = limitData.Amount
		utils.OKexAccount.PostDataLimit.Price = limitData.Price
		/*
		// 签名
		utils.OKexAccount.OKexLimitMd5Sign()
		fmt.Printf("%s买入价格：%s\n  数量：%s", account.PostDataLimit.Market, account.PostDataLimit.Price, account.PostDataLimit.Amount)
		//fmt.Println("account.Sign_trade_limit_buy:", account.Sign)
		utils.OKexAccount.OKexTradeLimit()
		*/
	}

}
