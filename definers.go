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
	cycleTime []float32
	chgOver   []float32 //changeover duration per product
	mpMatch   [][]int   //list of machines that each product can be produced in
	curDemand [][]int   //demand of each item in each period
	cumDemand [][]int   //total demand of each item up to each period
}
type Chromosome struct {
	machineLayer  [][]int
	lotsizeLayer  [][]int
	curProduction [][]int
	cumProduction [][]int
	last          [][]int
	mInvInd       [][][]int //inverse index for machines; shows at which [product][period] couples the machine is used in. Always keep sorted by period
	mpInvInd      [][][]int
	availability  [][]float32
	utilization   [][]float32
	objective     float32
}

func (prob *Problem) readInit() {
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
