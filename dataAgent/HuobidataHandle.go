package dataAgent

import (
	"net/http"
	"github.com/astaxie/beego/logs"
	"github.com/PuerkitoBio/goquery"
	"encoding/json"
	"time"
	"strconv"
	"tradeRobot/utils"
	"tradeRobot/models"
	"math/rand"
)

type HuobiDepthRes struct {
	Status string      `json:"status"`
	Ch     string      `json:"ch"`
	Ts     float64     `json:"ts"`
	Tick   *HuobiDepth `json:"tick"`
}

type HuobiDepth struct {
	Bids [][2]float64 `json:"bids"`
	Asks [][2]float64 `json:"asks"`
	//Ts      float64      `json:"ts"`
	//Version float64      `json:"version"`
}

type HuobiTrades struct {
	Datas []*Trade `json:"data"`
}

type Trade struct {
	Id   int64   `json:"id"`
	Data []*Data `json:"data"`
}

type Data struct {
	Amount    float64 `json:"amount"`
	Ts        float64 `json:"ts"`
	id        float64
	Price     float64 `json:"price"`
	Direction string  `json:"direction"`
}

// 获取Huobi单个种的数据(币币交易)depth,并且发送到redis
func GetCionDepthHuobi(symbols [][2]string) {
Loop:
	for {
		time.Sleep(time.Millisecond * 200)
		for _, symbol := range symbols {
			depthUrl := "https://api.huobipro.com/market/depth?symbol=" + symbol[0] + symbol[1] + "&type=step0"
			depthRes, err := http.Get(depthUrl)
			if err != nil {
				logs.Error("http.Get failed for huobi err:", err)
				depthRes.Body.Close()
				continue Loop
			}
			depthDoc, err := goquery.NewDocumentFromReader(depthRes.Body)
			if err != nil {
				logs.Error("goquery.NewDocumentFromReader failed err:", err)
				depthRes.Body.Close()
				continue Loop
			}

			tempDepth := &HuobiDepthRes{
				Tick: &HuobiDepth{},
			}
			curDepth := &HuobiDepth{}

			err = json.Unmarshal([]byte(depthDoc.Text()), tempDepth)
			if err != nil {
				logs.Error("json.Unmarshal huobi depth failed err:", err)
				depthRes.Body.Close()
				continue Loop
			}
			curDepth.Asks = tempDepth.Tick.Asks[len(tempDepth.Tick.Asks)-6:len(tempDepth.Tick.Asks)]
			curDepth.Bids = tempDepth.Tick.Bids[:6]

			data, err := json.Marshal(curDepth)
			if err != nil {
				logs.Error("json.marshal huobi depth failed err:", err)
				depthRes.Body.Close()
				continue Loop
			}
			// 将获得的depths数据以json格式发送到redis
			curSymbol := symbol[0] + "_" + symbol[1]
			SendDepthsToRedis2ZT(string(data), curSymbol)

			depthRes.Body.Close()
		}
	}
}

func GetTradesHuobi(symbols [][2]string, size string) {
	var preTradeId int64
Loop:
	for {
		time.Sleep(time.Millisecond * 200)
		for _, symbol := range symbols {
			okexSymbol := symbol[0] + "_" + symbol[1]
			symbol := symbol[0] + symbol[1]
			tradesUrl := "https://api.huobi.pro/market/history/trade?symbol=" + symbol + "&size=" + size
			tradesRes, err := http.Get(tradesUrl)

			if err != nil {
				logs.Error("http.Get trades failed from huobi err:", err)
				continue Loop
			}
			doc, err := goquery.NewDocumentFromReader(tradesRes.Body)
			if err != nil {
				logs.Error("goquery.NewDocumentFromReader failed err:", err)
				continue Loop
			}
			var huobiTrades = &HuobiTrades{
				Datas: []*Trade{
				},
			}
			err = json.Unmarshal([]byte(doc.Text()), huobiTrades)
			if err != nil {
				logs.Error("json.Unmarshal huobi trades failed err:", err)
				continue Loop
			}
			for _, trade := range huobiTrades.Datas {
				// 最新的交易的第一个交易ID和原来最后一交易ID比较，大于则说明拿到的trades都是有效的
				if trade.Id > preTradeId {
					rand.Seed(time.Now().Unix())
					for _, data := range trade.Data {
						postDataLimit := &utils.ZTPostDataLimit{}
						postDataLimit.Market = okexSymbol

						if data.Direction == "buy" {
							if rand.Intn(10)%2 == 0 {
								postDataLimit.Side = "2"
								// 买 设置价格
								price := data.Price*6.88 + rand.Float64() + rand.Float64()
								p := strconv.FormatFloat(price, 'E', -1, 64)
								postDataLimit.Price = p
							} else {
								postDataLimit.Side = "2"
								// 设置价格
								price := data.Price*6.88 - rand.Float64() - rand.Float64()
								p := strconv.FormatFloat(price, 'E', -1, 64)
								postDataLimit.Price = p
							}
						} else {
							if rand.Intn(10)%2 == 0 {
								postDataLimit.Side = "1"
								// 卖 设置价格
								price := data.Price*6.88 - rand.Float64() - rand.Float64()
								p := strconv.FormatFloat(price, 'E', -1, 64)
								postDataLimit.Price = p
							} else {
								postDataLimit.Side = "1"
								// 设置价格
								price := data.Price*6.88 + rand.Float64() + rand.Float64()
								p := strconv.FormatFloat(price, 'E', -1, 64)
								postDataLimit.Price = p
							}
						}
						a := strconv.FormatFloat(data.Amount, 'E', -1, 64)
						postDataLimit.Amount = a

						data, err := json.Marshal(postDataLimit)
						if err != nil {
							logs.Error("json.Marshal trade failed err:", err)
							tradesRes.Body.Close()
							continue Loop
						}
						//fmt.Println(string(data))
						models.TradesChan <- string(data)
					}
				}
			}
			if len(huobiTrades.Datas) != 0 {
				preTradeId = huobiTrades.Datas[0].Id
			}
			tradesRes.Body.Close()
		}
	}
}
