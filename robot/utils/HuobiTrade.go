package utils

import (
	"net/url"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"fmt"
	"tradeRobot/robot/models"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"sort"
	"time"
	"strings"
	"github.com/astaxie/beego/logs"
	"encoding/json"
	"strconv"
	"sync"
)

type HuobiRestfulApiRequest struct {
	PostDataLimit  *HuobiPostDataLimit
	PostDataMarket *HuobiPostDataMarket
	PostPataCancle *HuobiPostDataCancel
	GetDataPending *HuobiGetDataPending
	GetTradesDeal  *HuobiGetTradesDeal
}

type HuobiGetTradesDeal struct {
	Symbol string
	Zize   string
}

type HuobiTradesDealReturn struct {
	Status string       `json:"status"`
	Data   []*TradeDeal `json:"data"`
}

type TradeDeal struct {
	Id            int    `json:"id"`
	Symbol        string `json:"symbol"`
	Type          string `json:"type"`
	Price         string `json:"price"`
	Filled_amount string `json:"filled-amount"`
	Filled_fees   string `json:"filled-fees"`
	Created_at    int    `json:"created-at"`
}

type HuobiGetDataPending struct {
	Account_id string
	Symbol     string
	Size       int
}

type HuobiPostDataCancel struct {
	Order_id string
}

type HuobiPostDataLimit struct {
	Account_id string //账户 ID
	Amount     string //下单数量
	Price      string
	Symbol     string // 交易对 btcusdt
	Type       string // 市价卖, buy-limit：限价买, sell-limit
}

type HuobiPostDataMarket struct {
	Symbol string // 市场名称 交易对
	Type   string // 买卖类型：限价单(buy-limit/sell-limit) 市价单(buy_market/sell_market)
	Price  string // 下单总价格
}

type HuobiTradeResp struct {
	Result   bool `json:"result"`   // true代表成功返回
	Order_id int  `json:"order_id"` // 订单ID
}

