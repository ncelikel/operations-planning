package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
)

func (ch *Chromosome) copyChromosome(c2 Chromosome) {
	/*	machineLayer  [][]int //machineLayer[product][period]=machine -> indicates that the product will be produced on period on machine. If -1; product is not produced on that period on any machine
		lotsizeLayer  [][]int //lotsizeLayer[product][period]=lotsize -> how many units to be produced
		curProduction [][]int //total production on that [product][period]
		cumProduction [][]int //cumulative production of product up to that period
		last          [][]int //last[machine][period] = product; which product's mold will be kept on machine at the end of period
		mInvInd       [][][]int //inverse index for machines; shows at which [product][period] couples the machine is used in. Always keep sorted by period
		mpInvInd      [][][]int //inverse index for [machine][period]; which gives the list of products produced on that machine in that period
		availability  [][]float32 //availability[machine][period]; total time remaining for production after subtracting required changeovers.
		utilization   [][]float32 //utilization[machine][period]; what percentage of the total time will be efficient.
		objective     float32 //sum of all costs; deficit, inventory and changeover */
	nProduct := len(c2.machineLayer)
	nPeriod := len(c2.machineLayer[0])
	nMachine := len(c2.mpInvInd)
	ch.machineLayer = make([][]int, nProduct)
	ch.lotsizeLayer = make([][]int, nProduct)
	ch.curProduction = make([][]int, nProduct)
	ch.cumProduction = make([][]int, nProduct)
	for i := 0; i < nProduct; i++ {
		ch.machineLayer[i] = make([]int, nPeriod)
		ch.lotsizeLayer[i] = make([]int, nPeriod)
		ch.curProduction[i] = make([]int, nPeriod)
		ch.cumProduction[i] = make([]int, nPeriod)

		copy(ch.machineLayer[i], c2.machineLayer[i])
		copy(ch.lotsizeLayer[i], c2.lotsizeLayer[i])
		copy(ch.curProduction[i], c2.curProduction[i])
		copy(ch.cumProduction[i], c2.cumProduction[i])
	}
	ch.last = make([][]int, nMachine)
	ch.mpInvInd = make([][][]int, nMachine)
	ch.availability = make([][]float32, nMachine)
	ch.utilization = make([][]float32, nMachine)
	for i := 0; i < nMachine; i++ {
		ch.last[i] = make([]int, nPeriod)
		ch.mpInvInd[i] = make([][]int, nPeriod)
		ch.availability[i] = make([]float32, nPeriod)
		ch.utilization[i] = make([]float32, nPeriod)

		copy(ch.last[i], c2.last[i])
		copy(ch.availability[i], c2.availability[i])
		copy(ch.utilization[i], c2.utilization[i])
		for j := 0; j < nPeriod; j++ {
			ch.mpInvInd[i][j] = make([]int, len(c2.mpInvInd[i][j]))
			copy(ch.mpInvInd[i][j], c2.mpInvInd[i][j])
		}
	}
	ch.objective = c2.objective
}

type chJSON struct {
	MachineLayer  [][]int     //machineLayer[product][period]=machine -> indicates that the product will be produced on period on machine. If -1; product is not produced on that period on any machine
	LotsizeLayer  [][]int     //lotsizeLayer[product][period]=lotsize -> how many units to be produced
	CurProduction [][]int     //total production on that [product][period]
	CumProduction [][]int     //cumulative production of product up to that period
	Last          [][]int     //last[machine][period] = product; which product's mold will be kept on machine at the end of period
	MInvInd       [][][]int   //inverse index for machines; shows at which [product][period] couples the machine is used in. Always keep sorted by period
	MpInvInd      [][][]int   //inverse index for [machine][period]; which gives the list of products produced on that machine in that period
	Availability  [][]float32 //availability[machine][period]; total time remaining for production after subtracting required changeovers.
	Utilization   [][]float32 //utilization[machine][period]; what percentage of the total time will be efficient.
	Objective     float32     //sum of all costs; deficit, inventory and changeover
}

func writeChromosome(v Chromosome, fn string) {
	jsondat := &chJSON{MachineLayer: v.machineLayer, LotsizeLayer: v.lotsizeLayer, CurProduction: v.curProduction, CumProduction: v.cumProduction,
		Last: v.last, MpInvInd: v.mpInvInd, Availability: v.availability,
		Utilization: v.utilization, Objective: v.objective}
	encjson, _ := json.MarshalIndent(jsondat, "", " ")

	_ = ioutil.WriteFile(fn+".json", encjson, 0644)
}

