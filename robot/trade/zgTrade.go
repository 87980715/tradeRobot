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
	rand.Seed(time.Now().UnixNano())
	postDataLimit := &utils.ZTPostDataLimit{}
Loop:
	for {
		time.Sleep(500 * time.Millisecond)
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
			//账户进行限价交易
			account := utils.ZTAccount
			account.PostDataLimit.Market = postDataLimit.Market
			account.PostDataLimit.Side = postDataLimit.Side
			switch postDataLimit.Market {
			case "BTC_CNT":
				p, _ := strconv.ParseFloat(postDataLimit.Price, 64)
				account.PostDataLimit.Price = fmt.Sprintf("%."+strconv.Itoa(4)+"f", p)
				// 交易数量Amount 设置
				a, _ := strconv.ParseFloat(postDataLimit.Amount, 64)
				amount := a * models.TradeAmountMultiple
				account.PostDataLimit.Amount = fmt.Sprintf("%."+strconv.Itoa(4)+"f", amount)
			case "AE_CNT":
				p, _ := strconv.ParseFloat(postDataLimit.Price, 64)
				account.PostDataLimit.Price = fmt.Sprintf("%."+strconv.Itoa(4)+"f", p)
				// 交易数量Amount 设置
				a, _ := strconv.ParseFloat(postDataLimit.Amount, 64)
				amount := a * models.TradeAmountMultiple
				account.PostDataLimit.Amount = fmt.Sprintf("%."+strconv.Itoa(4)+"f", amount)
			default:
				p, _ := strconv.ParseFloat(postDataLimit.Price, 64)
				account.PostDataLimit.Price = fmt.Sprintf("%."+strconv.Itoa(4)+"f", p)
				// 交易数量Amount 设置
				a, _ := strconv.ParseFloat(postDataLimit.Amount, 64)
				amount := a * models.TradeAmountMultiple
				account.PostDataLimit.Amount = fmt.Sprintf("%."+strconv.Itoa(4)+"f", amount)
			}
			// 签名
			account.ZTLimitMd5Sign()
			account.ZTTradeLimit()
			logs.Info("%s ZT挂单价格：%s  数量：%s\n", account.PostDataLimit.Market, account.PostDataLimit.Price, account.PostDataLimit.Amount)
		}
	}
}

// 取消的ZTorders
func CanleOrdersZG(symbol []string, ctx context.Context) {
	rand.Seed(time.Now().Unix())
	market := symbol[0] + "_" + symbol[1]
	for {
		time.Sleep(4000 * time.Millisecond)
		select {
		case <-ctx.Done():
			return
		default:
			account := utils.ZTAccount
			account.PostDataQueryPending.Limit = 30
			account.PostDataQueryPending.Market = market
			account.PostDataQueryPending.Offset = 150 // 300
			account.ZTQueryPendingMd5Sign()
			postDatas := account.ZTQueryPending()
			// 分开处理
			if postDatas != nil {
				for _, postData := range postDatas {
					account.PostDataCancel.Market = postData.Market
					account.PostDataCancel.Order_id = postData.Order_id
					account.ZTCancelMd5Sign()
					// 取消之后重新挂单
					if account.ZTCancelOrder() {
						rand.Seed(time.Now().Unix())
						if rand.Intn(20)%2 == 0 {
							var postDataLimit = &utils.ZTPostDataLimit{}
							postDataLimit.Market = postData.Market
							postDataLimit.Amount = postData.Amount
							if postData.Side == 1 {
								postDataLimit.Side = "1"
								p, _ := strconv.ParseFloat(postData.Price, 64)
								price := p * (1 - 0.002)
								postDataLimit.Price = fmt.Sprintf("%."+strconv.Itoa(4)+"f", price)
							} else {
								postDataLimit.Side = "2"
								p, _ := strconv.ParseFloat(postData.Price, 64)
								price := p * (1 + 0.002)
								postDataLimit.Price = fmt.Sprintf("%."+strconv.Itoa(4)+"f", price)
							}
							data, _ := json.Marshal(postDataLimit)
							models.ZGTradesChan <- string(data)
						}
					}
				}
			}
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