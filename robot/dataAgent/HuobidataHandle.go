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
	"context"
	"strings"
	"github.com/tebeka/selenium/chrome"
	"github.com/tebeka/selenium"
	"math/rand"
	"fmt"
	"sync"
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

var WRMuLock sync.RWMutex

func GetTradesHuobiMtEth(symbol []string, size string, ctx context.Context) {
Loop:
	for {
		time.Sleep(time.Millisecond * 500)
		select {
		case <-ctx.Done():
			return
		default:
			// 获取火币mteth交易对的价格
			symbol_mteth := strings.ToLower(symbol[0]) + "eth"
			key := symbol_mteth + "TradeId"
			tradesUrl := "https://api.huobi.br.com/market/history/trade?symbol=" + symbol_mteth + "&size=" + size

			tradesRes, err := http.Get(tradesUrl)
			if err != nil {
				logs.Error("http.Get trades failed from huobi err:", err)
				if tradesRes != nil {
					tradesRes.Body.Close()
				}
				continue Loop
			}
			mtEthDoc, err := goquery.NewDocumentFromReader(tradesRes.Body)
			if tradesRes != nil {
				tradesRes.Body.Close()
			}

			if err != nil {
				logs.Error("goquery.NewDocumentFromReader mtEthDoc failed err:", err)
				continue Loop
			}
			var mtEthhuobiTrades = &HuobiTrades{
				Datas: []*Trade{
				},
			}
			// fmt.Println("mtEthDoc:",mtEthDoc.Text())
			err = json.Unmarshal([]byte(mtEthDoc.Text()), mtEthhuobiTrades)
			if err != nil {
				logs.Error("json.Unmarshal huobi  mtEthDoc trades failed err:", err)
				continue Loop
			}

			time.Sleep(time.Millisecond * 500)
			ZGSymbol := symbol[0] + "_" + symbol[1]
			// 获取火币ethusdt交易对的价格
			symbol_ethusdt := "ethusdt"
			key = symbol_ethusdt + "TradeId"
			tradesUrl = "https://api.huobi.br.com/market/history/trade?symbol=" + symbol_ethusdt + "&size=" + size

			tradesRes, err = http.Get(tradesUrl)
			if err != nil {
				logs.Error("http.Get trades failed from huobi err:", err)
				if tradesRes != nil {
					tradesRes.Body.Close()
				}
				continue Loop
			}
			ethUsdtDoc, err := goquery.NewDocumentFromReader(tradesRes.Body)
			if tradesRes != nil {
				tradesRes.Body.Close()
			}
			if err != nil {
				logs.Error("goquery.NewDocumentFromReader failed err:", err)
				continue Loop
			}

			var ethUsdthuobiTrades = &HuobiTrades{
				Datas: []*Trade{
				},
			}
			err = json.Unmarshal([]byte(ethUsdtDoc.Text()), ethUsdthuobiTrades)
			if err != nil {
				logs.Error("json.Unmarshal huobi  trades failed err:", err)
				continue Loop
			}

			WRMuLock.Lock()
			models.EthPrice["huobi"] = ethUsdthuobiTrades.Datas[0].Data[0].Price
			WRMuLock.Unlock()

			WRMuLock.RLock()
			usdtPrice := models.UsdtPrice["huobi"]
			WRMuLock.RUnlock()

			for _, trade := range mtEthhuobiTrades.Datas {
				// 最新的交易的第一个交易ID和原来最后一交易ID比较，大于则说明拿到的trades都是有效的
				if trade.Id > models.HuoPreTradeId[key] {
					rand.Seed(time.Now().Unix())
					for _, data := range trade.Data {
						postDataLimit := &utils.ZTPostDataLimit{}
						postDataLimit.Market = ZGSymbol

						if data.Direction == "buy" {
							postDataLimit.Side = "2"
							if rand.Intn(10)%2 == 0 {
								// 买 设置价格 增加买卖深度
								price := data.Price * models.EthPrice["huobi"] * usdtPrice * (1 + 0.005 )
								p := strconv.FormatFloat(price, 'E', -1, 64)
								postDataLimit.Price = p
							} else {
								price := data.Price*models.EthPrice["huobi"]*usdtPrice - (rand.Float64() / 1000)
								p := strconv.FormatFloat(price, 'E', -1, 64)
								postDataLimit.Price = p
							}
						} else {
							postDataLimit.Side = "1"
							if rand.Intn(10)%2 == 0 {
								// 卖 设置价格
								price := data.Price * models.EthPrice["huobi"] * usdtPrice * (1 - 0.005)
								p := strconv.FormatFloat(price, 'E', -1, 64)
								postDataLimit.Price = p
							} else {
								price := data.Price*models.EthPrice["huobi"]*usdtPrice + (rand.Float64() / 1000)
								p := strconv.FormatFloat(price, 'E', -1, 64)
								postDataLimit.Price = p
							}
						}
						a := strconv.FormatFloat(data.Amount, 'E', -1, 64)
						postDataLimit.Amount = a

						d, err := json.Marshal(postDataLimit)
						if err != nil {
							logs.Error("json.Marshal trade failed err:", err)
							tradesRes.Body.Close()
							continue Loop
						}
						// ------测试 ------
						// fmt.Println(string(d))
						models.ZGTradesChan <- string(d)
					}
				}
			}
			if len(mtEthhuobiTrades.Datas) != 0 {
				models.HuoPreTradeId[key] = mtEthhuobiTrades.Datas[0].Id
			}
		}
	}
}

