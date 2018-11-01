package trade

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"tradeRobot/robot/utils"
)

var OrdersChan = make(chan string, 100)

// 限价交易
func TradeLimitOKex() {
Loop:
	for {
		order := <-OrdersChan
		limitData := &utils.OKexPostDataLimit{}
		err := json.Unmarshal([]byte(order), limitData)
		if err != nil {
			logs.Error("json.Unmarshal depth failed err:", err)
			continue Loop
		}
		//账户进行限价交易
		utils.OKexAccount.PostDataLimit.Instrument_id = limitData.Instrument_id
		utils.OKexAccount.PostDataLimit.Type = limitData.Type
		utils.OKexAccount.PostDataLimit.Amount = limitData.Amount
		utils.OKexAccount.PostDataLimit.Price = limitData.Price

		// 签名
		utils.OKexAccount.OKexLimitMd5Sign()
		//fmt.Println("account.Sign_trade_limit_buy:", account.Sign)
		/*
		utils.OKexAccount.OKexTradeLimit()
		*/
	}
}
