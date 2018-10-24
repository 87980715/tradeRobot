package dataAgent

import (
	"github.com/astaxie/beego/logs"
	"encoding/json"
	"tradeRobot/models"
	"strings"
	"tradeRobot/initialize"
	"tradeRobot/utils"
)

func SendTickersToRedis(data, cionName string, toCionNames map[int]string) {
	conn := initialize.RedisPool.Get()
	defer conn.Close()
	for _, toCionName := range toCionNames {
		args := cionName + "_" + toCionName + "_" + "tickers"
		_, err := conn.Do("LPUSH", args, data)
		if err != nil {
			logs.Error("LPUSH tickers failed, err:", err)
		}
	}
}

func SendDepthsToRedis2ZT(data, symbol string) {

	conn := initialize.RedisPool.Get()
	defer conn.Close()

	_, err := conn.Do("LPUSH", symbol, data)
	if err != nil {
		logs.Error("LPUSH depth failed, err:", err)
	}
}

func SendOrdersToredis(records []*utils.Record)  {
	for _, record := range records {
		strs := strings.Split(record.Market, "_")
		// 转换成固定格式的交易对
		symbol := strs[0] + "_" + "usdt"
		// 判断交易对属于Huobi 还是 Okex
		if _, ok := models.OKexSymbols[symbol]; ok {
			if record.Type == 1 {
				okexPostDataLimit := &utils.OKexPostDataLimit{}
				okexPostDataLimit.Amount = record.Amount
				okexPostDataLimit.Price = record.Price
				okexPostDataLimit.Instrument_id = symbol
				okexPostDataLimit.Type = "buy"
			}
			okexPostDataLimit := &utils.OKexPostDataLimit{}
			okexPostDataLimit.Amount = record.Amount
			okexPostDataLimit.Price = record.Price
			okexPostDataLimit.Instrument_id = symbol
			okexPostDataLimit.Type = "sell"

			dataBytes, err := json.Marshal(okexPostDataLimit)
			conn := initialize.RedisPool.Get()
			defer conn.Close()
			_, err = conn.Do("LPUSH", "OKex", string(dataBytes))
			if err != nil {
				logs.Error("SendOrdersToRedis failed err:", err)
				return
			}
		}
	}
}

func  SendTradesToRedis(trade string)  {
	conn := initialize.RedisPool.Get()
	defer conn.Close()
	_, err := conn.Do("LPUSH", "T", trade)
	if err != nil {
		logs.Error("SendOrdersToRedis failed err:", err)
		return
	}
}

/*
symbol := strs[0] + "usdt"
finisdOrderMap[HuobiSymbol] = tempArray
tempArray[0] = v.Amount
tempArray[1] = v.Price
dataBytes,_ := json.Marshal(finisdOrderMap)
SendOrdersToRedis(string(dataBytes),HuobiSymbol)
if err != nil {
	logs.Error("SendOrdersToRedis failed err:",err)
	return
}
HuobiSymbolsChan <- HuobiSymbol
}
return */
