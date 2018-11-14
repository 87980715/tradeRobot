package models

// 防止暂停后启动时，id 都为零

var HuoPreDealId = make(map[string]int64)
var ZGPreDealId = make(map[string]int64)
var HuoPreTradeId = make(map[string]int64)
