package dataAgent

import (
	"tradeRobot/robot/initialize"
	"tradeRobot/robot/models"
	"context"
)

func AgentServerRun(symbol []string,ctx context.Context) {

	go func(symbol []string,ctx context.Context) {
		initialize.AppConfig.Wg.Add(1)
		GetTradesHuobi(symbol, models.Huobi_OrdersSize,ctx)
	}(symbol,ctx)

	go func(symbol []string, ctx context.Context) {
		initialize.AppConfig.Wg.Add(1)
		GetDealOrdersZG(symbol,ctx)
	}(symbol,ctx)

	go func(symbol []string,ctx context.Context) {
		QueryDealIds(symbol,ctx)
	}(symbol,ctx)

	go func(ctx context.Context){
		initialize.AppConfig.Wg.Add(1)
		QueryRealDealZG(ctx)
	}(ctx)

	go func(symbol []string,ctx context.Context) {
		initialize.AppConfig.Wg.Add(1)
		ZGInsertToDB(symbol ,ctx)
	}(symbol ,ctx)

	// 获取usdt价格
	go func(ctx context.Context){
		initialize.AppConfig.Wg.Add(1)
		GetHuobiUsdtPrice(ctx)
	}(ctx)
}
