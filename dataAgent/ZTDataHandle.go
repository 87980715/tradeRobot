package dataAgent

import (
	"tradeRobot/utils"
	"github.com/astaxie/beego/logs"
	"encoding/json"

	"strconv"

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

// 从ZG交易所中，获取账户真实成交记录,并进行处处理
func GetFinishedOrdersFromZT(account *utils.ZTRestfulApiRequest) {
	for {
		orders, err := account.ZTOrderFinished()
		if err != nil {
			logs.Error("GetFinishedOrdersFromZT failed err:", err)
			continue
		}
		var curOrdersFinished = &utils.ZTOrderFinishedResp{
			Result: &utils.Result{
				Records: make([]*utils.Record, 0),
			},
		}
		err = json.Unmarshal([]byte(orders), curOrdersFinished)
		if err != nil {
			logs.Error("json.Unmarshal curOrdersFinished failed err:", err)
			continue
		}
		records := curOrdersFinished.Result.Records
		offset := len(curOrdersFinished.Result.Records)

		SendOrdersToredis(records)

		preOffset, _ := strconv.Atoi(account.PostDataOrderFinished.Offset)
		NexOffset := preOffset + offset
		account.PostDataOrderFinished.Offset = strconv.Itoa(NexOffset)
	}
}


