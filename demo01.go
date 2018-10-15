package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"fmt"
	"encoding/json"
	"strings"
	"net/url"
	"io/ioutil"
)

type CurTicker struct {
	Date   string  `json:"date"`
	Ticker *Ticker `json:"ticker"`
}

type Ticker struct {
	Buy  string `json:"buy"`
	High string `json:"high"`
	Last string `json:"last"`
	Low  string `json:"low"`
	Sell string `json:"sell"`
	Vol  string `json:"vol"`
}

type Depth struct {
	Asks [][2]float64 `json:"asks"`
	Bids [][2]float64 `json:"bids"`
}

func main() {
	//GetCionData("ltc")
	//GetDepthData("true")
	//md5Hash()
	post()
}


func GetCionData(cionName string) {
	url := "https://www.okex.com/api/v1/ticker.do?symbol=" + cionName + "_btc"
	res, err := http.Get(url)
	if err != nil {
		logs.Error("http.Get failed for OKex err:", err)
	}
	var curTicker = &CurTicker{}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	fmt.Println(doc.Text())
	err = json.Unmarshal([]byte(doc.Text()), &curTicker)
	if err != nil {
		fmt.Println("json.Unmarshal failed:", err)
	}
	fmt.Println("curTicker:", curTicker.Ticker.Buy)
}

func GetDepthData(cionName string) {
	url := "https://www.okex.com/api/v1/depth.do?symbol=" + cionName + "_btc"
	res, err := http.Get(url)
	if err != nil {
		logs.Error("http.Get failed for OKex err:", err)
	}
	var curDepth = &Depth{}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	//fmt.Println(doc.Text())
	err = json.Unmarshal([]byte(doc.Text()), &curDepth)
	if err != nil {
		fmt.Println("json.Unmarshal failed:", err)
	}
	fmt.Println("curDepth:", curDepth.Bids)
	fmt.Println("curDepth:", len(curDepth.Bids))
	fmt.Println("curDepth:", curDepth.Asks)
	fmt.Println("curDepth:", len(curDepth.Asks))
}

func post() {
	rd := ioutil.NopCloser(strings.NewReader(url.Values{}.Encode()))
	fmt.Println("--------")
	url := "https://www.zg.com/api/v1/private/user"
	resp,err := http.Post(url,"multipart/form-data",rd)

	if err != nil {
		fmt.Println("http.Post failed err:",err)
		return
	}
	if resp.StatusCode == http.StatusOK {
		Doc, _ := goquery.NewDocumentFromReader(resp.Body)
		fmt.Println("DOC:",Doc.Text())
	}
	fmt.Println("err")
}

