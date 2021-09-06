package main


import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	var prob Problem
	prob.readInit()
	chex := prob.randInit(0.2,[]float32{0.8,0.9})
	fmt.Println(chex.machineLayer)
	fmt.Println(chex.lotsizeLayer)
	fmt.Println(chex.mpInvInd)
	fmt.Println(chex.last)
}
