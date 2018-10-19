package utils

import (
	"net/url"
	"strings"
	"github.com/astaxie/beego/logs"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"fmt"
	"tradeRobot/models"
)

type HuobiRestfulApiRequest struct {
	API_KEY    string
	SECRET_KEY string
	Sign       string

	PostDataLimit  *HuobiPostDataLimit
	PostDataMarket *HuobiPostDataMarket
	PostPataCancle *HuobiPostDataCancel
}

type HuobiPostDataCancel struct {
	Symbol   string
	Order_id string
}

type HuobiPostDataLimit struct {
	Symbol string // 市场名称 交易对
	Type   string // 买卖类型：限价单(buy/sell) 市价单(buy_market/sell_market)
	Amount string // 数量
	Price  string
}

type HuobiPostDataMarket struct {
	Symbol string // 市场名称 交易对
	Type   string // 买卖类型：限价单(buy/sell) 市价单(buy_market/sell_market)
	Price  string // 下单总价格
}

type HuobiTradeResp struct {
	Result   bool `json:"result"`   // true代表成功返回
	Order_id int  `json:"order_id"` // 订单ID
}

type HuobiTradeCancelResp struct {
	Result   bool `json:"result"`
	Order_id int  `json:"order_id"`
}

type HuobiOrderSInfoResp struct {
	Result bool         `json:"result"`
	Orders []*OrderInfo `json:"orders"`
}

type HuobiOrderInfo struct {
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

type HuobiTimeStamp struct {
	Iso   string
	epoch string
}

// Huobi 获取用户资产
func (r *HuobiRestfulApiRequest) HuobiGetUserAssets() {
	v := url.Values{}
	v.Set("api_key", r.API_KEY)
	v.Set("sign", r.Sign)
	rd := ioutil.NopCloser(strings.NewReader(v.Encode()))

	url := models.Huobi_API_URL + "userinfo.do"
	resp, err := http.Post(url, models.Huobi_Content_type, rd)
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

func sign() {

}