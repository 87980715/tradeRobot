package main

import (
	_ "tradeRobot/admin/routers"
	"github.com/astaxie/beego"
	"tradeRobot/robot/initialize"
	"github.com/astaxie/beego/logs"
	"runtime"
	"flag"
)



func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	fileName := flag.String("fileName", "/root/go/src/tradeRobot/robot/conf/app.conf", "configFilePath")
	flag.Parse()

	err := initialize.InitRobot(*fileName)
	if err != nil {
		logs.Error("init robot failed err:",err)
		panic("init robot failed")
	}

	logs.Debug("init conf succ")


	beego.Run()
}

