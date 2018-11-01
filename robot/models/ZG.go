package models

var (
	ZGUserID string

	//ZG_SECRET_KEY = "d8fedb73-17d1-4a94-a6f6-f45bc4dd79a0"
	ZG_SECRET_KEY = "8XUkytGxGH2ctvTHhrylBzho1lEz6YrS7AM8"
	//ZG_API_KEY = "ed5d8197-26db-45be-b1ce-719f13847b6c"
	ZG_API_KEY = "ZBC7CAodYDdHYxSI1h8xLqyqTxnV3dt92dpr"

	ZG_Content_type = "application/x-www-form-urlencoded"
	//ZG_API_URL = "https://www.zg.com/api/v1/private/"
	ZG_API_URL = "http://47.99.74.117:8883/api/v1/private/"
    ZGTradesChan = make(chan string, 100)
    ZGDealChan = make(chan string, 100)
)
