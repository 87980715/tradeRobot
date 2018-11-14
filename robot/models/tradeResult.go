package models

type HuobiTradeResults struct {
	Id          string  `gorm:"type:int;primary_key;AUTO_INCREMENT"`
	User_id     string  `grom:"type:int;index:user_id"`		 // 用户ID,索引，用于查询
	Trade_id    string  `grom:"type:int;"`
	Symbol      string  `gorm:"type:char(10);not null;"`
	Type        string  `gorm:"type:char(15);not null;"`      // 交易类型
	Price       string  `gorm:"type:decimal(15,8);not null;"` // 交易价格
	Deal_amount string  `gorm:"type:decimal(15,8);not null;"` // 交易数量
	Deal_fees   string  `gorm:"type:decimal(15,8);not null;"` // 交易费用
	Created_at  string  `gorm:"type:decimal(15,0);not null;"`
	Total       string  `gorm:"type:decimal(18,8);not null;"` // 成交总价
}

type ZGTradeResults struct {
	Id          string  `gorm:"type:int;primary_key;AUTO_INCREMENT"`
	User_id     string  `grom:"type:int;index:user_id"`
	Trade_id    string  `grom:"type:int;"`
	Symbol      string  `gorm:"type:char(10);not null;"`
	Type        string  `gorm:"type:char(15);not null;"`
	Price       string  `gorm:"type:decimal(15,8);not null;"`
	Deal_amount string  `gorm:"type:decimal(15,8);not null;"`
	Deal_fees   string  `gorm:"type:decimal(15,8);not null;"`
	Created_at  string  `gorm:"type:decimal(15,0);not null;"`
	Total       string  `gorm:"type:decimal(18,8);not null;"`
}
