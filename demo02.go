package main

import (
	"tradeRobot/utils"
	"fmt"
)

func main() {

	var account = utils.OKexRestfulApiRequest{}

	//account.API_KEY = "ed5d8197-26db-45be-b1ce-719f13847b6c"
	//account.SECRET_KEY = "d8fedb73-17d1-4a94-a6f6-f45bc4dd79a0"


	account.SECRET_KEY = "E624AB53C566F452D9AB074DB3B60BE0"
	account.API_KEY = "05b1f7b1-088d-4799-b702-d67d4d040105"

	account.OKexQureyMd5Sign()
	//sign := "B11018CC7816DEAE425A9A0967EA3141"
	//sign2 := "B11018CC7816DEAE425A9A0967EA3141"
	//account.OKexQureyMd5Sign()
	fmt.Println(account.Sign)
	account.OKexGetUserAssets()
}
