package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/logs"
	"tradeRobot/robot/models"
	"github.com/jinzhu/gorm"
	"strconv"
	"sync"
)

var (
	DB   *gorm.DB
	once sync.Once

	MariaDB *sql.DB
)

func main() {
	dealIds := queryDealIds(1805)

	qureyDealOerder(dealIds)

}

func initMysql() *sql.DB {
	db, err := sql.Open("mysql", "exchange:exchange@tcp(47.99.74.117:3306)/trade_history?charset=utf8")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return db
}

func queryDealIds(orderId int) []int {
	db := initMysql()
	defer db.Close()
	row, err := db.Query("select deal_id from user_deal_history_26 where order_id = ?", orderId)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Columns()
	var dealId = 0
	dealIds := make([]int, 0)
	for row.Next() {
		row.Scan(&dealId)
		dealIds = append(dealIds, dealId)
	}
	return dealIds
}

func qureyDealOerder(dealIds []int) {
	db := initMysql()
	defer db.Close()
	n := 0
	var CreatedAt, UserId, TradeId, Symbol, Type, Price, DealAmount, DealFee, Total string
	for _, id := range dealIds {
		row, err := db.Query("select COUNT(deal_id) from user_deal_history_26 where deal_id = ?", id)
		if err != nil {
			fmt.Println(err)
			return
		}
		for row.Next() {
			row.Scan(&n)
			// 真实交易过滤
			if n != 2 {
				row, err := db.Query("select time,user_id,market,order_id,side,price,amount,deal,deal_fee  from user_deal_history_26 where deal_id = ?", id)
				if err != nil {
					fmt.Println(err)
					return
				}
				for row.Next() {
					row.Scan(&CreatedAt, &UserId,  &Symbol, &TradeId,&Type, &Price, &DealAmount,&Total,&DealFee)
					tradeResult := &models.ZGTradeResults{}
					//tradeResult.Id = "1"
					tradeResult.User_id = UserId
					tradeResult.Trade_id = TradeId
					tradeResult.Symbol = Symbol
					tradeResult.Type = Type
					tradeResult.Price = Price
					tradeResult.Deal_amount = DealAmount
					tradeResult.Deal_fees = DealFee
					tradeResult.Created_at = CreatedAt
					tradeResult.Total = Total

					ZGInsertToDB(tradeResult)
				}
				row.Columns()
			}
			row.Columns()
		}
	}
}

func ZGInsertToDB(t *models.ZGTradeResults) {

	db,_:= loadDB()
	fmt.Println("连接成功")
	if err := db.Create(t).Error; err != nil {
		logs.Error("insert failed into Huobi tradeResult ")
		return
	}
	fmt.Println("插入成功")
	db.Close()
}


type DBConfig struct {
	User     string `default:"root"`
	Password string `default:""`
	Name     string `required:"true"`
	Port     uint   `default:"3306"`
	DbName   string `required:"true"`
	Charset  string `default:"utf8"`
}

func loadDB()(*gorm.DB,error) {
	dbConfig := new(DBConfig)
	dbConfig.Name = "127.0.0.1"
	dbConfig.User = "root"
	dbConfig.Password = "root"
	dbConfig.Port = 3306
	dbConfig.DbName = "tradeRobot"
	db, err := GetDBConnection(dbConfig)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetDBConnection(conf *DBConfig) (db *gorm.DB, err error) {
	once.Do(func() {
		str := conf.User + ":" +
			conf.Password + "@tcp(" +
			conf.Name + ":" +
			strconv.FormatUint(uint64(conf.Port), 10) + ")/" +
			conf.DbName + "?" +
			conf.Charset + "&parseTime=True&loc=Local"
		db, err = gorm.Open("mysql", str)
		if err != nil {
			logs.Error("gorm open mysql failed err:",err)
		}

		db.DB().SetMaxIdleConns(200)
		db.DB().SetMaxOpenConns(100)
	})
	return
}