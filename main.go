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
	chromosome_1 := prob.randInit(0.2, 0.85)
	_ = chromosome_1.yieldAll(prob)
	chromosome_2 := prob.randInit(0.2, 0.85)
	_ = chromosome_2.yieldAll(prob)

	children:=chromosome_1.blockCrossover(0.2,prob,chromosome_2)
	_=children[0].yieldAll(prob)
	_=children[1].yieldAll(prob)
	fmt.Println(chromosome_1.machineLayer,chromosome_1.lotsizeLayer,chromosome_1.last,chromosome_1.objective)
	fmt.Println(chromosome_2.machineLayer,chromosome_2.lotsizeLayer,chromosome_2.last,chromosome_2.objective)

	fmt.Println(children[0].machineLayer,children[0].lotsizeLayer,children[0].last,children[0].objective)
	fmt.Println(children[1].machineLayer,children[1].lotsizeLayer,children[1].last,children[1].objective)
	}
