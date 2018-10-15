package dataAgent

import (
	"tradeRobot/initialize"
	"time"
)

func AgentServerRun(coinNames map[int]string, toCionNames map[int]string) {

	go func(){
		initialize.AppConfig.Wg.Add(1)
		for {
			time.Sleep(500*time.Millisecond)
			GetDataFromOkex(coinNames,toCionNames)
		}
	}()
}
