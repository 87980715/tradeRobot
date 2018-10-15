package utils

import (
	"sort"
	"encoding/hex"
	"crypto/md5"
	"strings"
	"net/url"
	"net/http"
	"io/ioutil"
	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego/logs"
	"fmt"
)

type RestfulApiRequest struct {
	API_KEY    string
	SECRET_KEY string
	Sign       string

	ResPostDataLimit *PostDataLimit
	ResPostDataMarket *PostDataMarket
}

type PostDataLimit struct {
	Market     string //市场名称 交易对
	Side       string //1为ASK卖出，2为BID 买入
	Amount     string // 数量
	Price      string
}

type PostDataMarket struct {
	Market     string //市场名称 交易对
	Side       string //1为ASK卖出，2为BID 买入
}

const (
	API_URL      = "https://www.zg.com/api/v1/private/"
	Content_type = "application/x-www-form-urlencoded"
)

func (r *RestfulApiRequest) Md5Sign() {

	var data = make(map[string]string)
	// 判断账户参数，并以此判断账户将进行何种校验
	// 进行限价交易时需要设置的参数
	if r.ResPostDataLimit != nil {
		data["api_key"] = r.API_KEY
		//data["secret_key"] = r.SECRET_KEY
		data["market"] = r.ResPostDataLimit.Market
		data["side"] = r.ResPostDataLimit.Side
		data["amount"] = r.ResPostDataLimit.Amount

		data["price"] = r.ResPostDataLimit.Price
	}
	// 判断账户参数，并以此判断账户将进行何种校验
	// 进行市价交易时需要设置的参数
	if r.ResPostDataMarket != nil {
		data["api_key"] = r.API_KEY
		//data["secret_key"] = r.SECRET_KEY
		data["market"] = r.ResPostDataLimit.Market
		data["side"] = r.ResPostDataLimit.Side
	}

	// 查询账户时需要设置的参数
	data["api_key"] = r.API_KEY
	//data["secret_key"] = r.SECRET_KEY

	var signStr string
	tempSlice := make([]string, 0)
	for key := range data {
		tempSlice = append(tempSlice, key)
	}
	sort.Strings(tempSlice)
	for _, v := range tempSlice {
		signStr += v + "=" + data[v] + "&"
	}
	signStr = signStr + "secret_key=" + r.SECRET_KEY

	hash := md5.Sum([]byte(signStr))
	hashed := hash[:]
	sign := hex.EncodeToString(hashed)
	r.Sign = strings.ToUpper(sign)
}

// 获取用户资产
func (r *RestfulApiRequest) GetUserAssets() {
	v := url.Values{}
	v.Set("api_key", r.API_KEY)
	v.Set("secret_key", r.SECRET_KEY)
	v.Set("sign", r.Sign)
	rd := ioutil.NopCloser(strings.NewReader(v.Encode()))
	url := API_URL + "user"
	//fmt.Println("url:", url)
	resp, err := http.Post(url, Content_type, rd)
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
func (r *RestfulApiRequest) TradeLimit() {
	v := url.Values{}
	v.Set("api_key",r.API_KEY)
	v.Set("secret_key",r.SECRET_KEY)
	v.Set("sign",r.Sign)
	v.Set("market",r.ResPostDataLimit.Market)
	v.Set("side",r.ResPostDataLimit.Side)
	v.Set("amount",r.ResPostDataLimit.Amount)
	v.Set("price",r.ResPostDataLimit.Price)

	rd := ioutil.NopCloser(strings.NewReader(v.Encode()))
	url := API_URL + "trade/limit"
	resp, err := http.Post(url, Content_type, rd)
	if err != nil {
		logs.Error("http.Post tradeLimit failed err:", err)
		return
	}
	if resp.StatusCode == http.StatusOK {
		logs.Info("交易成功",resp.Body)
	}
	logs.Info("交易失败",resp.Body)
}

// 进行市价交易
func (r *RestfulApiRequest) TradeMarket() {
	v := url.Values{}
	v.Set("api_key",r.API_KEY)
	v.Set("secret_key",r.SECRET_KEY)
	v.Set("sign",r.Sign)
	v.Set("market",r.ResPostDataLimit.Market)
	v.Set("side",r.ResPostDataLimit.Side)

	rd := ioutil.NopCloser(strings.NewReader(v.Encode()))
	url := API_URL + "trade/market"
	resp, err := http.Post(url, Content_type, rd)
	if err != nil {
		logs.Error("http.Post tradeLimit failed err:", err)
		return
	}

	if resp.StatusCode == http.StatusOK {
		logs.Info("交易成功",resp.Body)
	}

	logs.Info("交易失败",resp.Body)
}