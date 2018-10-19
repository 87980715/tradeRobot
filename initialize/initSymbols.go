package initialize

import "tradeRobot/models"

func initSymbol() {
	models.AllSymbols= models.GetSymbols(models.OKexSymbolsArray,models.HuobiSymbolsArray)
}