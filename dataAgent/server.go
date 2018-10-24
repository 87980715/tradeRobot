package dataAgent

import (
	"tradeRobot/initialize"
	"tradeRobot/models"
)

func AgentServerRun() {


	if len(models.OKexSymbolsArray) != 0 {
		go func() {
			initialize.AppConfig.Wg.Add(1)
			GetTradesOKex(models.OKexSymbolsArray,models.LIMITED)
		}()
	}

	if len(models.HuobiSymbolsArray) != 0 {
		go func() {
			initialize.AppConfig.Wg.Add(1)
			GetTradesHuobi(models.OKexSymbolsArray,models.LIMITED)
		}()
	}

	/*
	if len(models.HuobiSymbolsArray) != 0 {
		go func() {
			initialize.AppConfig.Wg.Add(1)
			GetCionDepthHuobi(models.HuobiSymbolsArray)
		}()
	}
	*/
	/*
	go func() {
		initialize.AppConfig.Wg.Add(1)
		GetFinishedOrdersFromZT()
	}()
	*/
}
