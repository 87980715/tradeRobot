package utils

import (
	"strings"
	"net/url"
	"github.com/astaxie/beego/logs"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"fmt"
	"encoding/json"
	"tradeRobot/robot/models"
	"time"
)

type OKexRestfulApiRequest struct {
	API_KEY    string
	SECRET_KEY string
	Sign       string

	PostDataLimit  *OKexPostDataLimit
	PostDataMarket *OKexPostDataMarket
	PostPataCancle *OKexPostDataCancel
}

type OKexPostDataCancel struct {
	Symbol   string
	Order_id string
}

type OKexPostDataLimit struct {
	Instrument_id string // 市场名称 交易对
	Type          string // 买卖类型：限价单(buy/sell) 市价单(buy_market/sell_market)
	Amount        string // 数量
	Price         string
	Side          string
}

type OKexPostDataMarket struct {
	Symbol string // 市场名称 交易对
	Type   string // 买卖类型：限价单(buy/sell) 市价单(buy_market/sell_market)
	Price  string // 下单总价格
}

type TradeResp struct {
	Result   bool `json:"result"`   // true代表成功返回
	Order_id int  `json:"order_id"` // 订单ID
}

type TradeCancelResp struct {
	Result   bool `json:"result"`
	Order_id int  `json:"order_id"`
}

type OrderSInfoResp struct {
	Result bool         `json:"result"`
	Orders []*OrderInfo `json:"orders"`
}

type OrderInfo struct {
	Amount      float64 `json:"amount"`
	Avg_price   float64 `json:"avg_price"`
	Create_date float64 `json:"create_date"`
	Deal_amount float64 `json:"deal_amount"`
	Order_id    float64 `json:"order_id"`
	Orders_id   float64 `json:"orders_id"`
	Price       float64 `json:"price"`
	Status      float64 `json:"status"`
	Symbol      string  `json:"symbol"`
	Type        string  `json:"type"`
}

type timeStamp struct {
	Iso   string
	epoch string
}

var OKexOrders = make(chan *OKexPostDataLimit,100)

func (r *OKexRestfulApiRequest) OKexQureyMd5Sign() {
	var data = make(map[string]string)

	// 查询账户时需要设置的参数
	data["api_key"] = r.API_KEY
	//data["X-SITE-ID"] = "1"
	r.Sign = ZGSign(data, r.SECRET_KEY)
}

func (r *OKexRestfulApiRequest) OKexMarketMd5Sign() {
	var data = make(map[string]string)
	if r.PostDataMarket != nil {
		data["api_key"] = r.API_KEY
		data["Symbol"] = r.PostDataMarket.Symbol
		data["type"] = r.PostDataMarket.Type
	}
	r.Sign = ZGSign(data, r.SECRET_KEY)
}

func (r *OKexRestfulApiRequest) OKexLimitMd5Sign() {
	var data = make(map[string]string)
	// 判断账户参数，并以此判断账户将进行何种校验
	// 进行限价交易时需要设置的参数
	if r.PostDataLimit != nil {
		data["api_key"] = r.API_KEY
		//data["symbol"] = r.PostDataLimit.Symbol
		data["type"] = r.PostDataLimit.Type
		data["amount"] = r.PostDataLimit.Amount
		data["price"] = r.PostDataLimit.Price
	}
	r.Sign = ZGSign(data, r.SECRET_KEY)
}


// OKex 获取用户资产
func (r *OKexRestfulApiRequest) OKexGetUserAssets() {
	v := url.Values{}
	v.Set("api_key", r.API_KEY)
	v.Set("sign", r.Sign)
	rd := ioutil.NopCloser(strings.NewReader(v.Encode()))

	url := models.Okex_API_URL + "userinfo.do"

	resp, err := http.Post(url, models.Okex_Content_type, rd)
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

// OKex 进行限价交易
func (r *OKexRestfulApiRequest) OKexTrade(time time.Duration, ) {
	v := url.Values{}
	v.Set("instrument_id", r.PostDataLimit.Instrument_id)
	v.Set("type", r.PostDataLimit.Type)
	v.Set("price", r.PostDataLimit.Price)
	v.Set("size", r.PostDataLimit.Amount)
	v.Set("sign", r.Sign)

	rd := ioutil.NopCloser(strings.NewReader(v.Encode()))
	url := models.Okex_API_URL + "trade.do"
	resp, err := http.Post(url, models.Okex_Content_type, rd)
	if err != nil {
		logs.Error("http.Post OKex tradeLimit failed err:", err)
		return
	}
	// todo 需要继续完善，拿到订单号，用于之后撤销撤销定单
	if resp.StatusCode == http.StatusOK {
		logs.Info("交易成功", resp.Body)
	}
	logs.Info("交易失败", resp.Body)
}

// OKex 进行市价交易
// 市价卖单不传price
func (r *OKexRestfulApiRequest) OKexTradeMarket() {
	v := url.Values{}
	v.Set("api_key", r.API_KEY)
	v.Set("symbol", r.PostDataMarket.Symbol)
	v.Set("type", r.PostDataMarket.Type)
	// 交易数量 市价买单不传amount,市价买单需传price作为买入总金额
	v.Set("price", r.PostDataMarket.Price)
	v.Set("sign", r.Sign)

	rd := ioutil.NopCloser(strings.NewReader(v.Encode()))
	url := models.Okex_API_URL + "trade.do"
	resp, err := http.Post(url, models.Okex_Content_type, rd)
	if err != nil {
		logs.Error("http.Post tradeMarket failed err:", err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	var CurCancelResp = &TradeCancelResp{}
	err = json.Unmarshal([]byte(doc.Text()), CurCancelResp)
	if err != nil {
		logs.Error("json.Unmarshal MarketResp failed")
	}

	// todo 需要继续完善，拿到订单号，用于之后撤销撤销定单
	if resp.StatusCode == http.StatusOK {
		logs.Info("交易成功", resp.Body)
	}

	logs.Info("交易失败", resp.Body)
}

// 撤销订单
func (r *OKexRestfulApiRequest) OKexTradeCancel() bool {
	v := url.Values{}
	v.Set("api_key", r.API_KEY)
	v.Set("symbol", r.PostPataCancle.Symbol)
	v.Set("order_id", r.PostPataCancle.Order_id)
	v.Set("sign", r.Sign)

	rd := ioutil.NopCloser(strings.NewReader(v.Encode()))
	url := models.Okex_API_URL + "cancel_order.do"
	resp, err := http.Post(url, models.Okex_Content_type, rd)
	if err != nil {
		logs.Error("http.Post cancel trade failed err:", err)
		return false
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	var curCancelResp = &TradeCancelResp{Result: true}
	err = json.Unmarshal([]byte(doc.Text()), curCancelResp)
	if err != nil {
		logs.Error("json.Unmarshal CurCancelResp failed")
		return false
	}
	if !curCancelResp.Result {
		logs.Info("cancel order failed")
	}
	return true
}
