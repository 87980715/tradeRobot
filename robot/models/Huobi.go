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
	Huobi_API_URL           = "api.huobi.pro"
	//Huobi_Post_Content_type = "application/json"
	//Huobi_Get_Content_type  = "application/x-www-form-urlencoded"

	Huobi_PendingOrdersSize = 15 // huobi 每次未成交订单获取条数
	Huobi_FilledOrdersSize  = 15 // huobi 每次已成交订单获取条数
	Huobi_OrdersSize = "10" // 查询火币交易历史
)
var (
	Huobi_AccessKeyId string = "eca1800e-e94af9ce-c2a77a7b-7a8b4"
	Huobi_Secretkey   string = "5658bce4-10643a8d-7e62938f-24139"
	Huobi_Account_ID string

	TradeAmountMultiple float64
	TradeInspectTime    int64
	TradePriceAdjust    float64
	HuobiUserID         string
)
