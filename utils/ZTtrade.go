package utils

import (
	"strings"
	"net/url"
	"net/http"
	"io/ioutil"
	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego/logs"
	"fmt"
	"sort"
	"crypto/md5"
	"encoding/hex"
	"tradeRobot/models"
)

type ZTRestfulApiRequest struct {
	API_KEY    string
	SECRET_KEY string
	Sign       string

	PostDataOrderFinished *ZTPostDataOrderFinished
	PostDataLimit         *ZTPostDataLimit
	PostDataMarket        *ZTPostDataMarket
}

type ZTPostDataLimit struct {
	Market string //市场名称 交易对
	Side   string //1为ASK卖出，2为BID 买入
	Amount string // 数量
	Price  string
}

type ZTPostDataMarket struct {
	Market string //市场名称 交易对
	Side   string //1为ASK卖出，2为BID 买入
	Amount string
}

type ZTPostDataOrderFinished struct {
	Market     string
	Start_time string
	End_time   string
	Offset     string
	Limit      string
	Side       string
}

type ZTOrderFinishedResp struct {
	Code    int
	Message string
	Result  *Result
}

type Result struct {
	Limit   int
	Offset  int
	Records []*Record
}

type Record struct {
	Amount     string
	Ctime      float64 // create time
	Deal_fee   string
	Deal_money string
	Deal_stock string
	Ftime      float64
	Id         int64
	Maker_fee  string
	Market     string
	Price      string
	Side       int
	Source     string
	Taker_fee  string
	Type       int
	User       int64
}


func (r *ZTRestfulApiRequest) ZTLimitMd5Sign() {

	var data = make(map[string]string)
	// 判断账户参数，并以此判断账户将进行何种校验
	// 进行限价交易时需要设置的参数
	if r.PostDataLimit != nil {
		data["api_key"] = r.API_KEY
		data["market"] = r.PostDataLimit.Market
		data["side"] = r.PostDataLimit.Side
		data["amount"] = r.PostDataLimit.Amount
		data["price"] = r.PostDataLimit.Price
	}
	r.Sign = Sign(data, r.SECRET_KEY)
}

func (r *ZTRestfulApiRequest) ZTMarketMd5Sign() {

	var data = make(map[string]string)

	// 进行市价交易时需要设置的参数
	if r.PostDataMarket != nil {
		data["api_key"] = r.API_KEY
		data["market"] = r.PostDataMarket.Market
		data["side"] = r.PostDataMarket.Side
	}
	r.Sign = Sign(data, r.SECRET_KEY)
}

func (r *ZTRestfulApiRequest) ZTQueyMd5Sign() {

	var data = make(map[string]string)

	// 查询账户时需要设置的参数
	data["api_key"] = r.API_KEY

	r.Sign = Sign(data, r.SECRET_KEY)
}

// 获取用户资产
func (r *ZTRestfulApiRequest) ZTGetUserAssets() {
	v := url.Values{}
	v.Set("api_key", r.API_KEY)
	v.Set("secret_key", r.SECRET_KEY)
	v.Set("sign", r.Sign)
	rd := ioutil.NopCloser(strings.NewReader(v.Encode()))
	assetsUrl := models.ZG_API_URL + "user"

	resp, err := http.Post(assetsUrl, models.ZG_Content_type, rd)
	if err != nil {
		logs.Error("http.Post GetUserAssets failed err:", err)
		return
	}
	defer resp.Body.Close()

	Doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logs.Error("goquery.NewDocumentFromReader failed err:", err)
		return
	}
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Doc.Text:", Doc.Text())
	} else {
		fmt.Println("Doc.Text:", Doc.Text())
	}
}

// 进行限价交易
func (r *ZTRestfulApiRequest) ZTTradeLimit() {
	v := url.Values{}
	v.Set("api_key", r.API_KEY)
	v.Set("secret_key", r.SECRET_KEY)
	v.Set("sign", r.Sign)
	v.Set("market", r.PostDataLimit.Market)
	v.Set("side", r.PostDataLimit.Side)
	v.Set("amount", r.PostDataLimit.Amount)
	v.Set("price", r.PostDataLimit.Price)

	rd := ioutil.NopCloser(strings.NewReader(v.Encode()))
	limitUrl := models.ZG_API_URL + "trade/limit"
	resp, err := http.Post(limitUrl, models.ZG_Content_type, rd)
	if err != nil {
		logs.Error("http.Post tradeLimit failed err:", err)
		return
	}
	if resp.StatusCode == http.StatusOK {

		_, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			logs.Error("goquery.NewDocumentFromReader failed err:", err)
			return
		}
		logs.Info("交易成功")
		//fmt.Println(doc.Text())
	}
}

// 进行市价交易
func (r *ZTRestfulApiRequest) ZTTradeMarket() {
	v := url.Values{}
	v.Set("api_key", r.API_KEY)
	v.Set("secret_key", r.SECRET_KEY)
	v.Set("sign", r.Sign)
	v.Set("market", r.PostDataLimit.Market)
	v.Set("side", r.PostDataLimit.Side)
	v.Set("amount", r.PostDataLimit.Amount)

	rd := ioutil.NopCloser(strings.NewReader(v.Encode()))
	marketUrl := models.ZG_API_URL + "trade/market"
	resp, err := http.Post(marketUrl, models.ZG_Content_type, rd)
	if err != nil {
		logs.Error("http.Post tradeLimit failed err:", err)
		return
	}
	if resp.StatusCode == http.StatusOK {
		_, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			logs.Error("goquery.NewDocumentFromReader failed err:", err)
			return
		}
		logs.Info("市价交易成功")
		//logs.Info(doc.Text())
	}
}

// 查询已成交订单
func (r *ZTRestfulApiRequest) ZTOrderFinished() (string,error) {

	v := url.Values{}
	v.Set("market", r.PostDataOrderFinished.Market)
	v.Set("start_time", r.PostDataOrderFinished.Start_time)
	v.Set("end_time", r.PostDataOrderFinished.End_time)
	v.Set("offset", r.PostDataOrderFinished.Offset)
	v.Set("limit", r.PostDataOrderFinished.Limit)
	v.Set("side", r.PostDataOrderFinished.Side)

	v.Set("sign", r.Sign)
	v.Set("secret_key", r.SECRET_KEY)
	v.Set("api_key", r.API_KEY)

	rd := ioutil.NopCloser(strings.NewReader(v.Encode()))
	OrdersUrl := models.ZG_API_URL + "order/finished"
	resp, err := http.Post(OrdersUrl, models.ZG_Content_type, rd)
	if err != nil {
		logs.Error("http.Post tradeLimit failed err:", err)
		return "",err
	}

	if resp.StatusCode != http.StatusOK {
		logs.Info("query finished orders failed", resp.Body)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logs.Error("goquery.NewDocumentFromReader failed err:", err)
		return "",err
	}

	return doc.Text(),nil
}


func Sign(data map[string]string, secretKey string) (sign string) {

	var signStr string
	tempSlice := make([]string, 0)
	for key := range data {
		tempSlice = append(tempSlice, key)
	}
	sort.Strings(tempSlice)
	for _, v := range tempSlice {
		signStr += v + "=" + data[v] + "&"
	}
	signStr = signStr + "secret_key=" + secretKey

	// fmt.Println("signStr:",signStr)

	hash := md5.Sum([]byte(signStr))
	hashed := hash[:]
	sign = strings.ToUpper(hex.EncodeToString(hashed))

	return
}
