package utils

import (
	"strings"
	"net/url"
	"net/http"
	"io/ioutil"
	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego/logs"
	"sort"
	"crypto/md5"
	"encoding/hex"
	"tradeRobot/robot/models"
	"encoding/json"
	"time"
	"strconv"
	"math/rand"
)

type ZTRestfulApiRequest struct {
	API_KEY    string
	SECRET_KEY string
	Sign       string

	PostDataQueryPending  *ZTPostDataQureyPending
	PostDataCancel        *ZTPostDataCancel
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
	Code    int     `json:"code"`
	Message string  `json:"message"`
	Result  *Result `json:"result"`
}

type Result struct {
	Limit   int       `json:"limit"`
	Offset  int       `json:"offset"`
	Records []*Record `json:"records"`
}

type Record struct {
	Amount     string `json:"amount"`
	Ctime      float64 `json:"ctime"` // create time
	Deal_fee   string `json:"deal_fee"`
	Deal_money string `json:"deal_money"`
	Deal_stock string `json:"deal_stock"`
	Ftime      float64 `json:"ftime"`
	Id         int64 `json:"id"`
	Maker_fee  string `json:"maker_fee"`
	Market     string `json:"market"`
	Price      string `json:"price"`
	Side       int `json:"side"`
	Source     string `json:"source"`
	Taker_fee  string `json:"taker_fee"`
	Type       int `json:"type"`
	User       int64 `json:"user"`
}

type ZTPostDataQureyPending struct {
	Market string
	Offset int
	Limit  int
}

type PendingResult struct {
	Records []PeningRecord `json:"records"`
}

type PendingOrder struct {
	Result PendingResult `json:"result"`
}

type PeningRecord struct {
	Ctime  float64 `json:"ctime"`
	Id     int64   `json:"id"`
	Market string  `json:"market"`
}

type ZTPostDataCancel struct {
	Market   string
	Order_id int64
}

func (r *ZTRestfulApiRequest) ZTCancelMd5Sign() {

	var data = make(map[string]string)
	data["api_key"] = r.API_KEY
	data["market"] = r.PostDataCancel.Market
	data["order_id"] = strconv.Itoa(int(r.PostDataCancel.Order_id))
	r.Sign = ZGSign(data, r.SECRET_KEY)
}

func (r *ZTRestfulApiRequest) ZTQueryPendingMd5Sign() {

	var data = make(map[string]string)
	data["api_key"] = r.API_KEY
	data["market"] = r.PostDataQueryPending.Market
	data["limit"] = strconv.Itoa(r.PostDataQueryPending.Limit)
	data["offset"] = strconv.Itoa(r.PostDataQueryPending.Offset)
	//fmt.Println("data",data)
	r.Sign = ZGSign(data, r.SECRET_KEY)
}

func (r *ZTRestfulApiRequest) ZTQueryDealMd5Sign() {

	var data = make(map[string]string)
	data["api_key"] = r.API_KEY
	data["market"] = r.PostDataOrderFinished.Market
	data["limit"] = r.PostDataOrderFinished.Limit
	data["offset"] = r.PostDataOrderFinished.Offset
	data["start_time"] = r.PostDataOrderFinished.Start_time
	data["end_time"] = r.PostDataOrderFinished.End_time
	data["side"] = r.PostDataOrderFinished.Side

	r.Sign = ZGSign(data, r.SECRET_KEY)
}

func (r *ZTRestfulApiRequest) ZTLimitMd5Sign() {

	var data = make(map[string]string)
	data["api_key"] = r.API_KEY
	data["market"] = r.PostDataLimit.Market
	data["side"] = r.PostDataLimit.Side
	data["amount"] = r.PostDataLimit.Amount
	data["price"] = r.PostDataLimit.Price
	r.Sign = ZGSign(data, r.SECRET_KEY)
}

func (r *ZTRestfulApiRequest) ZTMarketMd5Sign() {

	var data = make(map[string]string)
	data["api_key"] = r.API_KEY
	data["market"] = r.PostDataMarket.Market
	data["side"] = r.PostDataMarket.Side
	r.Sign = ZGSign(data, r.SECRET_KEY)
}

func (r *ZTRestfulApiRequest) ZTQueyMd5Sign() {

	var data = make(map[string]string)
	data["api_key"] = r.API_KEY

	r.Sign = ZGSign(data, r.SECRET_KEY)
}

