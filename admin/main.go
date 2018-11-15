package main

import (
	_ "tradeRobot/admin/routers"
	"github.com/astaxie/beego"
	"tradeRobot/robot/initialize"
	"github.com/astaxie/beego/logs"
	"runtime"
)



func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	err := initialize.InitRobot()
	if err != nil {
		logs.Error("init robot failed err:",err)
		panic("init robot failed")
	}

	logs.Debug("init conf succ")


	beego.Run()
}

