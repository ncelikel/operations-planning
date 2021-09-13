package main

import (
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

	children:=chromosome_1.blockCrossover(0.3,prob,chromosome_2)
	_=children[0].yieldAll(prob)
	_=children[1].yieldAll(prob)

	writeChromosome(chromosome_1,"parent_1")
	writeChromosome(chromosome_2,"parent_2")
	writeChromosome(children[0],"child_1")
	writeChromosome(children[1],"child_2")

	}
