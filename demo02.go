package main

import (
	"net/http"
	"fmt"
)

func main() {
	f1()
}

func f1() {

	resp,err := http.Get("https://www.huobi.com/zh-cn/mt_eth/exchange/")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
}