func insertSorted(sortpool []int, valuepool []float32) {
	for i := 0; i < len(sortpool)-1; i++ {
		if valuepool[sortpool[len(sortpool)-2-i]] > valuepool[sortpool[len(sortpool)-1-i]] {
			swap(sortpool, len(valuepool)-2-i, len(valuepool)-1-i)
		} else {
			break
		}
	}
}

func swap(arr []int, a int, b int) {
	tmp := arr[b]
	arr[b] = arr[a]
	arr[a] = tmp
}

func mean(arr []float32) float32{
	total:=float32(0.0)
	for _,v:=range arr {
		total+=v
	}
	return(total/float32(len(arr)))
}




func (isl *Island) iterate(n int,prob Problem,waiting_convergence bool,PARENT_POOL_SIZE int,MATING_POOL_SIZE int,LOTSIZE_CONSERVATION float32,MUTATION_PROB float64,SIZE_COEFF float32) {
		var prevBest float32
		var newBest float32
		mutation_rand:=0.0
		
		for iter:=0;iter<n;iter++{
			prevBest=isl.parent_objective_pool[isl.parent_ranking_pool[0]]
			mating_pool:=[]Chromosome{}
			mating_objective_pool:=[]float32{}
			for mating:=0;mating<MATING_POOL_SIZE/2;mating++{
				children:=isl.parent_pool[isl.parent_ranking_pool[mating]].blockCrossover(SIZE_COEFF,prob,isl.parent_pool[rand.Intn(PARENT_POOL_SIZE)],LOTSIZE_CONSERVATION)
				
				mating_pool=append(mating_pool, Chromosome{})
				mating_pool[len(mating_pool)-1].copyChromosome(children[0])
				mating_pool[len(mating_pool)-1].yieldAll(prob)
				mating_objective_pool=append(mating_objective_pool, mating_pool[len(mating_pool)-1].objective)
	
				mating_pool=append(mating_pool, Chromosome{})
				mating_pool[len(mating_pool)-1].copyChromosome(children[1])
				mating_pool[len(mating_pool)-1].yieldAll(prob)
				mating_objective_pool=append(mating_objective_pool, mating_pool[len(mating_pool)-1].objective)		
			}
			for i:=0;i<len(mating_pool);i++ {
				if(rand.Float64()<MUTATION_PROB){
					mutation_rand=rand.Float64()
					if(mutation_rand<0.2){
						mating_pool=append(mating_pool, mating_pool[i].utilMutation(0.08,prob))
					} else if (mutation_rand<0.7){
						mating_pool=append(mating_pool, mating_pool[i].lotsizeMutation(0.10,prob))
					} else {
						mating_pool=append(mating_pool, mating_pool[i].machineMutation(0.10,prob,0.25,6))
					}
					mating_pool[len(mating_pool)-1].yieldAll(prob)
					mating_objective_pool = append(mating_objective_pool, mating_pool[len(mating_pool)-1].objective)
				}
			}
			for i,child:=range mating_pool{
				if(mating_objective_pool[i]<isl.parent_objective_pool[isl.parent_ranking_pool[len(isl.parent_ranking_pool)-1]]){
					isl.parent_pool[isl.parent_ranking_pool[len(isl.parent_ranking_pool)-1]]=Chromosome{}
					isl.parent_pool[isl.parent_ranking_pool[len(isl.parent_ranking_pool)-1]].copyChromosome(child)
					isl.parent_objective_pool[isl.parent_ranking_pool[len(isl.parent_ranking_pool)-1]]=child.objective
					insertSorted(isl.parent_ranking_pool,isl.parent_objective_pool)
				}			
			}
			newBest=isl.parent_objective_pool[isl.parent_ranking_pool[0]]
			if(prevBest*0.98<=newBest){
				isl.stop_meter+=0.05
			} else {isl.stop_meter-=0.03}
			if(isl.stop_meter>1.0 && !waiting_convergence){
				isl.stop_meter=0.0
				isl.iteration_continue=false
			}
			fmt.Println("best:",isl.parent_objective_pool[isl.parent_ranking_pool[0]])
			fmt.Println("mean:",mean(isl.parent_objective_pool))
			fmt.Println("worst:",isl.parent_objective_pool[isl.parent_ranking_pool[PARENT_POOL_SIZE-1]])
			fmt.Println()
		}
}