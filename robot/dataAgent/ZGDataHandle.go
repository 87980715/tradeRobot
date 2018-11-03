package dataAgent

import (
	"tradeRobot/robot/utils"
	"github.com/astaxie/beego/logs"
	"encoding/json"
	"strconv"
	"fmt"
	"tradeRobot/robot/models"
	"strings"
	"time"
	"context"
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

var ZGTradeResult = make(chan *models.ZGTradeResults, 100)

// 从ZG交易所中，获取账户真实成交记录,并进行处处理
func GetFinishedOrdersFromZT(symbol []string, ctx context.Context) {
	var preTradeId int64
	for {
		select {
		case <-ctx.Done():
			return
		default:
			orders := QueryDealZG(symbol)
			var ordersFinished = &utils.ZTOrderFinishedResp{
				Result: &utils.Result{
					Records: make([]*utils.Record, 10),
				},
			}
			err := json.Unmarshal([]byte(orders), ordersFinished)
			if err != nil {
				logs.Error("json.Unmarshal curOrdersFinished failed err:", err)
				continue
			}
			records := ordersFinished.Result.Records
			for _, r := range records {
				if r.Id > preTradeId {
					dealIds := QueryDealIds(r.Id, r.User)
					//fmt.Println("dealIds:", dealIds)
					if len(dealIds) != 0 {
						record := QureyDealOerder(dealIds, r.User)
						// ------测试----
						//fmt.Println("record：", record)
						if record != nil {
							var postDataLimit = &utils.HuobiPostDataLimit{}
							str := strings.Split(record.Market, "_")
							postDataLimit.Symbol = str[0] + "usdt"
							if record.Side == 1 {
								postDataLimit.Type = "buy-limit"
								p, _ := strconv.ParseFloat(record.Price, 64)
								// 买单价格 ，去除成交手续费
								d, _ := strconv.ParseFloat(record.Deal_fee, 64)
								price := (p/6.88)*(1-0.002) - d
								postDataLimit.Price = fmt.Sprintf("%."+strconv.Itoa(4)+"f", price)
							} else {
								postDataLimit.Type = "sell-limit"
								p, _ := strconv.ParseFloat(record.Price, 64)
								// 卖单价格 ，包含成交手续费
								d, _ := strconv.ParseFloat(record.Deal_fee, 64)
								price := (p/6.88)*(1+0.002) + d
								postDataLimit.Price = fmt.Sprintf("%."+strconv.Itoa(4)+"f", price)
							}
							postDataLimit.Amount = record.Amount
							postDataLimit.Account_id = models.Huobi_Account_ID
							utils.HuobiOrders <- postDataLimit
						}
					}
				}
			}
			// 去重
			preTradeId = records[0].Id
		}
	}
}

// 查询ZG已成交订单
func QueryDealZG(symbol []string) string {

	account := utils.ZTAccount
	account.PostDataOrderFinished.Market = strings.ToUpper(symbol[0] + "_" + "CNZ")
	account.PostDataOrderFinished.Side = "0"
	account.PostDataOrderFinished.Limit = "10"
	account.PostDataOrderFinished.Offset = "0"
	account.PostDataOrderFinished.Start_time = "0"
	account.PostDataOrderFinished.Start_time = "0"

	account.ZTQueryDealMd5Sign()
	orders := account.ZTOrderFinished()

	return orders

}

func QueryDealIds(orderId, userId int64) ([]int) {

	db, err := utils.LoadExchangeDB()
	if err != nil {
		logs.Error("initMaria failed err:", err)
		return nil
	}
	defer db.Close()
	table := fmt.Sprintf("user_deal_history_%d", userId%100)
	queryStr := "select deal_id from " + table + " where order_id = " + strconv.Itoa(int(orderId))

	row, err := db.Query(queryStr)
	if err != nil {
		logs.Error("select deal_id form user_deal_history_26 failed err:", err)
		return nil
	}
	defer row.Columns()
	var dealId = 0
	dealIds := make([]int, 0)
	for row.Next() {
		row.Scan(&dealId)
		dealIds = append(dealIds, dealId)
	}
	return dealIds
}

func QureyDealOerder(dealIds []int, userId int64) (data *utils.Record) {
	db, err := utils.LoadExchangeDB()
	if err != nil {
		logs.Error("initMaria failed err:", err)
		return
	}
	defer db.Close()
	n := 0
	table := fmt.Sprintf("user_deal_history_%d", userId%100)
	var CreatedAt, UserId, TradeId, Symbol, Type, Price, DealAmount, DealFee, Total string
	for _, id := range dealIds {
		queryStr := "select COUNT(deal_id) from " + table + " where deal_id = " + strconv.Itoa(id)
		row, err := db.Query(queryStr)
		if err != nil {
			logs.Error(" select count(deal_id) failed err:", err)
			return
		}
		for row.Next() {
			row.Scan(&n)
			// -------测试-----
			// fmt.Println("n:", n)
			// 提取真实交易
			if n != 2 {
				qureyStr := "select time,user_id,market,order_id,side,price,amount,deal,deal_fee from " + table + " where deal_id = " + strconv.Itoa(id)
				r, err := db.Query(qureyStr)
				if err != nil {
					logs.Error(" select count(deal_id) failed err:", err)
					return
				}
				for r.Next() {
					row.Scan(&CreatedAt, &UserId, &Symbol, &TradeId, &Type, &Price, &DealAmount, &Total, &DealFee)
					// huobi 交易所需数据
					data.Amount = DealAmount
					data.Price = Price
					side, _ := strconv.Atoi(Type)
					data.Side = side
					data.Market = Symbol
					data.Deal_fee = DealFee

					tradeResult := &models.ZGTradeResults{}
					tradeResult.Type = Type
					tradeResult.Created_at = CreatedAt
					tradeResult.User_id = UserId
					tradeResult.Trade_id = TradeId
					tradeResult.Symbol = Symbol
					tradeResult.Price = Price
					tradeResult.Deal_amount = DealAmount
					tradeResult.Deal_fees = DealFee
					tradeResult.Total = Total

					ZGTradeResult <- tradeResult
					//fmt.Println(CreatedAt, UserId, TradeId, Symbol, Type, Price, DealAmount, DealFee, CreatedAt, Total)
				}
				r.Columns()
			}
			row.Columns()
		}
	}
	return
}

func ZGInsertToDB(ctx context.Context) {
	db, err := utils.LoadRobotDB()
	if err != nil {
		logs.Error("loadDB failed err:", err)
		return
	}
	defer db.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(time.Second)
			tradeResult := <-ZGTradeResult
			if err := db.Create(tradeResult).Error; err != nil {
				logs.Error("insert failed into Huobi tradeResult ")
				return
			}
		}
	}
}

/*
// 判断交易对是否属于OKex
				if _,ok:= models.OKexSymbols[record.Market];ok{
					var postDataLimit = &utils.OKexPostDataLimit{}

					str := strings.Split(record.Market,"_")
					instrument_id := str[0] + "-" + "usdt"

					p,_:= strconv.ParseFloat(record.Price,64)
					// 设置挂单价格，包含成交手续费
					price := (p/6.88)* (1+0.001)
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

					utils.OKexOrders <- postDataLimit
				}
 */
