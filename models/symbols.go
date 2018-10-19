package models


var OKexSymbolsArray,HuobiSymbolsArray [][2]string

var OKexSymbols,HuobiSymbols map[string]int

var AllSymbols []string
// 分别获取获取所有的交易对
func GetSymbols(okexSymbols,huobiSymbols [][2]string) ([]string) {
	var tempSymbols = make([]string,0)
	for _,v := range OKexSymbolsArray {
		key := v[0] + "_" +v[1]
		tempSymbols = append(tempSymbols,key)
	}
	for _,v := range HuobiSymbolsArray {
		key := v[0] + "_" +v[1]
		tempSymbols = append(tempSymbols,key)
	}
	return tempSymbols
}

// 还缺少