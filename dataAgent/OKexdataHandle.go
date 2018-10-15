package dataAgent

import (
	"net/http"
	"github.com/astaxie/beego/logs"
	"github.com/PuerkitoBio/goquery"
)

// 获取Okex单个种的数据(币币交易),ticker,depth,并且发送到redis
func GetCionDataOKex(cionName string, toCionNames map[int]string) {

	for _, toCionName := range toCionNames {
		/*
		tickerUrl := "https://www.okex.com/api/v1/ticker.do?symbol=" + cionName + "_" + toCionName
		tickerRes, err := http.Get(tickerUrl)
		if err != nil {
			logs.Error("http.Get failed for OKex err:", err)
			return
		}
		tickerDoc, err := goquery.NewDocumentFromReader(tickerRes.Body)
		if err != nil {
			logs.Error("goquery.NewDocumentFromReader failed err:", err)
			return
		}
		//将获得的tickers数据以json格式发送到redis
		go func(tickerDoc *goquery.Document) {
			SendTickersToRedis(tickerDoc.Text(), cionName, toCionNames)
		}(tickerDoc)
		*/

		depthUrl := "https://www.okex.com/api/v1/depth.do?symbol=" + cionName + "_" + toCionName
		depthRes, err := http.Get(depthUrl)
		if err != nil {
			logs.Error("http.Get failed for OKex err:", err)
			return
		}
		depthDoc, err := goquery.NewDocumentFromReader(depthRes.Body)
		if err != nil {
			logs.Error("goquery.NewDocumentFromReader failed err:", err)
			return
		}
		// 将获得的depths数据以json格式发送到redis
		go func(depthDoc *goquery.Document) {
			SendDepthsToRedis(depthDoc.Text(), cionName, toCionNames)
		}(depthDoc)
	}
}

// 从Okex交易所中，获取指定币种的数据
func GetDataFromOkex(coinNames map[int]string, toCionNames map[int]string) { //map[string]*OKexC2CData {
	for _, coinName := range coinNames {
		GetCionDataOKex(coinName, toCionNames)
	}
}

// 获取所有指定的交易所的数据
// 将所有的需要数据进行处理
/*
func GetAndSendDataTotal(Exchanges map[string]bool, coinNames map[int]string, toCionNames map[int]string) {
	// 判断是否需要获取
	ExchangesDataMap := make(map[string]map[string]*OKexC2CData)
	_, ok := Exchanges["OKex"]
	if ok {
		OKexDataMap := GetDataFromOkex(coinNames, toCionNames)
		ExchangesDataMap["OKex"] = OKexDataMap
	}
	_, ok = Exchanges["Huobi"]
	if ok {
		OKexDataMap := GetDataFromHuobi(coinNames, toCionNames)
		ExchangesDataMap["Huobi"] = OKexDataMap
	}
	//dataChan <- ExchangesDataMap
}
*/
