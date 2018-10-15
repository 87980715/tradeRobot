package dataAgent

import (

)

// 火币response
type HuobiRes struct {
	Status string `json:"status"`
	Ch string `json:"ch"`
	Ts int64 `json:"ts"`
	Tick *HuobiDepth `json:"tick"`
}

type HuobiDepth struct {
	Asks [][2]float64 `json:"asks"`
	Bids [][2]float64 `json:"bids"`
	Ts int64 `json:"ts"`
	version int64 `json:"version"`
}

//// 获取Huobi单个币的数据(币币交易),depth
//func GetCionDataHuobi(cionName string, toCionNames map[int]string) *OKexC2CData {
//
//	var c2cDepthDataHuobi = &OKexC2CData{}
//	//tickerMap := make(map[string]*CurTicker)
//	depthMap := make(map[string]*OKexDepth)
//	curDepth := &OKexDepth{}
//	for _, toCionName := range toCionNames {
//
//		/*tickerUrl := "https://www.okex.com/api/v1/ticker.do?symbol=" + cionName + "_" + toCionName
//		tickerRes, err := http.Get(tickerUrl)
//		if err != nil {
//			logs.Error("http.Get failed for OKex err:", err)
//		}
//		tickerDoc, err := goquery.NewDocumentFromReader(tickerRes.Body)
//
//		var curTicker = &CurTicker{}
//		err = json.Unmarshal([]byte(tickerDoc.Text()), &curTicker)
//		if err != nil {
//			fmt.Println("json.Unmarshal failed:", err)
//		}
//		key := cionName + "ticker"
//		tickerMap[key] = curTicker
//		c2cDataHuobi.tickers = append(c2cDataHuobi.tickers, tickerMap)*/
//
//		depthUrl := "https://api.huobipro.com/market/depth?symbol="+cionName+toCionName+"btc&type=step1"
//		depthRes, err := http.Get(depthUrl)
//		if err != nil {
//			logs.Error("http.Get failed for OKex err:", err)
//		}
//		if depthRes.Status == "OK" {
//			depthDoc, err := goquery.NewDocumentFromReader(depthRes.Body)
//			err = json.Unmarshal([]byte(depthDoc.Text()), &curDepth)
//			if err != nil {
//				fmt.Println("json.Unmarshal depthDoc.Text() failed:", err)
//			}
//
//			key := cionName + "depth"
//			fmt.Println("curDepth.Asks",curDepth.Asks[len(curDepth.Asks)-5:len(curDepth.Asks)])
//			fmt.Println("curDepth.Bids",curDepth.Bids[:5])
//
//			depthMap[key] = curDepth
//			c2cDepthDataHuobi.Depths = append(c2cDepthDataHuobi.Depths, depthMap)
//		}
//	}
//	return c2cDepthDataHuobi
//}
//// 获取各个Okex交易所指定的所有币种的数据
//func GetDataFromHuobi(coinNames map[int]string, toCionNames map[int]string) map[string]*OKexC2CData {
//	HuobiDataMap := make(map[string]*OKexC2CData)
//	for _, coinName := range coinNames {
//		go func(coinName string, toCionNames map[int]string) {
//			HuobiDataMap[coinName] = GetCionDataOKex(coinName, toCionNames)
//		}(coinName, toCionNames)
//	}
//	return HuobiDataMap
//}
