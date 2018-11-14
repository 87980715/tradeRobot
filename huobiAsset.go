package main

import (
	"tradeRobot/robot/utils"
	"fmt"
	"time"
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"tradeRobot/robot/models"
)

func main() {

	data := make(map[string]string)
	data["AccessKeyId"] = "b2af8f9f-4ac75b4d-4fce0763-1c789"
	data["SignatureVersion"] = "2"
	data["SignatureMethod"] = "HmacSHA256"
	data["Timestamp"] = time.Now().UTC().Format("2006-01-02T15:04:05")

	sign := utils.HuobiSign(data, "GET", models.Huobi_API_URL, "/v1/account/accounts/4821321/balance", "80653212-f2b1fc55-f1af7577-b9a3f")
	data["Signature"] = sign

	strUrl := "https://" + models.Huobi_API_URL + "/v1/account/accounts/4821321/balance?" + utils.Map2UrlQuery(utils.MapValueEncodeURI(data))
	fmt.Println(strUrl)
	resp, err := http.Get(strUrl)
	if err != nil {
		fmt.Println("err:", err)
	}
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	fmt.Println(doc.Text())
}
