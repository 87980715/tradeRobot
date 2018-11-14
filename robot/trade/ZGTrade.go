package trade

import (
	"github.com/astaxie/beego/logs"
	"fmt"
	"strconv"
	"strings"
	"encoding/json"
	"tradeRobot/robot/utils"
	"tradeRobot/robot/models"
	"time"
	"math/rand"
	"context"
)

// 限价交易
func TradeLimitZG(ctx context.Context) {
	postDataLimit := &utils.ZTPostDataLimit{}
Loop:
	for {
		time.Sleep(150 * time.Millisecond)
		select {
		case <-ctx.Done():
			return
		default:
			trade := <-models.ZGTradesChan
			err := json.Unmarshal([]byte(trade), postDataLimit)
			if err != nil {
				logs.Error("json.Unmarshal trade failed err:", err)
				continue Loop
			}
			symbolStrs := strings.Split(postDataLimit.Market, "_")
			symbol := strings.ToUpper(symbolStrs[0]) + "_" + "CNT"
			//账户进行限价交易
			account := utils.ZTAccount
			account.PostDataLimit.Market = symbol
			account.PostDataLimit.Side = postDataLimit.Side
			p, _ := strconv.ParseFloat(postDataLimit.Price, 64)
			account.PostDataLimit.Price = fmt.Sprintf("%."+strconv.Itoa(4)+"f", p)
			// 交易数量Amount 设置
			a, _ := strconv.ParseFloat(postDataLimit.Amount, 64)
			account.PostDataLimit.Amount = fmt.Sprintf("%."+strconv.Itoa(4)+"f", a*models.TradeAmountMultiple)
			// 签名
			account.ZTLimitMd5Sign()
			account.ZTTradeLimit()
			logs.Info("%s ZT挂单价格：%s  数量：%s\n", account.PostDataLimit.Market, account.PostDataLimit.Price, account.PostDataLimit.Amount)
		}
	}
}

// 市价交易
func TradeMarketZG() {
	postDataLimit := &utils.ZTPostDataLimit{}
Loop:
	for {
		trade := <-models.ZGTradesChan
		err := json.Unmarshal([]byte(trade), postDataLimit)
		if err != nil {
			logs.Error("json.Unmarshal depth failed err:", err)
			continue Loop
		}
		symbolStrs := strings.Split(postDataLimit.Market, "_")
		symbol := strings.ToUpper(symbolStrs[0]) + "_" + "CNT"

		// 交易
		account := utils.ZTAccount
		account.PostDataMarket.Market = symbol // key 就是交易对
		account.PostDataMarket.Side = postDataLimit.Side
		// 设置数量
		amount := postDataLimit.Amount
		account.PostDataLimit.Amount = fmt.Sprintf("%."+strconv.Itoa(4)+"f", amount)
		account.ZTMarketMd5Sign()
		account.ZTTradeMarket()
	}
}

// 取消的ZTorders
func CanleOrdersZG(symbol []string,ctx context.Context) {
	rand.Seed(time.Now().Unix())
	market := symbol[0] + "_" + symbol[1]
	for {
		time.Sleep(4000 * time.Millisecond)
		select {
		case <-ctx.Done():
			return
		default:
			account := utils.ZTAccount
			account.PostDataQueryPending.Limit = 20
			account.PostDataQueryPending.Market = market
			account.PostDataQueryPending.Offset = 200 // 300
			account.ZTQueryPendingMd5Sign()
			postDatas := account.ZTQueryPending()
			// 分开处理
			if postDatas != nil {
				for _, postData := range postDatas {
					account.PostDataCancel.Market = postData.Market
					account.PostDataCancel.Order_id = postData.Order_id
					account.ZTCancelMd5Sign()
					account.ZTCancelOrder()
				}
			}
		}
	}
}