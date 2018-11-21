package trade

import (
	"tradeRobot/robot/initialize"
	"context"
)

func TradeServerRun(symbol []string,ctx context.Context) {

	go func(ctx context.Context) {
		initialize.AppConfig.Wg.Add(1)
		TradeLimitZG(ctx)
	}(ctx)

	go func(symbol []string,ctx context.Context) {
		initialize.AppConfig.Wg.Add(1)
		CanleOrdersZG(symbol,ctx)
	}(symbol,ctx)


	go func(ctx context.Context) {
		initialize.AppConfig.Wg.Add(1)
		TradeLimitHuobi(ctx)
	}(ctx)

	go func(symbol []string,ctx context.Context) {
		initialize.AppConfig.Wg.Add(1)
		TradeCancelHuobi(symbol,ctx)
	}(symbol,ctx)

	go func(symbol []string,ctx context.Context) {
		HuobiInsertToDB(symbol,ctx)
	}(symbol,ctx)

}
