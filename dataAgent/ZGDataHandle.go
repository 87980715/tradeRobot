package dataAgent

import (
	"tradeRobot/utils"
	"github.com/astaxie/beego/logs"
	"encoding/json"
	"strconv"
	"fmt"
	"tradeRobot/models"
	"strings"
)

type ZTCurTicker struct {
	TimeStamp string    `json:"time_stamp"`
	Ticker    *ZTTicker `json:"ticker"`
}

type ZTTicker struct {
	Buy    string `json:"buy"`
	High   string `json:"high"`
	Last   string `json:"last"`
	Low    string `json:"low"`
	Sell   string `json:"sell"`
	Symbol string `json:"symbol"`
	Vol    string `json:"vol"`
}

var okexOrders = make(chan *utils.OKexPostDataLimit,100)
var huobiOrders = make(chan *utils.HuobiPostDataLimit,100)

// 从ZG交易所中，获取账户真实成交记录,并进行处处理
func GetFinishedOrdersFromZT(account *utils.ZTRestfulApiRequest) {
	var preTradeId int64
	for {
		orders, err := account.ZTOrderFinished()
		if err != nil {
			logs.Error("GetFinishedOrdersFromZT failed err:", err)
			continue
		}
		var ordersFinished = &utils.ZTOrderFinishedResp{
			Result: &utils.Result{
				Records: make([]*utils.Record, 0),
			},
		}
		err = json.Unmarshal([]byte(orders), ordersFinished)
		if err != nil {
			logs.Error("json.Unmarshal curOrdersFinished failed err:", err)
			continue
		}
		records := ordersFinished.Result.Records
		for _,record := range records {
			if record.Id > preTradeId {
				// 判断交易对是否属于OKex
				if _,ok:= models.OKexSymbols[record.Market];ok{
					var postDataLimit = &utils.OKexPostDataLimit{}
					str := strings.Split(record.Market,"_")
					instrument_id := str[0] + "-" + "usdt"
					p,_:= strconv.ParseFloat(record.Price,64)
					price := p / 6.88
					postDataLimit.Price = fmt.Sprintf("%."+strconv.Itoa(4)+"f", price)
					postDataLimit.Amount = record.Amount
					postDataLimit.Instrument_id = instrument_id
					if record.Side == 1 {
						postDataLimit.Side = "buy"
					}else{
						postDataLimit.Side = "sell"
						}
					postDataLimit.Type = "limit"
					postDataLimit.Side = strconv.Itoa(record.Side)
					postDataLimit.Price = record.Price

					okexOrders <- postDataLimit
				}
				// 属于火币的交易对
				var postDataLimit = &utils.HuobiPostDataLimit{}
				str := strings.Split(record.Market,"_")
				postDataLimit.Symbol = str[0] + "usdt"
				p,_:= strconv.ParseFloat(record.Price,64)
				price := p / 6.88
				postDataLimit.Price = fmt.Sprintf("%."+strconv.Itoa(4)+"f", price)
				if record.Side == 1 {
					postDataLimit.Type = "buy-limit"
				} else {
					postDataLimit.Type = "sell-limit"
				}
				postDataLimit.Amount = record.Amount
				postDataLimit.Account_id = models.Huobi_Account_ID

				huobiOrders <- postDataLimit
			}
		}
		//
		preTradeId = records[0].Id
	}
}

// 获取未成交的订单