type HuobiTradeCancelReturn struct {
	Result   bool  `json:"result"`
	Order_id int64 `json:"order_id"`
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

type HuobiLimitTradeReturn struct {
	Status string `json:"status"`
}

type HuobiPendingOrdersReturn struct {
	Status string                     `json:"status"`
	Data   []*PendingOrdersReturnData `json:"data"`
}

type PendingOrdersReturnData struct {
	Id                 int    `json:"id"`
	Symbol             string `json:"symbol"`
	Account_id         int    `json:"account-id"`
	Amount             string `json:"amount"`
	Price              string `json:"price"`
	Created_at         int    `json:"created-at"` // 下单时间（毫秒）
	Type               string `json:"type"`       // 订单类型
	Filled_amount      string `json:"filled-amount"`
	Filled_cash_amount string `json:"filled-cash-amount"`
	Filled_fees        string `json:"filled-fees"`
}

type HuobiCanleReturn struct {
	Status string `json:"status"`
}

var (
	HuobiOrders = make(chan *HuobiPostDataLimit, 100)
	RMuLock sync.RWMutex
)

// Huobi 获取用户资产
func (r *HuobiRestfulApiRequest) HuobiGetUserAssets() {
	signParams := make(map[string]string)
	signParams["AccessKeyId"] = models.Huobi_AccessKeyId
	signParams["SignatureVersion"] = "2"
	signParams["SignatureMethod"] = "HmacSHA256"
	signParams["Timestamp"] = time.Now().UTC().Format("2006-01-02T15:04:05")

	sign := HuobiSign(signParams, "GET", models.Huobi_API_URL, "/v1/account/accounts", models.Huobi_Secretkey)
	signParams["Signature"] = sign

	strUrl := "https://" + models.Huobi_API_URL + "/v1/account/accounts?" + Map2UrlQuery(MapValueEncodeURI(signParams))
	resp, err := http.Get(strUrl)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	fmt.Println(doc.Text())
}
var n int
// Huobi 进行限价交易
func (r *HuobiRestfulApiRequest) HuobiLimitTrade() {
	signParams := make(map[string]string)
	signParams["AccessKeyId"] = models.Huobi_AccessKeyId
	signParams["SignatureVersion"] = "2"
	signParams["SignatureMethod"] = "HmacSHA256"
	signParams["Timestamp"] = time.Now().UTC().Format("2006-01-02T15:04:05")

	sign := HuobiSign(signParams, "POST", models.Huobi_API_URL, "/v1/order/orders/place", models.Huobi_Secretkey)
	signParams["Signature"] = sign
	strUrl := "https://" + models.Huobi_API_URL + "/v1/order/orders/place?" + Map2UrlQuery(MapValueEncodeURI(signParams))

	mapParams := make(map[string]string)
	mapParams["account-id"] = r.PostDataLimit.Account_id
	mapParams["amount"] = r.PostDataLimit.Amount
	mapParams["price"] = r.PostDataLimit.Price
	mapParams["symbol"] = r.PostDataLimit.Symbol
	mapParams["type"] = r.PostDataLimit.Type
	// 测试
	// fmt.Println("mapParams:",mapParams)

	bytesParams, _ := json.Marshal(mapParams)
	jsonParams := string(bytesParams)

	rd := strings.NewReader(jsonParams)
	req, err := http.NewRequest("POST", strUrl, rd)
	if err != nil {
		logs.Error("http new request cancel order failed err:", err)
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept-Language", "zh-cn")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logs.Error("http.Post GetUserAssets failed err:", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		// 测试
		// fmt.Println(doc.Text())
		if err != nil {
			logs.Error(" go qurey new document from reader failed err:", err)
			return
		}
		var limitTradeReturn = &HuobiLimitTradeReturn{}
		err = json.Unmarshal([]byte(doc.Text()), limitTradeReturn)
		if err != nil {
			logs.Error(" go query limit trade return  from reader failed err:", err)
			return
		}
		n++
		for {
			if limitTradeReturn.Status == "ok" {
				logs.Info("火币交易挂单成功")
				break
			}
			if n > 5 {
				logs.Info("火币挂单失败")
				break
			}

		}
	}
}

// huobi查询已成交订单 并写入数据库
func (r *HuobiRestfulApiRequest) HuobiTradesDeal() {

	signParams := make(map[string]string)
	signParams["AccessKeyId"] = models.Huobi_AccessKeyId
	signParams["SignatureVersion"] = "2"
	signParams["SignatureMethod"] = "HmacSHA256"
	signParams["Timestamp"] = time.Now().UTC().Format("2006-01-02T15:04:05")

	signParams["symbol"] = r.GetTradesDeal.Symbol
	signParams["size"] = strconv.Itoa(models.Huobi_FilledOrdersSize)

	sign := HuobiSign(signParams, "GET", models.Huobi_API_URL, "/v1/order/matchresults", models.Huobi_Secretkey)
	signParams["Signature"] = sign
	strUrl := "https://" + models.Huobi_API_URL + "/v1/order/matchresults?" + Map2UrlQuery(MapValueEncodeURI(signParams))

	client := &http.Client{}
	req, _ := http.NewRequest("GET", strUrl, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		logs.Error("http get pending orders failed err:", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		size := models.Huobi_FilledOrdersSize
		var tradesDealReturn = &HuobiTradesDealReturn{
			Data: make([]*TradeDeal, size),
		}
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			logs.Error(" go query filled orders from reader failed err:", err)
			return
		}

		// -测试------------
		// fmt.Println("火币已成交数据:",doc.Text())
		err = json.Unmarshal([]byte(doc.Text()), tradesDealReturn)
		if err != nil {
			logs.Error(" json unmarshal filled orders failed err:", err)
			return
		}

		key := signParams["symbol"] + "HuobiDealId"

		for _, order := range tradesDealReturn.Data {
			if order.Id > int(models.HuoPreDealId[key]) {
				// fmt.Println(order.Id,HuobiDealId)
				var tradeResult = &models.HuobiTradeResults{}
				a, err := strconv.ParseFloat(order.Filled_amount, 64)
				if err != nil {
					logs.Error(" strconv parseFloat order filled_amount  failed err:", err)
					return
				}
				p, err := strconv.ParseFloat(order.Price, 64)
				if err != nil {
					logs.Error(" strconv parseFloat order price failed err:", err)
					return
				}
				t := a * p
				total := fmt.Sprintf("%."+strconv.Itoa(8)+"f", t)

				tradeResult.User_id = models.HuobiUserID
				tradeResult.Trade_id = strconv.Itoa(order.Id)
				tradeResult.Symbol = order.Symbol
				tradeResult.Type = order.Type
				tradeResult.Price = order.Price
				tradeResult.Deal_amount = order.Filled_amount
				tradeResult.Deal_fees = order.Filled_fees
				tradeResult.Created_at = strconv.Itoa(order.Created_at)
				tradeResult.Total = total

				db := RobotDB
				if err = db.Create(tradeResult).Error; err != nil {
					logs.Error("insert failed into Huobi tradeResult ")
					return
				}
			}
		}
		// 去重
		if len(tradesDealReturn.Data) != 0 {
			models.HuoPreDealId[key] = int64(tradesDealReturn.Data[0].Id)
		}
	}
}

// Huobi查询未成交订单并取消满足条件的订单
func (r *HuobiRestfulApiRequest) HuobiCancelPendingOrders() {
	signParams := make(map[string]string)
	signParams["AccessKeyId"] = models.Huobi_AccessKeyId
	signParams["SignatureVersion"] = "2"
	signParams["SignatureMethod"] = "HmacSHA256"
	signParams["Timestamp"] = time.Now().UTC().Format("2006-01-02T15:04:05")

	signParams["account-id"] = r.GetDataPending.Account_id
	signParams["symbol"] = r.GetDataPending.Symbol
	signParams["size"] = strconv.Itoa(models.Huobi_PendingOrdersSize)

	sign := HuobiSign(signParams, "GET", models.Huobi_API_URL, "/v1/order/openOrders", models.Huobi_Secretkey)
	signParams["Signature"] = sign
	strUrl := "https://" + models.Huobi_API_URL + "/v1/order/openOrders?" + Map2UrlQuery(MapValueEncodeURI(signParams))

	client := &http.Client{}
	req, _ := http.NewRequest("GET", strUrl, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		logs.Error("http get pending orders failed err:", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		size := models.Huobi_PendingOrdersSize
		var pendingOrdersReturn = &HuobiPendingOrdersReturn{
			Data: make([]*PendingOrdersReturnData, size),
		}
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			logs.Error(" go query pending orders from reader failed err:", err)
			return
		}
		// 测试--
		//fmt.Println("火币未成交数据:",doc.Text())
		err = json.Unmarshal([]byte(doc.Text()), pendingOrdersReturn)
		if err != nil {
			logs.Error(" json unmarshal pending orders failed err:", err)
			return
		}
		if len(pendingOrdersReturn.Data) != 0 {
			for _, order := range pendingOrdersReturn.Data {
				// 获取当前时间，毫秒 ms
				curTime := time.Now().UnixNano() / 1e6
				createTime := order.Created_at
				// 超过500ms,未成交
				//fmt.Println("未成交时间:",curTime-int64(createTime),"参数时间:",models.TradeInspectTime)
				if curTime-int64(createTime) > models.TradeInspectTime {
					// 取消订单
					orderId := strconv.Itoa(order.Id)
					if r.HuobiCancelOrder(orderId) {
						RMuLock.RLock()
						usdtPrice := models.UsdtPrice["huobi"]
						ethPrice := models.EthPrice["huobi"]
						RMuLock.RUnlock()
						postDataLimit := &HuobiPostDataLimit{}
						postDataLimit.Account_id = strconv.Itoa(order.Account_id)
						postDataLimit.Symbol = order.Symbol
						postDataLimit.Type = order.Type
						if order.Type == "buy-limit" {
							p, _ := strconv.ParseFloat(order.Price, 64)
							price := p * (1 + models.TradePriceAdjust) * usdtPrice * ethPrice
							// 4
							postDataLimit.Price = fmt.Sprintf("%."+strconv.Itoa(8)+"f", price)
						} else {
							// 价格设置：降低价格1‰，重新挂单
							p, _ := strconv.ParseFloat(order.Price, 64)
							price := p * (1 - models.TradePriceAdjust) * usdtPrice * ethPrice
							// 4
							postDataLimit.Price = fmt.Sprintf("%."+strconv.Itoa(8)+"f", price)
						}
						// 数量设置：减去已成交的数量
						amount, _ := strconv.ParseFloat(order.Amount, 64)
						filledAmount, _ := strconv.ParseFloat(order.Filled_amount, 64)
						postDataLimit.Amount = strconv.FormatFloat(amount - filledAmount, 'E', -1, 64)

						//fmt.Println("cancel:",postDataLimit)
						HuobiOrders <- postDataLimit
					}
				}
			}
		}
	}
}

// Huobi取消订单
func (r *HuobiRestfulApiRequest) HuobiCancelOrder(orderId string) bool {
	signParams := make(map[string]string)
	signParams["AccessKeyId"] = models.Huobi_AccessKeyId
	signParams["SignatureVersion"] = "2"
	signParams["SignatureMethod"] = "HmacSHA256"
	signParams["Timestamp"] = time.Now().UTC().Format("2006-01-02T15:04:05")
	strRequestPath := "/v1/order/orders/" + orderId + "/submitcancel"

	sign := HuobiSign(signParams, "POST", models.Huobi_API_URL, strRequestPath, models.Huobi_Secretkey)
	signParams["Signature"] = sign
	strUrl := "https://" + models.Huobi_API_URL + strRequestPath + "?" + Map2UrlQuery(MapValueEncodeURI(signParams))

	mapParams := make(map[string]string)
	mapParams["order-id"] = orderId

	bytesParams, _ := json.Marshal(mapParams)
	jsonParams := string(bytesParams)

	rd := strings.NewReader(jsonParams)
	client := &http.Client{}
	req, err := http.NewRequest("POST", strUrl, rd)
	if err != nil {
		logs.Error("http new request cancel order failed err:", err)
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept-Language", "zh-cn")

	var flag = true
	resp, err := client.Do(req)
	if err != nil {
		logs.Error("http.Post huobi cancel order failed err:", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			logs.Error(" go qurey new document from cancel huobi order failed err:", err)
			return false
		}
		var cancelReturn = &HuobiCanleReturn{}
		err = json.Unmarshal([]byte(doc.Text()), cancelReturn)
		if err != nil {
			logs.Error(" json unmarshal  cancelReturn failed err:", err)
			return false
		}

		fmt.Println("cancelReturn.Status",cancelReturn.Status)

		if cancelReturn.Status != "ok" {
			logs.Error("cancelReturn status is not ok ")
			return false
		}
		return true
	} else {
		flag = false
	}
	return flag
}

// Huobi加密
func HuobiSign(mapParams map[string]string, strMethod, strHostUrl, strRequestPath, strSecretKey string) string {
	// 参数处理, 按API要求, 参数名应按ASCII码进行排序(使用UTF-8编码, 其进行URI编码, 16进制字符必须大写)
	mapCloned := make(map[string]string)
	for key, value := range mapParams {
		mapCloned[key] = url.QueryEscape(value)
	}
	strParams := Map2UrlQueryBySort(mapCloned)
	strPayload := strMethod + "\n" + strHostUrl + "\n" + strRequestPath + "\n" + strParams
	return ComputeHmac256(strPayload, strSecretKey)
}

func ComputeHmac256(strMessage string, strSecret string) string {
	key := []byte(strSecret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(strMessage))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func Map2UrlQueryBySort(mapParams map[string]string) string {
	var keys []string
	for key := range mapParams {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var strParams string
	for _, key := range keys {
		strParams += key + "=" + mapParams[key] + "&"
	}
	if 0 < len(strParams) {
		strParams = string([]rune(strParams)[:len(strParams)-1])
	}
	return strParams
}

func MapValueEncodeURI(mapValue map[string]string) map[string]string {
	for key, value := range mapValue {
		valueEncodeURI := url.QueryEscape(value)
		mapValue[key] = valueEncodeURI
	}
	return mapValue
}

func Map2UrlQuery(mapParams map[string]string) string {
	var strParams string
	for key, value := range mapParams {
		strParams += key + "=" + value + "&"
	}
	if 0 < len(strParams) {
		strParams = string([]rune(strParams)[:len(strParams)-1])
	}
	return strParams
}
