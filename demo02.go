package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {

	rand.Seed(time.Now().UnixNano())
	a := rand.Float64()
	b := rand.Float64()
	c := rand.Intn(10)
	d := rand.Intn(10)

	fmt.Println(a,b,c,d)
}
