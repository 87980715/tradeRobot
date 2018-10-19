package dataAgent

import (
	"net/http"
	"github.com/astaxie/beego/logs"
	"github.com/PuerkitoBio/goquery"
	"encoding/json"
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
	//Ts      float64      `json:"ts"`
	//Version float64      `json:"version"`
}

// 子账号API Key 现不能绑定IP， 有效期为90天
// 获取Huobi单个种的数据(币币交易)depth,并且发送到redis
func GetCionDepthHuobi(symbols [][2]string) {
Loop:
	for {
		//time.Sleep(time.Millisecond*200)
		for _, symbol := range symbols {
			depthUrl := "https://api.huobipro.com/market/depth?symbol=" + symbol[0] + symbol[1] + "&type=step0"
			//fmt.Println(depthUrl)
			depthRes, err := http.Get(depthUrl)
			if err != nil {
				logs.Error("http.Get failed for OKex err:", err)
				continue Loop
			}
			depthDoc, err := goquery.NewDocumentFromReader(depthRes.Body)
			if err != nil {
				logs.Error("goquery.NewDocumentFromReader failed err:", err)
				continue Loop
			}

			tempDepth := &HuobiDepthRes{
				Tick: &HuobiDepth{},
			}
			curDepth := &HuobiDepth{}

			err = json.Unmarshal([]byte(depthDoc.Text()), tempDepth)
			if err != nil {
				logs.Error("json.Unmarshal huobi depth failed err:", err)
				continue Loop
			}
			curDepth.Asks = tempDepth.Tick.Asks[len(tempDepth.Tick.Asks)-6:len(tempDepth.Tick.Asks)]
			curDepth.Bids = tempDepth.Tick.Bids[:6]

			data, err := json.Marshal(curDepth)
			if err != nil {
				logs.Error("json.marshal huobi depth failed err:", err)
				continue Loop
			}
			// 将获得的depths数据以json格式发送到redis
			curSymbol := symbol[0] + "_" + symbol[1]
			//fmt.Println(string(data))
			SendDepthsToRedis2ZT(string(data), curSymbol)
		}
	}
}
