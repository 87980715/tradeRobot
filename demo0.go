package main

import (
	"fmt"
	"strconv"
)

func main() {

	amount,_:= strconv.ParseFloat("2985.18000000", 64)
	fmt.Println(amount)
}