// 获取用户资产
func (r *ZTRestfulApiRequest) ZTGetUserAssets() (error,string) {
	v := url.Values{}
	v.Set("api_key", r.API_KEY)
	v.Set("secret_key", r.SECRET_KEY)
	v.Set("sign", r.Sign)
	rd := ioutil.NopCloser(strings.NewReader(v.Encode()))
	assetsUrl := models.ZG_API_URL + "user"

	resp, err := http.Post(assetsUrl, models.ZG_Content_type, rd)
	if err != nil {
		logs.Error("http.Post GetUserAssets failed err:", err)
		return err ,""
	}
	defer resp.Body.Close()

	Doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		logs.Error("goquery.NewDocumentFromReader failed err:", err)
		return err ,""
	}
	if resp.StatusCode == http.StatusOK {
		return nil,Doc.Text()
	} else {
		return err ,""
	}
}

var tradeNum int64
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
		if resp != nil {
			resp.Body.Close()
		}
		logs.Error("http.Post tradeLimit failed err:", err)
		return
	}

	_, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		logs.Error("goquery.NewDocumentFromReader failed err:", err)
		return
	}
	logs.Info("ZG挂单成功")
	tradeNum ++
	logs.Info("tradeNum:", tradeNum)
	if resp != nil {
		resp.Body.Close()
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

	if resp != nil {
		defer resp.Body.Close()
	}

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
	}
	// resp.Body.Close()
}

// 查询已成交订单
func (r *ZTRestfulApiRequest) ZTOrderFinished() (string) {

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
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		logs.Error("http post finished order failed err:", err)
		return ""
	}

	if resp.StatusCode == http.StatusOK {
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			logs.Error("goquery.NewDocumentFromReader failed err:", err)
		}
		return doc.Text()

	}
	logs.Error("查询用户已成交订单失败")
	return ""
}

func (r *ZTRestfulApiRequest) ZTQueryPending() []*ZTPostDataCancel {
	v := url.Values{}
	v.Set("market", r.PostDataQueryPending.Market)
	v.Set("offset", strconv.Itoa(r.PostDataQueryPending.Offset))
	v.Set("limit", strconv.Itoa(r.PostDataQueryPending.Limit))

	v.Set("sign", r.Sign)
	v.Set("secret_key", r.SECRET_KEY)
	v.Set("api_key", r.API_KEY)

	rd := ioutil.NopCloser(strings.NewReader(v.Encode()))
	OrdersUrl := models.ZG_API_URL + "order/pending"
	resp, err := http.Post(OrdersUrl, models.ZG_Content_type, rd)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		logs.Error("http.Post zt query pending failed err:", err)
		return nil
	}

	// ----测试----
	// fmt.Println("resp.StatusCode:", resp.StatusCode)
	if resp.StatusCode == http.StatusOK {
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		var pendingOrders = &PendingOrder{
			Result: PendingResult{
				Records: make([]PeningRecord, r.PostDataQueryPending.Limit),
			},
		}
		err = json.Unmarshal([]byte(doc.Text()), pendingOrders)
		if err != nil {
			logs.Info("unmarshal zt query pending failed err:", err)
		}

		var cancelOrders []*ZTPostDataCancel
		rand.Seed(time.Now().UnixNano())
		randNum := float64(rand.Intn(5) + 10)

		if pendingOrders.Result.Records != nil {
			for _, record := range pendingOrders.Result.Records {
				curTime := float64(time.Now().Unix())
				if curTime-record.Ctime > randNum {
					var postData = &ZTPostDataCancel{}
					postData.Market = record.Market
					postData.Order_id = record.Id
					cancelOrders = append(cancelOrders, postData)
				}
			}
		}
		logs.Info("查询未成交订单成功")
		return cancelOrders
	}
	time.Sleep(3 * time.Second)
	logs.Error("查询未成交订单失败")
	return nil
}

var cancelNum int64
// 取消订单
func (r *ZTRestfulApiRequest) ZTCancelOrder() {

	v := url.Values{}
	v.Set("market", r.PostDataCancel.Market)
	v.Set("order_id", strconv.Itoa(int(r.PostDataCancel.Order_id)))

	v.Set("sign", r.Sign)
	v.Set("secret_key", r.SECRET_KEY)
	v.Set("api_key", r.API_KEY)

	rd := ioutil.NopCloser(strings.NewReader(v.Encode()))
	OrdersUrl := models.ZG_API_URL + "trade/cancel"
	resp, err := http.Post(OrdersUrl, models.ZG_Content_type, rd)
	if err != nil {
		logs.Error("http.Post cancel order failed err:", err)
		return
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp.StatusCode == http.StatusOK {
		_, err = goquery.NewDocumentFromReader(resp.Body)
		logs.Info("取消订单成功")
		cancelNum ++
		logs.Info("cancelNum:", cancelNum)
		resp.Body.Close()
		return
	}
	logs.Error("取消订单失败")
	//resp.Body.Close()
	return
}

func ZGSign(data map[string]string, secretKey string) (sign string) {

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

	hash := md5.Sum([]byte(signStr))
	hashed := hash[:]
	sign = strings.ToUpper(hex.EncodeToString(hashed))
	return
}
