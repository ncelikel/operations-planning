package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)


func main() {
	rand.Seed(time.Now().UnixNano())
	var prob Problem
	prob.readInit()
	PARENT_POOL_SIZE:=120
	MATING_POOL_SIZE:=240
	SIZE_COEFF:=float32(0.15)
	LOTSIZE_CONSERVATION:=float32(0.7)
	MUTATION_PROB:=0.5

	parent_pool:=[]Chromosome{}
	parent_objective_pool:=[]float32{}
	parent_ranking_pool:=[]int{}
	parent_pool_worst:=float32(0.0)
	parent_pool_best:=float32(math.Inf(1))
	for i:=0;i<PARENT_POOL_SIZE;i++{
		parent_pool=append(parent_pool, prob.randInit(0.2,0.85))
		parent_pool[i].yieldAll(prob)
		parent_objective_pool=append(parent_objective_pool, parent_pool[i].objective)
		if(parent_objective_pool[i]>parent_pool_worst){
			parent_pool_worst=parent_objective_pool[i]
		}
		if(parent_objective_pool[i]<parent_pool_best){
			parent_pool_best=parent_objective_pool[i]
		}
		parent_ranking_pool=append(parent_ranking_pool, i)
		insertSorted(parent_ranking_pool,parent_objective_pool)
	}

	mutation_rand:=0.0

	for iter:=0;iter<10000;iter++{
		mating_pool:=[]Chromosome{}
		mating_objective_pool:=[]float32{}
		for mating:=0;mating<MATING_POOL_SIZE/2;mating++{
			children:=parent_pool[parent_ranking_pool[mating]].blockCrossover(SIZE_COEFF,prob,parent_pool[rand.Intn(PARENT_POOL_SIZE)],LOTSIZE_CONSERVATION)
			
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
			if(mating_objective_pool[i]<parent_objective_pool[parent_ranking_pool[len(parent_ranking_pool)-1]]){
				parent_pool[parent_ranking_pool[len(parent_ranking_pool)-1]]=Chromosome{}
				parent_pool[parent_ranking_pool[len(parent_ranking_pool)-1]].copyChromosome(child)
				parent_objective_pool[parent_ranking_pool[len(parent_ranking_pool)-1]]=child.objective
				insertSorted(parent_ranking_pool,parent_objective_pool)
			}
		}
		fmt.Println("best:",parent_objective_pool[parent_ranking_pool[0]])
		fmt.Println("mean:",mean(parent_objective_pool))
		fmt.Println("25th:",parent_objective_pool[parent_ranking_pool[PARENT_POOL_SIZE/4]])
		fmt.Println("50th:",parent_objective_pool[parent_ranking_pool[PARENT_POOL_SIZE/2]])
		fmt.Println("75th:",parent_objective_pool[parent_ranking_pool[PARENT_POOL_SIZE*3/4]])
		fmt.Println("worst:",parent_objective_pool[parent_ranking_pool[PARENT_POOL_SIZE-1]])
		fmt.Println()
	}

	writeChromosome(parent_pool[parent_ranking_pool[0]],"best_chromosome")
}
