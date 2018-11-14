package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"tradeRobot/robot/utils"
	"time"
)

var id = make(map[string]int)

func main() {
	db, err := utils.LoadExchangeDB()

	defer db.Close()
	for {
		time.Sleep(1*time.Second)
		if err != nil {
			logs.Error("initMaria failed err:", err)
			return
		}

		table := fmt.Sprintf("user_deal_history_%d", 26%100)

		queryStr := "select deal_id from " + table + " order by id desc limit 10 "

		fmt.Println(time.Now().Second())
		row, err := db.Query(queryStr)
		if err != nil {
			logs.Error("select deal_id form %s failed err %v:", table, err)
			return
		}
		defer row.Columns()
		var dealId = 0
		dealIds := make([]int, 0)
		for row.Next() {
			row.Scan(&dealId)
			dealIds = append(dealIds, dealId)
		}
		fmt.Println(dealIds)
	}
}
