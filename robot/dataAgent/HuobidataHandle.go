package dataAgent

import (
	"net/http"
	"github.com/astaxie/beego/logs"
	"github.com/PuerkitoBio/goquery"
	"encoding/json"
	"time"
	"strconv"
	"tradeRobot/robot/utils"
	"tradeRobot/robot/models"
	"math/rand"
	"context"
)

var StopChan chan string

type HuobiDepthRes struct {
	Status string      `json:"status"`
	Ch     string      `json:"ch"`
	Ts     float64     `json:"ts"`
	Tick   *HuobiDepth `json:"tick"`
}

type HuobiDepth struct {
	Bids [][2]float64 `json:"bids"`
	Asks [][2]float64 `json:"asks"`
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

func GetTradesHuobi(symbol []string, size string, ctx context.Context) {
	var preTradeId int64
Loop:
	for {
		time.Sleep(time.Millisecond * 200)
		select {
		case <-ctx.Done():
			//StopChan <- "stopped"
			return
		default:
			ZGSymbol := symbol[0] + "_" + symbol[1]
			symbol := symbol[0] + symbol[1]

			tradesUrl := "https://api.huobi.br.com/market/history/trade?symbol=" + symbol + "&size=" + size
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
						postDataLimit.Market = ZGSymbol
						if data.Direction == "buy" {
							if rand.Intn(10)%2 == 0 {
								postDataLimit.Side = "2"
								// 买 设置价格 增加买卖深度
								price := data.Price*6.88 + 0.25
								p := strconv.FormatFloat(price, 'E', -1, 64)
								postDataLimit.Price = p
							} else {
								postDataLimit.Side = "2"
								// 设置价格
								price := data.Price*6.88 - rand.Float64()
								price -= rand.Float64()
								p := strconv.FormatFloat(price, 'E', -1, 64)
								postDataLimit.Price = p
							}
						} else {
							if rand.Intn(10)%2 == 0 {
								postDataLimit.Side = "1"
								// 卖 设置价格
								price := data.Price*6.88 - 0.25
								p := strconv.FormatFloat(price, 'E', -1, 64)
								postDataLimit.Price = p
							} else {
								postDataLimit.Side = "1"
								// 设置价格
								price := data.Price*6.88 + rand.Float64()
								price += rand.Float64()
								p := strconv.FormatFloat(price, 'E', -1, 64)
								postDataLimit.Price = p
							}
						}
						// 防止特大单，出现平行线
						if data.Amount > 50 {
							data.Amount = 30
							a := strconv.FormatFloat(data.Amount, 'E', -1, 64)
							postDataLimit.Amount = a
						} else {
							a := strconv.FormatFloat(data.Amount, 'E', -1, 64)
							postDataLimit.Amount = a
						}
						data, err := json.Marshal(postDataLimit)
						if err != nil {
							logs.Error("json.Marshal trade failed err:", err)
							tradesRes.Body.Close()
							continue Loop
						}
						// ------测试 ------
						models.ZGTradesChan <- string(data)
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
