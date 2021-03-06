package trade

import (
	"time"
	"tradeRobot/robot/utils"
	"fmt"
	"tradeRobot/robot/models"
	"context"
	"strings"
	"strconv"
	"sync"
	"github.com/astaxie/beego/logs"
	"math/rand"
)

var RMuLock sync.RWMutex
// 限价交易
func TradeLimitHuobi(ctx context.Context) {
	rand.Seed(time.Now().UnixNano())
	for {
		time.Sleep(151 * time.Millisecond)
		select {
		case <-ctx.Done():
			return
		default:
			postDataLimit := <-utils.HuobiOrders
			// fmt.Println("amount:",postDataLimit.Amount)
			// 设置交易参数,精度需处理
			account := utils.HuobiAccount
			account.PostDataLimit.Symbol = postDataLimit.Symbol
			account.PostDataLimit.Type = postDataLimit.Type
			account.PostDataLimit.Account_id = strconv.FormatInt(models.UserID[postDataLimit.Symbol],10)

			RMuLock.RLock()
			ethPrice := models.EthPrice["huobi"]
			usdtPrice := models.UsdtPrice["huobi"]
			RMuLock.RUnlock()

			for ethPrice*usdtPrice == 0 {
				RMuLock.RLock()
				ethPrice = models.EthPrice["huobi"]
				usdtPrice = models.UsdtPrice["huobi"]
				RMuLock.RUnlock()
				if ethPrice*usdtPrice != 0 {
					break
				}
			}
			switch postDataLimit.Symbol {
			case "mteth":
				// 数量最小为0.01，2位小数
				amount, _ := strconv.ParseFloat(postDataLimit.Amount, 64)
				account.PostDataLimit.Amount = fmt.Sprintf("%."+strconv.Itoa(2)+"f", amount)
				// fmt.Println("amount:",account.PostDataLimit.Amount)
				// 价格，8位小数
				p, _ := strconv.ParseFloat(postDataLimit.Price, 64)
				price := p / (ethPrice * usdtPrice)
				// fmt.Println("price:",price)
				account.PostDataLimit.Price = fmt.Sprintf("%."+strconv.Itoa(8)+"f", price)
			case "aeeth":
				amount, _ := strconv.ParseFloat(postDataLimit.Amount, 64)
				account.PostDataLimit.Amount = fmt.Sprintf("%."+strconv.Itoa(4)+"f", amount)
				// 价格，6位小数
				p, _ := strconv.ParseFloat(postDataLimit.Price, 64)
				price := p / (ethPrice * usdtPrice)
				account.PostDataLimit.Price = fmt.Sprintf("%."+strconv.Itoa(6)+"f", price)
			case "ethusdt","btcusdt":
				amount,_:= strconv.ParseFloat(postDataLimit.Amount, 64)
				account.PostDataLimit.Amount = fmt.Sprintf("%."+strconv.Itoa(4)+"f", amount)
				// 价格，2位小数
				p,_:= strconv.ParseFloat(postDataLimit.Price, 64)
				price := p / usdtPrice
				account.PostDataLimit.Price = fmt.Sprintf("%."+strconv.Itoa(2)+"f", price)
			default:
				amount,_:= strconv.ParseFloat(postDataLimit.Amount, 64)
				account.PostDataLimit.Amount = fmt.Sprintf("%."+strconv.Itoa(4)+"f", amount)
				// 价格，2位小数
				p,_:= strconv.ParseFloat(postDataLimit.Price, 64)
				price := p / usdtPrice
				account.PostDataLimit.Price = fmt.Sprintf("%."+strconv.Itoa(2)+"f", price)
			}
			// 执行交易
			account.HuobiLimitTrade()
			logs.Info("%s 火币挂单价格：%s  数量：%s\n", account.PostDataLimit.Symbol, account.PostDataLimit.Price, account.PostDataLimit.Amount)
		}
	}
}

func TradeCancelHuobi(symbol []string, ctx context.Context) {
	for {
		time.Sleep(500 * time.Millisecond)
		select {
		case <-ctx.Done():
			return
		default:
			s := strings.ToLower(symbol [0] + symbol[2])
			acount := utils.HuobiAccount
			acount.GetDataPending.Symbol = s
			acount.GetDataPending.Account_id = strconv.FormatInt(models.UserID[s],10)
			acount.GetDataPending.Size = models.Huobi_PendingOrdersSize
			acount.HuobiCancelPendingOrders()
		}
	}
}

// 将已成交的交易插入数据库
func HuobiInsertToDB(symbol []string, ctx context.Context) {

	db := utils.RobotDB
	var preTradeResult models.HuobiTradeResults
	s := strings.ToLower(symbol[0] + symbol[2])
	db.Model(&models.HuobiTradeResults{}).Where("symbol = ?", s).Order("trade_id desc").Limit(1).Find(&preTradeResult)
	for {
		time.Sleep(300 * time.Millisecond)
		select {
		case <-ctx.Done():
			return
		default:
			acount := utils.HuobiAccount
			acount.GetTradesDeal.Symbol = s
			acount.HuobiTradesDeal(preTradeResult)
		}
	}
}

/*
func TradeLimitHuobi(ctx context.Context) {
	for {
		time.Sleep(150 * time.Millisecond)
		select {
		case <-ctx.Done():
			return
		default:
			postDataLimit := <-utils.HuobiOrders
			// 设置交易参数,精度需处理
			account := utils.HuobiAccount
			account.PostDataLimit.Symbol = postDataLimit.Symbol
			// 数量最小为0.001，四位小数
			amount,_:= strconv.ParseFloat(postDataLimit.Amount, 64)
			a := strconv.FormatFloat(amount, 'E', -1, 64)
			account.PostDataLimit.Amount = a
			// 价格，两位小数
			price,_:= strconv.ParseFloat(postDataLimit.Price, 64)
			account.PostDataLimit.Price = fmt.Sprintf("%."+strconv.Itoa(2)+"f", price)

			account.PostDataLimit.Type = postDataLimit.Type
			account.PostDataLimit.Account_id = models.Huobi_Account_ID
			// 执行交易
			//account.HuobiLimitTrade()
			logs.Info("%s 火币挂单价格：%s  数量：%s\n", account.PostDataLimit.Symbol, account.PostDataLimit.Price, account.PostDataLimit.Amount)
		}
	}
}
 */
