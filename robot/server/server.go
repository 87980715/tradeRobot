package server

import (
	"tradeRobot/robot/initialize"
	"tradeRobot/robot/dataAgent"
	"tradeRobot/robot/trade"
	"context"
)

func RobotRun(symbol []string,ctx context.Context) {

	initialize.AppConfig.Wg.Add(1)

	dataAgent.AgentServerRun(symbol,ctx)

	trade.TradeServerRun(symbol,ctx)

	initialize.AppConfig.Wg.Wait()
}