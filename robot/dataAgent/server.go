package dataAgent

import (
	"tradeRobot/robot/initialize"
	"tradeRobot/robot/models"
	"context"
)

func AgentServerRun(symbol []string,ctx context.Context) {

	go func(symbol []string,ctx context.Context) {
		initialize.AppConfig.Wg.Add(1)
		GetTradesHuobiMtEth(symbol, models.Huobi_OrdersSize,ctx)
	}(symbol,ctx)

	go func(ctx context.Context) {
		initialize.AppConfig.Wg.Add(1)
		GetDealOrdersZG(ctx)
	}(ctx)

	go func(symbol []string,ctx context.Context) {
		QueryDealIds(symbol,ctx)
	}(symbol,ctx)

	go func(ctx context.Context){
		initialize.AppConfig.Wg.Add(1)
		QueryRealDealZG(ctx)
	}(ctx)

	go func(ctx context.Context) {
		initialize.AppConfig.Wg.Add(1)
		ZGInsertToDB(ctx)
	}(ctx)

	// 获取usdt价格
	go func(ctx context.Context){
		initialize.AppConfig.Wg.Add(1)
		GetHuobiUsdtPrice(ctx)
	}(ctx)
}
