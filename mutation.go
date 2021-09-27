package main

import (
	"math/rand"
	"sort"
)

func (ch *Chromosome) utilMutation(sizeCoeff float32, p Problem) Chromosome {
	var mutated Chromosome
	mutated.copyChromosome(*ch)
	var availableDuration float32
	durationRatio := float32(0.0)
	var amt int

	var blockIndices [][]int
	blockIndices = append(blockIndices, []int{rand.Intn(p.nMachine), rand.Intn(p.nPeriod)})
	li := blockIndices[0]
	spreadProbability := 1.0 - (1.0 / (float32(p.nMachine) * float32(p.nPeriod) * sizeCoeff))
	dxy := []int{-1, 0, 1}
	for rand.Float32() < spreadProbability { //randomly spread to locations around
		li = []int{li[0] + dxy[rand.Intn(3)], li[0] + dxy[rand.Intn(3)]}
		if li[0] >= 0 && li[1] >= 0 && li[0] < p.nMachine && li[1] < p.nPeriod {
			blockIndices = append(blockIndices, li)
		}
	}
	allKeys := make([][]bool, p.nProduct)
	for i := 0; i < p.nProduct; i++ {
		allKeys[i] = make([]bool, p.nPeriod)
	}
	for _, location := range blockIndices {
		allKeys[location[0]][location[1]] = true
	}
	for imac, x := range allKeys {
		for iper, y := range x {
			if y {
				mutated.utilization[imac][iper] = rand.Float32()*0.2 + 0.75
				availableDuration = p.dPeriod * mutated.utilization[imac][iper]
				for _, prod := range mutated.mpInvInd[imac][iper] {
					if iper != 0 {
						if mutated.last[imac][iper-1] != prod {
							availableDuration -= p.chgOver[prod]
						}
					} else {
						availableDuration -= p.chgOver[prod]
					}
				}
				if availableDuration < 0.0 {
					for _, prod := range mutated.mpInvInd[imac][iper] {
						mutated.machineLayer[prod][iper] = -1
						mutated.lotsizeLayer[prod][iper] = 0
					}
					mutated.mpInvInd[imac][iper] = make([]int, 0)
					mutated.last[imac][iper] = -1
					break
				} else {
					for _, prod := range mutated.mpInvInd[imac][iper] {
						durationRatio += (float32(mutated.lotsizeLayer[prod][iper]) * float32(p.cycleTime[prod]) / float32(p.socket[prod]))
					}
					durationRatio = availableDuration / durationRatio
					for _, prod := range mutated.mpInvInd[imac][iper] {
						amt = mutated.lotsizeLayer[prod][iper]
						amt = int(float32(amt) * durationRatio)
						amt = amt - amt%int(p.socket[prod])
						mutated.lotsizeLayer[prod][iper] = amt
					}
				}
				mutated.availability[imac][iper] = availableDuration
			}
		}
	}
	return mutated
}

func (ch *Chromosome) lotsizeMutation(sizeCoeff float32, p Problem) Chromosome {
	var mutated Chromosome
	mutated.copyChromosome(*ch)
	var availableDuration float32

	var blockIndices [][]int
	blockIndices = append(blockIndices, []int{rand.Intn(p.nMachine), rand.Intn(p.nPeriod)})
	li := blockIndices[0]
	spreadProbability := 1.0 - (1.0 / (float32(p.nMachine) * float32(p.nPeriod) * sizeCoeff))
	dxy := []int{-1, 0, 1}
	for rand.Float32() < spreadProbability { //randomly spread to locations around
		li = []int{li[0] + dxy[rand.Intn(3)], li[0] + dxy[rand.Intn(3)]}
		if li[0] >= 0 && li[1] >= 0 && li[0] < p.nMachine && li[1] < p.nPeriod {
			blockIndices = append(blockIndices, li)
		}
	}
	allKeys := make([][]bool, p.nProduct)
	for i := 0; i < p.nProduct; i++ {
		allKeys[i] = make([]bool, p.nPeriod)
	}
	for _, location := range blockIndices {
		allKeys[location[0]][location[1]] = true
	}
	for imac, x := range allKeys {
		for iper, y := range x {
			if y {
				availableDuration = p.dPeriod * mutated.utilization[imac][iper]
				for _, prod := range mutated.mpInvInd[imac][iper] {
					if iper != 0 {
						if mutated.last[imac][iper-1] != prod {
							availableDuration -= p.chgOver[prod]
						}
					} else {
						availableDuration -= p.chgOver[prod]
					}
				}
				if availableDuration < 0.0 {
					for _, prod := range mutated.mpInvInd[imac][iper] {
						mutated.machineLayer[prod][iper] = -1
						mutated.lotsizeLayer[prod][iper] = 0
					}
					mutated.mpInvInd[imac][iper] = make([]int, 0)
					mutated.last[imac][iper] = -1
					break
				} else {
					amt := make([]float64, 0)
					for prodInd := 1; prodInd < len(mutated.mpInvInd[imac][iper]); prodInd++ { //divide the available duration into random intervals (production duration for each product in that [machine][period])
						amt = append(amt, rand.Float64())
					}
					sort.Float64s(amt)
					amt = append(amt, float64(1.0))
					amt = append([]float64{0.0}, amt...)
					for prodInd, prod := range mutated.mpInvInd[imac][iper] {
						mutated.lotsizeLayer[prod][iper] = int(float32((amt[prodInd+1] - amt[prodInd])) * availableDuration * p.socket[prod] / p.cycleTime[prod])
					}
				}
			}
		}
	}
	return mutated
}