func GetHuobiUsdtPrice() {

	var url = "https://otc.huobi.br.com/zh-cn/trade/buy-usdt/"
	for {

		opts := []selenium.ServiceOption{}
		caps := selenium.Capabilities{
			"browserName": "chrome",
		}

		imagCaps := map[string]interface{}{
			"profile.managed_default_content_settings.images": 2,
		}

		chromeCaps := chrome.Capabilities{
			Prefs: imagCaps,
			Path:  "",
			Args: []string{
				"--headless", // 设置Chrome无头模式,linux 必须设置，否则会报错
				"--no-sandbox",
				"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36", // 模拟user-agent，防反爬
			},
		}

		caps.AddChrome(chromeCaps)

		service, err := selenium.NewChromeDriverService("/opt/chrome/chromedriver", 9515, opts...)
		if err != nil {
			continue
		}
		defer service.Stop()

		webDriver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", 9515))
		if err != nil {
			continue
		}

		webDriver.Refresh()
		err = webDriver.Get(url)
		if err != nil {
			logs.Error("webDriver get failed")
			continue
		}

		t, err := webDriver.FindElement(selenium.ByXPATH, `//*[@id="app"]/div[1]/div[2]/div/div/div[2]/div[3]/div[1]/div/div[2]/div[4]`)
		if err != nil { //*[@id="app"]/div[1]/div[2]/div/div/div[2]/div[3]/div[1]/div/div[2]/div[4]
			continue
		}

		price, err := t.Text()
		s := strings.Split(price, " ")
		p, _ := strconv.ParseFloat(s[0], 64)

		WRMuLock.Lock()
		models.UsdtPrice["Huobi"] = p
		WRMuLock.Unlock()

		fmt.Println("usdtPrice:", p)
		time.Sleep(20 * time.Second)
	}
}

func GetTradesHuobiEthUsdt(symbol []string, size string, ctx context.Context) {
Loop:
	for {
		time.Sleep(time.Millisecond * 200)
		select {
		case <-ctx.Done():
			return
		default:
			ZGSymbol := symbol[0] + "_" + symbol[1]
			symbol := strings.ToLower(symbol[0]) + "usdt"
			key := symbol + "TradeId"
			tradesUrl := "https://api.huobi.br.com/market/history/trade?symbol=" + symbol + "&size=" + size
			tradesRes, err := http.Get(tradesUrl)
			if err != nil {
				logs.Error("http.Get trades failed from huobi err:", err)
				if tradesRes != nil {
					tradesRes.Body.Close()
				}
				continue Loop
			}
			doc, err := goquery.NewDocumentFromReader(tradesRes.Body)
			if tradesRes != nil {
				tradesRes.Body.Close()
			}
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
				if trade.Id > models.HuoPreTradeId[key] {
					rand.Seed(time.Now().Unix())
					for _, data := range trade.Data {
						postDataLimit := &utils.ZTPostDataLimit{}
						postDataLimit.Market = ZGSymbol
						if data.Direction == "buy" {
							if rand.Intn(10)%2 == 0 {
								postDataLimit.Side = "2"
								// 买 设置价格 增加买卖深度
								price := data.Price*models.UsdtPrice["Huobi"] + 0.25
								p := strconv.FormatFloat(price, 'E', -1, 64)
								postDataLimit.Price = p
							} else {
								postDataLimit.Side = "2"
								// 设置价格
								price := data.Price*models.UsdtPrice["Huobi"] - rand.Float64()
								price -= rand.Float64()
								p := strconv.FormatFloat(price, 'E', -1, 64)
								postDataLimit.Price = p
							}
						} else {
							if rand.Intn(10)%2 == 0 {
								postDataLimit.Side = "1"
								// 卖 设置价格
								price := data.Price*models.UsdtPrice["Huobi"] - 0.25
								p := strconv.FormatFloat(price, 'E', -1, 64)
								postDataLimit.Price = p
							} else {
								postDataLimit.Side = "1"
								// 设置价格
								price := data.Price*models.UsdtPrice["Huobi"] + rand.Float64()
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
				models.HuoPreTradeId[key] = huobiTrades.Datas[0].Id
			}
		}
	}
}
