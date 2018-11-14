package main

import (
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"fmt"
	"strings"
	"time"
)

var url = "https://otc.huobi.br.com/zh-cn/trade/buy-usdt/"

func main() {
	GetUsdtPrice()
}

func GetUsdtPrice() {
	for {

		fmt.Println("spider run")
		opts := []selenium.ServiceOption{}

		caps := selenium.Capabilities{
			"browserName": "chrome",
		}
		imagCaps := map[string]interface{}{
			"profile.managed_default_content_settings.images": 2,
		}
		chromeCaps := chrome.Capabilities{
			Prefs: imagCaps,
			Path:  "",
			Args: []string{
				"--headless", // 设置Chrome无头模式,linux 必须设置，否则会报错
				"--no-sandbox",
				"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36", // 模拟user-agent，防反爬
			},
		}

		caps.AddChrome(chromeCaps)

		// Users/lumingjian/Downloads/chromedriver
		service, err := selenium.NewChromeDriverService("/opt/chrome//chromedriver", 9515, opts...)
		if err != nil {
			fmt.Printf("Error starting the ChromeDriver server: %v", err)
		}

		defer service.Stop()

		webDriver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", 9515))
		if err != nil {
			fmt.Println("selenium new remote failed err: ",err)
			continue
		}

		webDriver.Refresh()

		err = webDriver.Get(url)
		if err != nil {
			fmt.Println("webdriver get failed err: ",err)
			continue
		}

		t, err := webDriver.FindElement(selenium.ByXPATH, `//*[@id="tickerCny_ticker_bar"]`)
		if err != nil {
			//continue
		}

		fmt.Println(t)
		if t != nil {
			price, _ := t.Text()

			s := strings.Split(price, " ")

			fmt.Println(s[0])
		}
		time.Sleep(time.Minute)
	}
}