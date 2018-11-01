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
	data["AccessKeyId"] = "eca1800e-e94af9ce-c2a77a7b-7a8b4"
	data["SignatureVersion"] = "2"
	data["SignatureMethod"] = "HmacSHA256"
	data["Timestamp"] = time.Now().UTC().Format("2006-01-02T15:04:05")

	sign := utils.HuobiSign(data, "GET", models.Huobi_API_URL, "/v1/account/accounts", "5658bce4-10643a8d-7e62938f-24139")
	data["Signature"] = sign

	strUrl := "https://" + models.Huobi_API_URL + "/v1/account/accounts?" + utils.Map2UrlQuery(utils.MapValueEncodeURI(data))
	fmt.Println(strUrl)
	resp, err := http.Get(strUrl)
	if err != nil {
		fmt.Println("err:", err)
	}
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	fmt.Println(doc.Text())

}
