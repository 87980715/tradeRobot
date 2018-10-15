package utils

import (

)

type FromDataRpcReq struct {
	Id     int64         `json:"id"`
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

/*func NewJsonRpcReq(method string, params []interface{}) JsonRpcReq {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return JsonRpcReq{
		Id:     time.Now().UnixNano() + rnd.Int63n(1000),
		Method: method,
		Params: params,
	}
}
*/

type JsonRpcRespError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	error
}

func (e *JsonRpcRespError) Error() string {
	return e.Message
}

type JsonRpcResp struct {
	Error  *JsonRpcRespError `json:"error"`
	Id     int64             `json:"id"`
	Result interface{}       `json:"result"`
}
/*
func TradeRpcCall(url string, method string, params []interface{}) (*JsonRpcResp, error) {
	v := url.Values{}
	v.Add("api_key", "ed5d8197-26db-45be-b1ce-719f13847b6c")
	v.Add("symbol","ZG_CNZ")
	v.Add("side","1")
	v.Add("","")
	v.Add("price","50")
	v.Add("amount", "0.02")

	signStr := "amount=1.0&api_key=c821db84-6fbd-11e4-a9e3-c86000d26d7c&price=680&side=1&symbol=ZG_CNZ&secret_key=secretKey"
	hash := md5.Sum([]byte(signStr))
	hashed := hash[:]
	sign := hex.EncodeToString(hashed)
	v.Add("sign",sign)

	rd := ioutil.NopCloser(strings.NewReader(v.Encode()))
	resp, err := http.Post(url, "application/x-www-form-urlencoded", rd)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		Doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			logs.Error("goquery.NewDocumentFromReader failed err:", err)
		}
	}
}
*/