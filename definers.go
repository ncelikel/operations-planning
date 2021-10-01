package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
)

type Problem struct {
	nMachine  int
	nProduct  int
	nPeriod   int
	dPeriod   float32
	socket    []float32 //how many sockets each machine has
	cycleTime []float32 //cycle time of products
	chgOver   []float32 //changeover duration per product
	mpMatch   [][]int   //list of machines that each product can be produced in
	curDemand [][]int   //demand of each item in each period
	cumDemand [][]int   //total demand of each item up to each period
}

type Chromosome struct {
	machineLayer  [][]int //machineLayer[product][period]=machine -> indicates that the product will be produced on period on machine. If -1; product is not produced on that period on any machine
	lotsizeLayer  [][]int //lotsizeLayer[product][period]=lotsize -> how many units to be produced
	curProduction [][]int //total production on that [product][period]
	cumProduction [][]int //cumulative production of product up to that period
	last          [][]int //last[machine][period] = product; which product's mold will be kept on machine at the end of period
	mInvInd       [][][]int //inverse index for machines; shows at which [product][period] couples the machine is used in. Always keep sorted by period
	mpInvInd      [][][]int //inverse index for [machine][period]; which gives the list of products produced on that machine in that period
	availability  [][]float32 //availability[machine][period]; total time remaining for production after subtracting required changeovers.
	utilization   [][]float32 //utilization[machine][period]; what percentage of the total time will be efficient.
	objective     float32 //sum of all costs; deficit, inventory and changeover
}

type Island struct{
	parent_pool []Chromosome
	parent_objective_pool []float32
	parent_ranking_pool []int
	stop_meter float32
	iteration_continue bool
}

func (prob *Problem) readInit() { //read problem description from local .json file
	data, err := ioutil.ReadFile("initializer.json")
	if err != nil {
		fmt.Println(err)
	}
	var result map[string]interface{}
	json.Unmarshal(data, &result)
	prob.nMachine = int(result["nMachine"].(float64))
	prob.nProduct = int(result["nProduct"].(float64))
	prob.nPeriod = int(result["nPeriod"].(float64))
	prob.dPeriod = float32(result["dPeriod"].(float64))
	socket := reflect.ValueOf(result["socket"])
	prob.socket = make([]float32, socket.Len())
	for i := 0; i < socket.Len(); i++ {
		prob.socket[i] = float32(socket.Index(i).Interface().(float64))
	}
	cycleTime := reflect.ValueOf(result["cycleTime"])
	prob.cycleTime = make([]float32, cycleTime.Len())
	for i := 0; i < cycleTime.Len(); i++ {
		prob.cycleTime[i] = float32(cycleTime.Index(i).Interface().(float64))
	}
	chgOver := reflect.ValueOf(result["chgOver"])
	prob.chgOver = make([]float32, chgOver.Len())
	for i := 0; i < chgOver.Len(); i++ {
		prob.chgOver[i] = float32(chgOver.Index(i).Interface().(float64))
	}
	mpMatch := reflect.ValueOf(result["mpMatch"])
	prob.mpMatch = make([][]int, mpMatch.Len())
	for i := 0; i < mpMatch.Len(); i++ {
		prob.mpMatch[i] = make([]int, reflect.ValueOf(mpMatch.Index(i).Interface()).Len())
		for j := 0; j < len(prob.mpMatch[i]); j++ {
			prob.mpMatch[i][j] = int(reflect.ValueOf(mpMatch.Index(i).Interface()).Index(j).Interface().(float64))
		}
	}
	curDemand := reflect.ValueOf(result["curDemand"])
	prob.curDemand = make([][]int, curDemand.Len())
	for i := 0; i < curDemand.Len(); i++ {
		prob.curDemand[i] = make([]int, reflect.ValueOf(curDemand.Index(i).Interface()).Len())
		for j := 0; j < len(prob.curDemand[i]); j++ {
			prob.curDemand[i][j] = int(reflect.ValueOf(curDemand.Index(i).Interface()).Index(j).Interface().(float64))
		}
	}
	cumDemand := reflect.ValueOf(result["cumDemand"])
	prob.cumDemand = make([][]int, cumDemand.Len())
	for i := 0; i < cumDemand.Len(); i++ {
		prob.cumDemand[i] = make([]int, reflect.ValueOf(cumDemand.Index(i).Interface()).Len())
		for j := 0; j < len(prob.cumDemand[i]); j++ {
			prob.cumDemand[i][j] = int(reflect.ValueOf(cumDemand.Index(i).Interface()).Index(j).Interface().(float64))
		}
	}
}
