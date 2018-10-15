package trade

import (
	"tradeRobot/initialize"
)

func TradeServerRun(coinNames map[int]string, toCoinNames map[int]string, accounts []*Account) {

	/*
	go func() {
		initialize.AppConfig.Wg.Add(1)
		ReadTickersFromRedis(coinNames,toCoinNames)
	}()
	*/
	go func() {
		initialize.AppConfig.Wg.Add(1)
		ReadDepthsFromRedis(coinNames,toCoinNames)
	}()

	go func() {
		initialize.AppConfig.Wg.Add(1)
		TradeLimitZT(coinNames,toCoinNames,accounts)
	}()

	/*go func() {
		initialize.AppConfig.Wg.Add(1)
		TradeMarketZT(coinNames,toCoinNames,accounts)
	}()
	*/
}
