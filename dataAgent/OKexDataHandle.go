package dataAgent

import (
	"net/http"
	"github.com/astaxie/beego/logs"
	"github.com/PuerkitoBio/goquery"
	"tradeRobot/models"
	"encoding/json"
	"encoding/gob"
	"bytes"
	"time"
)

type OKexDepth struct {
	Asks [][2]float64 `json:"asks"`
	Bids [][2]float64 `json:"bids"`
}

// 获取Okex单个种的数据(币币交易),ticker,depth,并且发送到redis
func GetCionDepthOKex(symbols [][2]string) {
Loop:
	for {
		time.Sleep(time.Millisecond*80)
		for _, symbol := range symbols {
			depthUrl := models.Okex_API_URL + "depth.do?symbol=" + symbol[0] + "_" + symbol[1]
			depthRes, err := http.Get(depthUrl)
			if err != nil {
				logs.Error("http.Get failed for OKex err:", err)
				continue
			}
			depthDoc, err := goquery.NewDocumentFromReader(depthRes.Body)
			if err != nil {
				logs.Error("goquery.NewDocumentFromReader failed err:", err)
				continue Loop
			}

			curDepth := &OKexDepth{}
			err = json.Unmarshal([]byte(depthDoc.Text()), curDepth)
			if err != nil {
				logs.Error("json.Unmarshal okex depth failed err:", err)
				continue Loop
			}

			// 截取数据多少可调整
			curDepth.Asks = curDepth.Asks[len(curDepth.Asks)-10:len(curDepth.Asks)]
			curDepth.Bids = curDepth.Bids[:10]

			data, err := json.Marshal(curDepth)
			if err != nil {
				logs.Error("json.marshal okex depth failed err:", err)
				continue Loop
			}
			// 将获得的depths数据以json格式发送到redis
			curSymbol := symbol[0] + "_" + symbol[1]
			//logs.Info(string(data))
			SendDepthsToRedis2ZT(string(data), curSymbol)
		}
	}
}

// 将结构体转化成字节数组
func Encode(data *OKexDepth) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}
	dataBytes := buf.Bytes()
	return dataBytes, nil
}
