package main

import (
	"fmt"
)

func main() {

	s := fmt.Sprintf("user_deal_history_%d",316%100)
	fmt.Println(s)

}
