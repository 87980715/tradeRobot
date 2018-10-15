package trade

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"math/rand"
	"time"
	"fmt"
	"strconv"
)

// 账户
// 1. 创建交易账户，分为限价交易账户和市价交易账户
// 2. SECRET_KEY，API_KEY


// 限价交易
func TradeLimitZT(coinNames map[int]string, toCionNames map[int]string, accounts []*Account) {
	curDepth := &OKexDepth{}
Loop:
	for {
		temMap := <-DepthsChan
		for _, toCionName := range toCionNames {
			for _, cionName := range coinNames {
				key := cionName + "_" + toCionName
				if depths, ok := temMap[key]; ok {
					// 根据coinName_toCoinName处理获取的depths数据
					err := json.Unmarshal([]byte(depths), curDepth)
					if err != nil {
						logs.Error("json.Unmarshal depthDoc.Text() failed:", err)
						continue Loop
					}
					//随机获取指定数量的限价买卖账户
					RandBuyAccounts,RandSoldAccounts := getAccounts(accounts)

					//asks,bids 各截取15个数据
					curDepth.Asks = curDepth.Asks[len(curDepth.Asks)-15:len(curDepth.Asks)]
					curDepth.Bids = curDepth.Bids[:15]

					//根据账户的数量，随机获取相同数量的的Asks，Bids
					curDepth.Asks, curDepth.Bids = GetRandArray(curDepth.Asks, curDepth.Bids,len(RandBuyAccounts))

					//账户进行限价买入
					for index,account := range RandBuyAccounts {
						account.ResPostDataLimit.Market = key
						account.ResPostDataLimit.Side = "1"
						account.ResPostDataLimit.Price = strconv.FormatFloat(curDepth.Asks[index][0], 'E', -1, 64)

						// 成交量设置
						account.ResPostDataLimit.Amount = strconv.FormatFloat(curDepth.Asks[index][1], 'E', -1, 64)
						account.Md5Sign()
						fmt.Printf("%s买入价格：%s\n",account.ResPostDataLimit.Market,account.ResPostDataLimit.Price)
						fmt.Println("account.Sign_trade_limit_buy:",account.Sign)
						//account.TradeLimit()
					}
					//账户进行限价卖出
					for index,account := range RandSoldAccounts {
						account.ResPostDataLimit.Market = key
						account.ResPostDataLimit.Side = "2"
						account.ResPostDataLimit.Price = strconv.FormatFloat(curDepth.Bids[index][0], 'E', -1, 64)

						account.ResPostDataLimit.Amount = strconv.FormatFloat(curDepth.Bids[index][1], 'E', -1, 64)
						account.Md5Sign()
						fmt.Printf("%s卖出价格：%s\n",account.ResPostDataLimit.Market,account.ResPostDataLimit.Price)
						fmt.Println("account.Sign_trade_limit_sold:",account.Sign)
						//account.TradeLimit()
					}
				}
			}
		}
	}
}


// 市价交易
func TradeMarketZT(coinNames map[int]string, toCionNames map[int]string, accounts []*Account) {
	curDepth := &OKexDepth{}
Loop:
	for {
		temMap := <-DepthsChan
		for _, toCionName := range toCionNames {
			for _, cionName := range coinNames {
				key := cionName + "_" + toCionName
				if depths, ok := temMap[key]; ok {
					err := json.Unmarshal([]byte(depths), curDepth)
					if err != nil {
						logs.Error("json.Unmarshal depthDoc.Text() failed:", err)
						continue Loop
					}

					// 随机获取指定数量的市价买卖账户
					RandBuyAccounts,RandSoldAccounts := getAccounts(accounts)

					// asks,bids 各截取15个数据
					curDepth.Asks = curDepth.Asks[len(curDepth.Asks)-15:len(curDepth.Asks)]
					curDepth.Bids = curDepth.Bids[:15]

					// 根据账户的数量，随机获取相同数量的的Asks，Bids
					curDepth.Asks, curDepth.Bids = GetRandArray(curDepth.Asks, curDepth.Bids,len(RandSoldAccounts))

					// 买
					for _,account := range RandBuyAccounts {
						account.ResPostDataMarket.Market = key // key 就是交易对
						account.ResPostDataMarket.Side = "1"

						account.Md5Sign()
						//fmt.Println("account.Sign_trade_market_sold:",account.Sign)
						//account.TradeMarket()
					}
					// 卖
					for _,account := range RandSoldAccounts {
						account.ResPostDataMarket.Market = key // key 就是交易对
						account.ResPostDataMarket.Side = "2"

						account.Md5Sign()
						//fmt.Println("account.Sign_trade_market_sold:",account.Sign)
						//account.TradeMarket()
					}
				}
			}
		}
	}
}

// 获取随机的跟账户相同数量的asks和bids,用于下单
func GetRandArray(rawAsks, rawBids [][2]float64,arrayLen int) (asks, bids [][2]float64) {

	rand.Seed(time.Now().UnixNano())
	randArray := make([]int, 0)

	for i := 0; i < arrayLen; i++ {
		num := rand.Intn(arrayLen)
		randArray = append(randArray, num)
	}

	for _, v := range randArray {
		asks = append(asks, rawAsks[v])
		bids = append(bids, rawBids[v])
	}
	return
}

// 在提供的账户中，随机挑选指定数量的账户
func getAccounts(rawAccounts []*Account) (RandBuyAccounts,RandSoldAccounts []*Account) {
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
