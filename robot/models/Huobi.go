package models

type HuobiAccountsData struct {
	ID     int64  `json:"id"`      // Account ID
	Type   string `json:"type"`    // 账户类型, spot: 现货账户
	State  string `json:"state"`   // 账户状态, working: 正常, lock: 账户被锁定
	UserID int64  `json:"user-id"` // 用户ID
}

type HuobiAccountsReturn struct {
	Status  string              `json:"status"` // 请求状态
	Data    []HuobiAccountsData `json:"data"`   // 用户数据
	ErrCode string              `json:"err-code"`
	ErrMsg  string              `json:"err-msg"`
}

var (
	Huobi_API_URL = "api.huobi.pro"
	// Huobi_Post_Content_type = "application/json"
	// Huobi_Get_Content_type  = "application/x-www-form-urlencoded"

	Huobi_PendingOrdersSize = 15   // huobi 每次未成交订单获取条数
	Huobi_FilledOrdersSize  = 15   // huobi 每次已成交订单获取条数
	Huobi_OrdersSize        = "10" // 查询火币交易历史

	Huobi_AccessKeyId string  = "b2af8f9f-4ac75b4d-4fce0763-1c789"
	Huobi_Secretkey   string  = "80653212-f2b1fc55-f1af7577-b9a3f"
	// Huobi_Account_ID  string = "4821321"

	TradeAmountMultiple float64
	TradeInspectTime    int64
	TradePriceAdjust    float64
	//HuobiUserID         string

	// 默认值设置
	UsdtPrice = map[string]float64{"huobi":6.84}
	EthPrice = make(map[string]float64)
)
