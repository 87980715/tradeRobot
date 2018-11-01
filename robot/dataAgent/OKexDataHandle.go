package dataAgent

import (
	"net/http"
	"github.com/astaxie/beego/logs"
	"github.com/PuerkitoBio/goquery"
	"tradeRobot/robot/models"
	"encoding/json"
	"time"
	"strconv"
	"tradeRobot/robot/utils"
	"math/rand"
)

type OKexDepth struct {
	Asks [][2]float64 `json:"asks"`
	Bids [][2]float64 `json:"bids"`
}

type OKexTrade struct {
	Time      string `json:"time"`
	Timestamp string `json:"timestamp"`
	Trade_id  string `json:"trade_id"`
	Price     string `json:"price"`
	Size      string `json:"size"`
	Side      string `json:"side"`
	symbol    string
}
/*
// 获取Okex单个种的数据(币币交易),ticker,depth,并且发送到redis
func GetCionDepthOKex(symbols [][2]string) {
Loop:
	for {
		time.Sleep(time.Millisecond*200)
		for _, symbol := range symbols {
			depthUrl := models.Okex_API_URL + "depth.do?symbol=" + symbol[0] + "_" + symbol[1]
			depthRes, err := http.Get(depthUrl)
			if err != nil {
				logs.Error("http.Get failed for OKex err:", err)
				depthRes.Body.Close()
				continue Loop
			}
			depthDoc, err := goquery.NewDocumentFromReader(depthRes.Body)
			if err != nil {
				logs.Error("goquery.NewDocumentFromReader failed err:", err)
				depthRes.Body.Close()
				continue Loop
			}

			curDepth := &OKexDepth{}
			err = json.Unmarshal([]byte(depthDoc.Text()), curDepth)
			if err != nil {
				logs.Error("json.Unmarshal okex depth failed err:", err)
				depthRes.Body.Close()
				continue Loop
			}

			// 截取数据多少可调整
			curDepth.Asks = curDepth.Asks[len(curDepth.Asks)-10:len(curDepth.Asks)]
			curDepth.Bids = curDepth.Bids[:10]

			data, err := json.Marshal(curDepth)
			if err != nil {
				logs.Error("json.marshal okex depth failed err:", err)
				depthRes.Body.Close()
				continue Loop
			}
			// 将获得的depths数据以json格式发送到redis
			curSymbol := symbol[0] + "_" + symbol[1]
			//logs.Info(string(data))
			SendDepthsToRedis2ZT(string(data), curSymbol)
			depthRes.Body.Close()
		}
	}
}
*/

func GetTradesOKex(symbols [][2]string, limit string) {
	var preTradeId int64
Loop:
	for {
		time.Sleep(time.Millisecond * 100)
		for _, symbol := range symbols {
			symbol := symbol[0] + "_" + symbol[1]
			tradesUrl := "https://www.okex.com/api/spot/v3/instruments/" + symbol + "/trades?limit=" + limit
			tradesRes, err := http.Get(tradesUrl)
			if err != nil {
				logs.Error("http.Get trades failed from OKex err:", err)
				continue Loop
			}
			doc, err := goquery.NewDocumentFromReader(tradesRes.Body)
			if err != nil {
				logs.Error("goquery.NewDocumentFromReader failed err:", err)
				continue Loop
			}
			l, _ := strconv.Atoi(limit)
			var curTrades = make([]*OKexTrade, l)
			err = json.Unmarshal([]byte(doc.Text()), &curTrades)
			if err != nil {
				logs.Error("OkexTrades:",doc.Text())
				logs.Error("json.Unmarshal okex trades failed err:", err)
				continue Loop
			}
			for _, trade := range curTrades {
				// 最新的交易的第一个交易ID和原来最后一交易ID比较，大于则说明拿到的trades都是有效的
				newTradeId, _ := strconv.Atoi(trade.Trade_id)
				//
				if int64(newTradeId) > preTradeId {
					rand.Seed(time.Now().Unix())
					postDataLimit := &utils.ZTPostDataLimit{}
					postDataLimit.Market = symbol
					if trade.Side == "buy" {
						// 买
						if rand.Intn(10)%2 == 0 {
							postDataLimit.Side = "2"
							p, _ := strconv.ParseFloat(trade.Price, 64)
							// 设置价格
							price := p*6.88 + rand.Float64()
							postDataLimit.Price = strconv.FormatFloat(price, 'E', -1, 64)
						}else{
							postDataLimit.Side = "2"
							p, _ := strconv.ParseFloat(trade.Price, 64)
							price := p*6.88 - rand.Float64()
							postDataLimit.Price = strconv.FormatFloat(price, 'E', -1, 64)
						}
					} else {
						// 卖
						if rand.Intn(10)%2 == 0 {
							postDataLimit.Side = "1"
							p, _ := strconv.ParseFloat(trade.Price, 64)
							// 设置价格
							price := p*6.88 + rand.Float64()
							postDataLimit.Price = strconv.FormatFloat(price, 'E', -1, 64)
						} else {
							postDataLimit.Side = "1"
							p, _ := strconv.ParseFloat(trade.Price, 64)
							price := p*6.88 - rand.Float64()
							postDataLimit.Price = strconv.FormatFloat(price, 'E', -1, 64)
						}
					}
					postDataLimit.Amount = trade.Size
					data, err := json.Marshal(postDataLimit)
					if err != nil {
						logs.Error("json.Marshal trade failed err:", err)
					}
					models.ZGTradesChan <- string(data)
				}
			}
			Id, _ := strconv.Atoi(curTrades[0].Trade_id)
			preTradeId = int64(Id)
			tradesRes.Body.Close()
		}
	}
}

