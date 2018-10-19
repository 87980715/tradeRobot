package trade

import (
	"tradeRobot/initialize"
)

func ZTTradeServerRun() {

	/*
	go func() {
		initialize.AppConfig.Wg.Add(1)
		ReadTickersFromRedis(coinNames,toCoinNames)
	}()
	*/
	go func() {
		initialize.AppConfig.Wg.Add(1)
		ReadDepthsFromRedis()
	}()
	go func() {
		initialize.AppConfig.Wg.Add(1)
		TradeLimitZT()
	}()
	/*
	go func() {
		initialize.AppConfig.Wg.Add(1)
		TradeMarketZT()
	}()
	*/
}
/*
func OKexTradeServerRun(coinNames map[int]string, toCoinNames map[int]string, accounts []*Account) {

	go func() {
		initialize.AppConfig.Wg.Add(1)
		ReadTickersFromRedis(coinNames,toCoinNames)
	}()

	go func() {
		initialize.AppConfig.Wg.Add(1)
		ReadDepthsFromRedis(coinNames,toCoinNames)
	}()

	go func() {
		initialize.AppConfig.Wg.Add(1)
		TradeLimitZT(coinNames,toCoinNames,accounts)
	}()

	go func() {
		initialize.AppConfig.Wg.Add(1)
		TradeMarketZT(coinNames,toCoinNames,accounts)
	}()

}
*/