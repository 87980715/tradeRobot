package trade

import (
	"github.com/astaxie/beego/logs"
	"time"
	"fmt"
	"strconv"
	"tradeRobot/models"
	"math/rand"
	"strings"
	"tradeRobot/dataAgent"
	"encoding/json"
	"tradeRobot/utils"
)

// 账户
// 1. 创建交易账户，分为限价交易账户和市价交易账户
// 2. SECRET_KEY，API_KEY

//var TickersChan = make(chan map[string]string,100)

var DepthsChan = make(chan map[string]string,100)

// 限价交易
func TradeLimitZT() {
	Loop:
	for {
		temMap := <-DepthsChan
		for _, tempSmybol := range models.AllSymbols {
			if depths, ok := temMap[tempSmybol]; ok {
				symbolStrs := strings.Split(tempSmybol, "_")
				symbol := strings.ToUpper(symbolStrs[0]) + "_" + "CNZ"
				curDepth := &dataAgent.OKexDepth{}
				err := json.Unmarshal([]byte(depths), curDepth)
				if err != nil {
					logs.Error("json.Unmarshal depth failed err:", err)
					continue Loop
				}
				//随机获取指定数量的限价买卖账户
				RandBuyAccounts, RandSoldAccounts := getAccounts(utils.ZTAccounts)

				//根据账户的数量，随机获取相同数量的的Asks，Bids
				curDepth.Asks, curDepth.Bids = GetRandArray(curDepth.Asks, curDepth.Bids, len(RandBuyAccounts))

				//账户进行限价买入
				for index, account := range RandBuyAccounts {
					account.PostDataLimit.Market = symbol
					account.PostDataLimit.Side = "1"

					// 设置价格
					price := curDepth.Bids[index][0] * 6.88
					account.PostDataLimit.Price = fmt.Sprintf("%."+strconv.Itoa(4)+"f", price)

					// 成交量设置
					amount := curDepth.Bids[index][1]
					account.PostDataLimit.Amount = fmt.Sprintf("%."+strconv.Itoa(4)+"f", amount)
					// 签名
					account.ZTLimitMd5Sign()
					fmt.Printf("%s买入价格：%s  数量：%s\n", account.PostDataLimit.Market, account.PostDataLimit.Price,account.PostDataLimit.Amount)
					//fmt.Println("account.Sign_trade_limit_buy:", account.Sign)
					account.ZTTradeLimit()
				}
				//账户进行限价卖出
				for index, account := range RandSoldAccounts {
					account.PostDataLimit.Market = symbol
					account.PostDataLimit.Side = "2"
					price := curDepth.Asks[index][0] * 6.88
					account.PostDataLimit.Price = fmt.Sprintf("%."+strconv.Itoa(4)+"f", price)

					amount := curDepth.Asks[index][1]
					account.PostDataLimit.Amount = fmt.Sprintf("%."+strconv.Itoa(4)+"f", amount)
					account.ZTLimitMd5Sign()
					fmt.Printf("%s卖出价格：%s  数量：%s\n", account.PostDataLimit.Market, account.PostDataLimit.Price,account.PostDataLimit.Amount)
					fmt.Println("---------------------------")
					account.ZTTradeLimit()
				}
			}
		}
	}
}

// 市价交易
func TradeMarketZT() {
	var curDepth = &dataAgent.OKexDepth{}
Loop:
	for {
		temMap := <-DepthsChan
		for _, tempsmybol := range models.AllSymbols {
			if depths, ok := temMap[tempsmybol]; ok {
				symbolStrs := strings.Split(tempsmybol, "_")
				symbol := strings.ToUpper(symbolStrs[0]) + "_" + "CNZ"
				err := json.Unmarshal([]byte(depths), curDepth)
				if err != nil {
					logs.Error("json.Unmarshal depth failed err:", err)
					continue Loop
				}
				// 随机获取指定数量的市价买卖账户
				RandBuyAccounts, RandSoldAccounts := getAccounts(utils.ZTAccounts)

				// 根据账户的数量，随机获取相同数量的的Asks，Bids
				curDepth.Asks, curDepth.Bids = GetRandArray(curDepth.Asks, curDepth.Bids, len(RandSoldAccounts))

				// 买
				for index, account := range RandBuyAccounts {
					account.PostDataMarket.Market = symbol // key 就是交易对
					account.PostDataMarket.Side = "1"
					// 设置数量
					amount := curDepth.Asks[index][1]
					account.PostDataLimit.Amount = fmt.Sprintf("%."+strconv.Itoa(4)+"f", amount)
					account.ZTMarketMd5Sign()
					//fmt.Println("account.Sign_trade_market_sold:",account.Sign)
					account.ZTTradeMarket()
				}
				// 卖
				for index, account := range RandSoldAccounts {
					account.PostDataMarket.Market = symbol // key 就是交易对
					account.PostDataMarket.Side = "2"
					// 设置数量
					amount := curDepth.Bids[index][1]
					account.PostDataLimit.Amount = fmt.Sprintf("%."+strconv.Itoa(4)+"f", amount)
					account.ZTMarketMd5Sign()
					account.ZTTradeMarket()
				}
			}
		}
	}
}

// 获取随机的跟账户相同数量的asks和bids,用于下单
func GetRandArray(rawAsks, rawBids [][2]float64, arrayLen int) (asks, bids [][2]float64) {

	rand.Seed(time.Now().UnixNano())
	randArray := make([]int, 0)

	for i := 0; i < arrayLen; i++ {
		num := rand.Intn(10)
		randArray = append(randArray, num)
	}

	for _, v := range randArray {
		asks = append(asks, rawAsks[v])
		bids = append(bids, rawBids[v])
	}
	return
}

// 在提供的账户中，随机挑选指定数量的账户
func getAccounts(rawAccounts []*utils.ZTAccount) (RandBuyAccounts, RandSoldAccounts []*utils.ZTAccount) {
	rand.Seed(time.Now().UnixNano())
	randArray1 := make([]int, 0)
	randArray2 := make([]int, 0)

	for i := 0; i < len(rawAccounts); i++ {
		num1 := rand.Intn(len(rawAccounts))
		randArray1 = append(randArray1, num1)
	}

	for i := 0; i < len(rawAccounts); i++ {
		num2 := rand.Intn(len(rawAccounts))
		randArray2 = append(randArray2, num2)
	}

	for _, v := range randArray1 {
		RandBuyAccounts = append(RandBuyAccounts, rawAccounts[v])
	}

	for _, v := range randArray1 {
		RandSoldAccounts = append(RandSoldAccounts, rawAccounts[v])
	}
	return
}
