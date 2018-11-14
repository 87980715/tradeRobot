package initialize

import (
	"time"
	"github.com/PuerkitoBio/goquery"
	"tradeRobot/robot/models"
	"tradeRobot/robot/utils"
	"net/http"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"strings"
	"net/url"
	"io/ioutil"
)

type HuobiAccountReturn struct {
	Data []*Data `json:"data"`
}

type Data struct {
	Id int `json:"id"`
}

type ZGAccountReturn struct {
	Code int `json:"code"`
}

type UserAssets struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Result  *Result `json:"result"`

}

type Result struct {
	User_id int `json:"user_id"`
}
func HuobiUserId() (userId int, err error) {

	data := make(map[string]string)
	data["AccessKeyId"] = models.Huobi_AccessKeyId
	data["SignatureVersion"] = "2"
	data["SignatureMethod"] = "HmacSHA256"
	data["Timestamp"] = time.Now().UTC().Format("2006-01-02T15:04:05")

	sign := utils.HuobiSign(data, "GET", models.Huobi_API_URL, "/v1/account/accounts", models.Huobi_Secretkey)
	data["Signature"] = sign

	strUrl := "https://" + models.Huobi_API_URL + "/v1/account/accounts?" + utils.Map2UrlQuery(utils.MapValueEncodeURI(data))

	resp, err := http.Get(strUrl)
	if err != nil {
		logs.Error("err:", err)
		return
	}
	var accountReturn = &HuobiAccountReturn{
		Data: make([]*Data,1),
	}

	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	err = json.Unmarshal([]byte(doc.Text()), accountReturn)
	if err != nil {
		logs.Error("err:", err)
		return
	}
	return accountReturn.Data[0].Id, nil
}

func ZGUserId() (userId int, err error) {

	account := utils.ZTAccount
	account.API_KEY = models.ZT_API_KEY
	account.SECRET_KEY = models.ZT_SECRET_KEY
	account.ZTQueyMd5Sign()

	err, data := account.ZTGetUserAssets()
	if err != nil {
		logs.Error("get user assets failed err:", err)
		return
	}

	var userAssets = &UserAssets{}
	err = json.Unmarshal([]byte(data), userAssets)
	if err != nil {
		logs.Error("json.Unmarshal user assets failed err:", err)
		return
	}
	userId = userAssets.Result.User_id

	return
}

func VerfiZGKey() (bool) {
	var flag = true
	var data = make(map[string]string)
	data["api_key"] = models.ZT_API_KEY
	sign := utils.ZGSign(data, models.ZT_SECRET_KEY)

	v := url.Values{}
	v.Set("api_key", models.ZT_API_KEY)
	v.Set("secret_key", models.ZT_SECRET_KEY)
	v.Set("sign", sign)
	rd := ioutil.NopCloser(strings.NewReader(v.Encode()))
	assetsUrl := models.ZG_API_URL + "user"

	resp, err := http.Post(assetsUrl, models.ZG_Content_type, rd)
	if err != nil {
		logs.Error("http.Post GetUserAssets failed err:", err)
		return false
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logs.Error("goquery.NewDocumentFromReader failed err:", err)
		return false
	}

	var accountReturn = &ZGAccountReturn{}
	err = json.Unmarshal([]byte(doc.Text()), accountReturn)
	if err != nil {
		logs.Error("err:", err)
		return false
	}

	if accountReturn.Code != 0 {
		flag = false
	}
	return flag
}

