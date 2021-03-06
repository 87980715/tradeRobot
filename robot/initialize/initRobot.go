package initialize

import (
	"github.com/astaxie/beego/logs"
	"fmt"
	"tradeRobot/robot/utils"
	)

func InitRobot(fileName string) (err error){


	err = InitConf("ini", fileName)
	if err != nil {
		fmt.Printf("init conf failed, err:%v\n", err)
		panic("init conf failed")
		return
	}

	err = InitLogger()
	if err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		panic("init logger failed")
		return
	}

	err = utils.InitRobotDB()
	if err != nil {
		logs.Error("init mysqldb failed, err:%v\n", err)
		panic("init db failed")
		return
	}

	utils.ExchangeDB,err = utils.LoadExchangeDB()
	if err != nil {
		logs.Error("init mysqldb failed, err:%v\n", err)
		panic("init exchangedb failed")
		return
	}

	//logs.Info("init Accounts success")

	return
}
