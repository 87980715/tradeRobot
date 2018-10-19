package main


import (
	"time"
	"math/rand"
	"fmt"
)



func main(){
	rand.Seed(time.Now().UnixNano())
	randArray := make([]int, 0)

	for i := 0; i < 1; i++ {
		num := rand.Intn(8)
		randArray = append(randArray, num)
	}
	fmt.Println(randArray)
}


