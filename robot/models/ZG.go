package models

var (

	UserID = make(map[string]int64)

	ZT_API_KEY string   = "xtx84xh65386czj6zyw3hwvsx0r3jzkd6amb"
	ZT_SECRET_KEY string = "yvv5vmrs5omdaliyidoy8baazqumlryxo1u7"

	ZG_Content_type = "application/x-www-form-urlencoded"

	ZG_API_URL = "https://www.zt.com/api/v1/private/"

    ZGTradesChan = make(chan string, 100)

    ZGQueryDealOrderSize = "20"

    UsdtCntRate float64
)


// test API
// ZG_SECRET_KEY = "d8fedb73-17d1-4a94-a6f6-f45bc4dd79a0"
// ZG_API_KEY = "ed5d8197-26db-45be-b1ce-719f13847b6c"
// test site
//ZG_API_URL = "http://47.99.74.117:8883/api/v1/private/"