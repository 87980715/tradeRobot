package trade

import (
	"time"
	"tradeRobot/robot/utils"
	"fmt"
	"tradeRobot/robot/models"
	"context"
)

// 限价交易
func TradeLimitHuobi(ctx context.Context) {
	for {
		time.Sleep(150 * time.Millisecond)
		select {
		case <-ctx.Done():
			return
		default:
			postDataLimit := <-utils.HuobiOrders
			// 设置交易参数
			account := utils.HuobiAccount
			account.PostDataLimit.Symbol = postDataLimit.Symbol
			account.PostDataLimit.Amount = postDataLimit.Amount
			account.PostDataLimit.Price = postDataLimit.Price
			account.PostDataLimit.Type = postDataLimit.Type
			account.PostDataLimit.Account_id = models.Huobi_Account_ID
			// 执行交易
			account.HuobiLimitTrade()
			fmt.Printf("%s 成交价格：%s  数量：%s\n", account.PostDataLimit.Symbol, account.PostDataLimit.Price, account.PostDataLimit.Amount)
		}
	}
}

func TradeCancelHuobi(symbol []string, ctx context.Context) {
	for {
		time.Sleep(100 * time.Millisecond)
		select {
		case <-ctx.Done():
			return
		default:
			s := symbol [0] + symbol [1]
			acount := utils.HuobiAccount
			acount.GetDataPending.Symbol = s
			acount.GetDataPending.Account_id = models.Huobi_Account_ID
			acount.GetDataPending.Size = models.Huobi_PendingOrdersSize
			acount.HuobiCancelPendingOrders()
		}
	}
}

// 将已成交的交易插入数据库
func HuobiInsertToDB(ctx context.Context) {
	for {
		time.Sleep(100 * time.Millisecond)
		select {
		case <-ctx.Done():
			return
		default:

			acount := utils.HuobiAccount
			acount.HuobiTradesDeal()
		}
	}
}
