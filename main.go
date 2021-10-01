package main

import (
	"math/rand"
	"time"
)


func main() {
	rand.Seed(time.Now().UnixNano())
	var prob Problem
	prob.readInit()
	PARENT_POOL_SIZE:=50
	MATING_POOL_SIZE:=100
	SIZE_COEFF:=float32(0.15)
	LOTSIZE_CONSERVATION:=float32(0.7)
	MUTATION_PROB:=0.5
	N_ISLANDS:=5
	MIGRATION_COEFF:=0.05

	islands:=[]Island{}
	for island:=0;island<N_ISLANDS;island++{
		islands=append(islands, Island{[]Chromosome{},[]float32{},[]int{},0.0,true})
		for i:=0;i<PARENT_POOL_SIZE;i++{
			islands[island].parent_pool=append(islands[island].parent_pool, prob.randInit(0.2,0.85))
			islands[island].parent_pool[i].yieldAll(prob)
			islands[island].parent_objective_pool=append(islands[island].parent_objective_pool, islands[island].parent_pool[i].objective)
			islands[island].parent_ranking_pool=append(islands[island].parent_ranking_pool, i)
			insertSorted(islands[island].parent_ranking_pool,islands[island].parent_objective_pool)
		}
	}
	waiting_convergence:=false
	var convergence_obj float32
	converged:=true
	for iter:=0;iter<100;iter++{
		for island:=0;island<N_ISLANDS;island++{
			if(islands[island].iteration_continue){ 
				islands[island].iterate(100,prob,waiting_convergence,PARENT_POOL_SIZE,MATING_POOL_SIZE,LOTSIZE_CONSERVATION,MUTATION_PROB,SIZE_COEFF)
			} else {if(!waiting_convergence) {waiting_convergence=true; convergence_obj=islands[island].parent_objective_pool[islands[island].parent_ranking_pool[0]]}}
		}
		if(waiting_convergence){ //some island has approached a plateau, and waiting for others to come closer to its best objective. When they do, migration will take place
			for island:=0;island<N_ISLANDS;island++{
				if(!islands[island].iteration_continue || islands[island].parent_objective_pool[islands[island].parent_ranking_pool[0]]<float32(convergence_obj)){
					converged=true
					islands[island].iteration_continue=false
				} else {converged=false}
			}
		}	
		if(waiting_convergence && converged){ //migration happens here
			converged=true
			waiting_convergence=false
			migrators:=[]Chromosome{}
			for island:=0;island<N_ISLANDS;island++{
				islands[island].iteration_continue=true //set parameters ready to restart convergence run
				islands[island].stop_meter=0.0
				for m:=0;m<int(float64(PARENT_POOL_SIZE)*MIGRATION_COEFF);m++{
					migrators=append(migrators,islands[island].parent_pool[islands[island].parent_ranking_pool[m]])
				}
			}
			for island:=0;island<N_ISLANDS;island++{
				for m:=0;m<int(float64(PARENT_POOL_SIZE)*MIGRATION_COEFF);m++{
					islands[island].parent_pool[islands[island].parent_ranking_pool[len(islands[island].parent_ranking_pool)-1]]=Chromosome{}
					islands[island].parent_pool[islands[island].parent_ranking_pool[len(islands[island].parent_ranking_pool)-1]].copyChromosome(migrators[m])
					islands[island].parent_objective_pool[islands[island].parent_ranking_pool[len(islands[island].parent_ranking_pool)-1]]=migrators[m].objective
					insertSorted(islands[island].parent_ranking_pool,islands[island].parent_objective_pool)
				}
			}
		}
	}
	pangaea:=Island{[]Chromosome{},[]float32{},[]int{},0.0,true}
	for island:=0;island<N_ISLANDS;island++{
		for chr:=0;chr<10;chr++{
			pangaea.parent_pool=append(pangaea.parent_pool, islands[island].parent_pool[islands[island].parent_ranking_pool[chr]])
			pangaea.parent_pool[len(pangaea.parent_pool)-1].yieldAll(prob)
			pangaea.parent_objective_pool=append(pangaea.parent_objective_pool, pangaea.parent_pool[len(pangaea.parent_pool)-1].objective)
			pangaea.parent_ranking_pool=append(pangaea.parent_ranking_pool, len(pangaea.parent_pool)-1)
			insertSorted(pangaea.parent_ranking_pool,pangaea.parent_objective_pool)
		}
	}
	pangaea.iterate(1000,prob,false,PARENT_POOL_SIZE,MATING_POOL_SIZE,LOTSIZE_CONSERVATION,MUTATION_PROB,SIZE_COEFF)


	writeChromosome(pangaea.parent_pool[pangaea.parent_ranking_pool[0]],"best_chromosome")
}