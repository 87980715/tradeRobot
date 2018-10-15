package trade


type OKexCurTicker struct {
	Date   string      `json:"date"`
	Ticker *OKexTicker `json:"ticker"`
}

type OKexTicker struct {
	Buy  string `json:"buy"`
	High string `json:"high"`
	Last string `json:"last"`
	Low  string `json:"low"`
	Sell string `json:"sell"`
	Vol  string `json:"vol"`
}

type OKexDepth struct {
	Asks [][2]float64 `json:"asks"`
	Bids [][2]float64 `json:"bids"`
}

// 储存币的数据
type OKexC2CData struct {
	Tickers []map[string]*OKexCurTicker `json:"tickers"`
	Depths  []map[string]*OKexDepth     `json:"depths"`
}

var TickersChan = make(chan map[string]string,100)

var DepthsChan = make(chan map[string]string,100)