//machineMutation changes randomly chosen machine layer entries up to a number of maxLoop, or until Bernoulli(sizeCoeff)=0
//if -1 entry is replaced by a new machine -> adds a random amount of lot size to newly added machines, then fits to max lot size
//if machine replaced by new machine -> add to new machine, fit the remaining to full in prev machine
//...same lot size retaining logic as applied everywhere
func (ch *Chromosome) machineMutation(sizeCoeff float32, p Problem, fillProb float32, maxLoop int) Chromosome {
	var mutated Chromosome
	mutated.copyChromosome(*ch)
	var availableDuration float32
	durationRatio := float32(0.0)
	var amt int
	nLoop := 0
	var iprod int
	var iper int
	var imac int
	var newMac int

	for rand.Float32() < sizeCoeff && nLoop < maxLoop {
		iprod = rand.Intn(p.nProduct)
		iper = rand.Intn(p.nPeriod)
		if mutated.machineLayer[iprod][iper] == -1 { //chose an empty cell; with fillProb a new cell will be added
			if rand.Float32() < fillProb {
				imac = p.mpMatch[iprod][rand.Intn(len(p.mpMatch[iprod]))]
				mutated.mpInvInd[imac][iper] = append(mutated.mpInvInd[imac][iper], iprod)
				mutated.machineLayer[iprod][iper] = imac
				mutated.lotsizeLayer[iprod][iper] = int(mutated.availability[imac][iper]*(1.0/float32(len(mutated.mpInvInd[imac][iper])))*(0.5+rand.Float32())) / int(p.cycleTime[iprod])
				mutated.lotsizeLayer[iprod][iper] = mutated.lotsizeLayer[iprod][iper] * int(p.socket[iprod])

				// assigned a randomized lot size to the newly added machineLayer entry. Now it needs repair to fit the available duration

				availableDuration = p.dPeriod * mutated.utilization[imac][iper]
				for _, prod := range mutated.mpInvInd[imac][iper] {
					if iper != 0 {
						if mutated.last[imac][iper-1] != prod {
							availableDuration -= p.chgOver[prod]
						}
					} else {
						availableDuration -= p.chgOver[prod]
					}
				}
				if availableDuration < 0.0 {
					for _, prod := range mutated.mpInvInd[imac][iper] {
						mutated.machineLayer[prod][iper] = -1
						mutated.lotsizeLayer[prod][iper] = 0
					}
					mutated.mpInvInd[imac][iper] = make([]int, 0)
					mutated.last[imac][iper] = -1
					break
				} else {
					durationRatio = 0.0
					for _, prod := range mutated.mpInvInd[imac][iper] {
						durationRatio += (float32(mutated.lotsizeLayer[prod][iper]) * float32(p.cycleTime[prod]) / float32(p.socket[prod]))
					}
					durationRatio = availableDuration / durationRatio //durationRatio is the available duration change by percentage in the mp
					for _, prod := range mutated.mpInvInd[imac][iper] {
						amt = mutated.lotsizeLayer[prod][iper]
						amt = int(float32(amt) * durationRatio)
						amt = amt - amt%int(p.socket[prod])
						mutated.lotsizeLayer[prod][iper] = amt
					}
				}
			} else {
				continue
			}
		} else {
			if rand.Float32() < fillProb { //changing the preexisting machine to a new one
				if len(p.mpMatch[iprod]) == 1 {
					continue //if no other machine alternative exists for that product, skip
				} else { 
					tmp := make([]int, 0)
					for _, v := range mutated.mpInvInd[imac][iper] {
						if v != iprod {
							tmp = append(tmp, v)
						}
					}
					mutated.mpInvInd[imac][iper] = make([]int, len(tmp))
					copy(mutated.mpInvInd[imac][iper], tmp)
					for v := rand.Intn(len(p.mpMatch[iprod])); p.mpMatch[iprod][v] == imac; v = rand.Intn(len(p.mpMatch[iprod])) {
						newMac = p.mpMatch[iprod][v]
					}
					mutated.mpInvInd[newMac][iper] = append(mutated.mpInvInd[newMac][iper],iprod)
					mutated.machineLayer[iprod][iper] = newMac

					if(len(mutated.mpInvInd[imac][iper])!=0){ //fit the previous period lot sizes to available duration
						availableDuration = p.dPeriod * mutated.utilization[imac][iper]
						for _, prod := range mutated.mpInvInd[imac][iper] {
							if iper != 0 {
								if mutated.last[imac][iper-1] != prod {
									availableDuration -= p.chgOver[prod]
								}
							} else {
								availableDuration -= p.chgOver[prod]
							}
						}
	
						durationRatio = 0.0
						for _, prod := range mutated.mpInvInd[imac][iper] {
							durationRatio += (float32(mutated.lotsizeLayer[prod][iper]) * float32(p.cycleTime[prod]) / float32(p.socket[prod]))
						}
						durationRatio = availableDuration / durationRatio
						for _, prod := range mutated.mpInvInd[imac][iper] {
							amt = mutated.lotsizeLayer[prod][iper]
							amt = int(float32(amt) * durationRatio)
							amt = amt - amt%int(p.socket[prod])
							mutated.lotsizeLayer[prod][iper] = amt
						}
					}
					//fit the lot sizes to available duration in the newly added machine's mp
					availableDuration = p.dPeriod * mutated.utilization[newMac][iper]
					for _, prod := range mutated.mpInvInd[newMac][iper] {
						if iper != 0 {
							if mutated.last[newMac][iper-1] != prod {
								availableDuration -= p.chgOver[prod]
							}
						} else {
							availableDuration -= p.chgOver[prod]
						}
					}
					if availableDuration < 0.0 {
						for _, prod := range mutated.mpInvInd[newMac][iper] {
							mutated.machineLayer[prod][iper] = -1
							mutated.lotsizeLayer[prod][iper] = 0
						}
						mutated.mpInvInd[newMac][iper] = make([]int, 0)
						mutated.last[newMac][iper] = -1
						break
					} else {
						durationRatio = 0.0
						for _, prod := range mutated.mpInvInd[newMac][iper] {
							durationRatio += (float32(mutated.lotsizeLayer[prod][iper]) * float32(p.cycleTime[prod]) / float32(p.socket[prod]))
						}
						durationRatio = availableDuration / durationRatio //durationRatio is the available duration change by percentage in the mp
						for _, prod := range mutated.mpInvInd[newMac][iper] {
							amt = mutated.lotsizeLayer[prod][iper]
							amt = int(float32(amt) * durationRatio)
							amt = amt - amt%int(p.socket[prod])
							mutated.lotsizeLayer[prod][iper] = amt
						}
					}
				}
			} else {//removing a preexisting machine to idle pp
				tmp := make([]int, 0)
				for _, v := range mutated.mpInvInd[imac][iper] {
					if v != iprod {
						tmp = append(tmp, v)
					}
				}
				mutated.mpInvInd[imac][iper] = make([]int, len(tmp))
				copy(mutated.mpInvInd[imac][iper], tmp)
				mutated.machineLayer[iprod][iper]=-1
				mutated.lotsizeLayer[iprod][iper]=0
				if(len(mutated.mpInvInd[imac][iper])!=0){ //fit the previous period lot sizes to available duration
					availableDuration = p.dPeriod * mutated.utilization[imac][iper]
					for _, prod := range mutated.mpInvInd[imac][iper] {
						if iper != 0 {
							if mutated.last[imac][iper-1] != prod {
								availableDuration -= p.chgOver[prod]
							}
						} else {
							availableDuration -= p.chgOver[prod]
						}
					}

					durationRatio = 0.0
					for _, prod := range mutated.mpInvInd[imac][iper] {
						durationRatio += (float32(mutated.lotsizeLayer[prod][iper]) * float32(p.cycleTime[prod]) / float32(p.socket[prod]))
					}
					durationRatio = availableDuration / durationRatio
					for _, prod := range mutated.mpInvInd[imac][iper] {
						amt = mutated.lotsizeLayer[prod][iper]
						amt = int(float32(amt) * durationRatio)
						amt = amt - amt%int(p.socket[prod])
						mutated.lotsizeLayer[prod][iper] = amt
					}
				}
			}
		}
	}
	return mutated
}