package dataAgent

import (
	"tradeRobot/robot/utils"
	"github.com/astaxie/beego/logs"
	"strconv"
	"fmt"
	"tradeRobot/robot/models"
	"strings"
	"time"
	"context"
	"database/sql"
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
var ZGTradeRecords = make(chan *utils.Record, 100)
var ZGDealIds = make(chan []int, 100)

// 从ZG交易所中，获取账户真实成交记录,并进行处处理
func GetDealOrdersZG(ctx context.Context) {
	for {
		time.Sleep(1000 * time.Millisecond)
		select {
		case <-ctx.Done():
			return
		default:
			record := <-ZGTradeRecords
			if record != nil {
				var postDataLimit = &utils.HuobiPostDataLimit{}
				str := strings.Split(record.Market, "_")
				postDataLimit.Symbol = strings.ToLower(str[0] + "eth")
				if record.Side == 1 {
					postDataLimit.Type = "buy-limit"
					p, _ := strconv.ParseFloat(record.Price, 64)
					// 买单价格 ，去除成交手续费
					price := p*(1-0.002)
					postDataLimit.Price = fmt.Sprintf("%."+strconv.Itoa(8)+"f", price)
				} else {
					postDataLimit.Type = "sell-limit"
					p, _ := strconv.ParseFloat(record.Price, 64)
					// 卖单价格 ，包含成交手续费
					price := p*(1+0.002)
					postDataLimit.Price = fmt.Sprintf("%."+strconv.Itoa(8)+"f", price)
				}
				postDataLimit.Amount = record.Amount
				postDataLimit.Account_id = models.Huobi_Account_ID
				//------测试-------
				// fmt.Println("火币挂单数据:", postDataLimit)
				utils.HuobiOrders <- postDataLimit

			}
		}
	}
}

// 查询ZT已成交订单
func QueryRealDealZG(ctx context.Context) {

	userId, _ := strconv.ParseInt(models.ZGUserID, 10, 64)
	for {
		time.Sleep(1000 * time.Millisecond)
		select {
		case <-ctx.Done():
			return
		default:
			dealIds := <-ZGDealIds
			if len(dealIds) != 0 {
				QureyDealOerder(utils.ExchangeDB, dealIds, userId)
			}
		}

	}
}

func QueryDealIds(symbol []string, ctx context.Context) {

	userId, _ := strconv.ParseInt(models.ZGUserID, 10, 64)
	for {
		time.Sleep(1000 * time.Millisecond)
		select {
		case <-ctx.Done():
			return
		default:
			key := symbol[0] + "cnt" + "ZGDealId"
			table := fmt.Sprintf("user_deal_history_%d", userId%100)
			queryStr := "select deal_id from " + table + " order by id desc limit " + models.ZGQueryDealOrderSize
			//fmt.Println("queryStr:",queryStr)
			row, err := utils.ExchangeDB.Query(queryStr)
			if err != nil {
				logs.Error("select deal_id form %s failed err %v:", table, err)
				return
			}
			var dealId = 0
			dealIds := make([]int, 0)
			for row.Next() {
				row.Scan(&dealId)
				if int64(dealId) > models.ZGPreDealId[key] {
					dealIds = append(dealIds, dealId)
				}
			}
			if len(dealIds) != 0 {
				models.ZGPreDealId[key] = int64(dealIds[0])
			}
			row.Close()
			ZGDealIds <- dealIds
		}
	}
}
func QureyDealOerder(db *sql.DB, dealIds []int, userId int64) {
	logs.Info("dealIds:", dealIds)
	m := make(map[int]int)
	// 统计每个dealId 出现的次数
	for _, v := range dealIds {
		if m[v] != 0 {
			m[v] ++
		} else {
			m[v] = 1
		}
	}
	table := fmt.Sprintf("user_deal_history_%d", userId%100)
	var UserId, TradeId, Symbol, Type, Price, DealAmount, DealFee, Total string
	var CreatedAt float64
	var n int
	// 出现两次的即为虚拟交易
	for k, v := range m {
		if v == 2 {
			continue
		}
		//fmt.Println("dealId:", k)
		queryStr := "select COUNT(deal_id) from " + table + " where deal_id = " + strconv.Itoa(k) + " and " + " user_id = " + strconv.Itoa(int(userId))
		row, err := db.Query(queryStr)
		if err != nil {
			logs.Error(" select count(deal_id) failed err:", err)
			return
		}
		for row.Next() {
			row.Scan(&n)
			// fmt.Println("n:", n)
			// 提取真实交易
			if n != 2 {
				qureyStr := "select time,user_id,market,order_id,side,price,amount,deal,deal_fee from " + table + " where deal_id = " + strconv.Itoa(k)
				r, err := db.Query(qureyStr)
				if err != nil {
					logs.Error(" select count(deal_id) failed err:", err)
					return
				}
				for r.Next() {
					r.Scan(&CreatedAt, &UserId, &Symbol, &TradeId, &Type, &Price, &DealAmount, &Total, &DealFee)
					// huobi 交易所需数据
					//fmt.Println("数据库数据:",CreatedAt, UserId, TradeId, Symbol, Type, Price, DealAmount, DealFee, Total)
					data := &utils.Record{}
					data.Amount = DealAmount
					data.Price = Price
					side, _ := strconv.Atoi(Type)
					data.Side = side
					data.Market = Symbol
					data.Deal_fee = DealFee
					ZGTradeRecords <- data
					// 写入数据库的数据
					tradeResult := &models.ZGTradeResults{}
					tradeResult.Type = Type
					CreatedAt = CreatedAt * 1000
					tradeResult.Created_at = strconv.FormatFloat(CreatedAt, 'E', -1, 64)
					tradeResult.User_id = UserId
					tradeResult.Trade_id = TradeId
					tradeResult.Symbol = Symbol
					tradeResult.Price = Price
					tradeResult.Deal_amount = DealAmount
					tradeResult.Deal_fees = DealFee
					tradeResult.Total = Total
					ZGTradeResult <- tradeResult
				}
				r.Close()
			}
			row.Close()
		}
	}
}

func ZGInsertToDB(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(time.Second)
			tradeResult := <-ZGTradeResult
			db := utils.RobotDB
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
