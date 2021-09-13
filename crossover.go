package main

import (
	"fmt"
	"math/rand"
)

func deleteSliceElement(a []int, i int) []int {
	a[i] = a[len(a)-1]
	a = a[:len(a)-1]
	return a
}
func isin(a int, b []int) bool {
	isinFlag := false
	for i := 0; i < len(b); i++ {
		if b[i] == a {
			isinFlag = true
			break
		}
	}
	return isinFlag
}
func inArr(a []int, b []int) []int {
	elements := make([]int, 0)
	for j := 0; j < len(a); j++ {
		if isin(a[j], b) {
			elements = append(elements, a[j])
		}
	}
	return elements
}

func (ch *Chromosome) blockCrossover(sizeCoeff float32, p Problem, c2 Chromosome) []Chromosome {
	var blockIndices [][]int
	blockIndices = append(blockIndices, []int{rand.Intn(p.nProduct), rand.Intn(p.nPeriod)})
	li := blockIndices[0]
	spreadProbability := 1.0 - (1.0 / (float32(p.nPeriod) * float32(p.nProduct) * sizeCoeff))
	dxy := []int{-1, 0, 1}
	for rand.Float32() < spreadProbability { //randomly spread to locations around
		li = []int{li[0] + dxy[rand.Intn(3)], li[0] + dxy[rand.Intn(3)]}
		if li[0] >= 0 && li[1] >= 0 && li[0] < p.nProduct && li[1] < p.nPeriod {
			blockIndices = append(blockIndices, li)
		}
	}
	allKeys := make([][]bool, p.nProduct)
	for i := 0; i < p.nProduct; i++ {
		allKeys[i] = make([]bool, p.nPeriod)
	}
	affectedMac := make(map[int]bool)
	affectedPer := make(map[int]bool)
	for _, location := range blockIndices {
		allKeys[location[0]][location[1]] = true
	}
	blockIndices = make([][]int, 0) //locations where change will be made
	for ix, x := range allKeys {
		for iy, y := range x {
			if y {
				blockIndices = append(blockIndices, []int{ix, iy})
				affectedMac[ch.machineLayer[ix][iy]] = true
				affectedMac[c2.machineLayer[ix][iy]] = true
				affectedPer[iy] = true
				if iy != 0 {
					affectedPer[iy-1] = true
				} else if iy != p.nPeriod-1 {
					affectedPer[iy+1] = true
				}
			}
		}

	}
	fmt.Println(blockIndices)
	affectedMac[-1] = true
	delete(affectedMac, -1)

	var child_1 Chromosome
	var child_2 Chromosome
	child_1.copyChromosome(*ch)
	child_2.copyChromosome(c2)

	var product, period int
	for _, location := range blockIndices { //swap the mpInvInd, machineLayer, and lotsizeLayer entries
		product = location[0]
		period = location[1]
		if child_1.machineLayer[product][period] != -1 { //delete product from mpInvInd if any product produced on that location
			for i, v := range child_1.mpInvInd[child_1.machineLayer[product][period]][period] {
				if v == product {
					child_1.mpInvInd[child_1.machineLayer[product][period]][period] = deleteSliceElement(child_1.mpInvInd[child_1.machineLayer[product][period]][period], i)
					break
				}
			}
		}
		if c2.machineLayer[product][period] != -1 { //add new product to the machine if any introduced
			child_1.mpInvInd[c2.machineLayer[product][period]][period] = append(child_1.mpInvInd[c2.machineLayer[product][period]][period], product)
		}

		if child_2.machineLayer[product][period] != -1 {
			for i, v := range child_2.mpInvInd[child_2.machineLayer[product][period]][period] {
				if v == product {
					child_2.mpInvInd[child_2.machineLayer[product][period]][period] = deleteSliceElement(child_2.mpInvInd[child_2.machineLayer[product][period]][period], i)
					break
				}
			}
		}
		if ch.machineLayer[product][period] != -1 {
			child_2.mpInvInd[ch.machineLayer[product][period]][period] = append(child_2.mpInvInd[ch.machineLayer[product][period]][period], product)
		}
		//now swap machine and lot size layers; the required adjustments for these will be done below in the repair procedure
		child_1.machineLayer[product][period] = c2.machineLayer[product][period]
		child_2.machineLayer[product][period] = ch.machineLayer[product][period]
		child_1.lotsizeLayer[product][period] = c2.machineLayer[product][period]
		child_2.lotsizeLayer[product][period] = ch.machineLayer[product][period]
	}
	minAP := p.nPeriod
	maxAP := 0
	for per, _ := range affectedPer { //repair the affected periods. min affected to max affected
		if per < minAP {
			minAP = per
		}
		if per > maxAP {
			maxAP = per
		}
	}
	var keepFlag bool
	var availableDuration float32
	mpList := make([]int, 0)
	alternatives := make([]int, 0)
	var durationRatio float32
	var amt int
	for mac, _ := range affectedMac { //restructure the last layer. if still compatible, keep. if not, choose another last for the mp. if unable to do so, -1
		for t := minAP; t <= maxAP; t++ {
			durationRatio = 0.0
			if t != p.nPeriod-1 {
				if isin(c2.last[mac][t], child_1.mpInvInd[mac][t]) && isin(c2.last[mac][t], child_1.mpInvInd[mac][t+1]) {
					child_1.last[mac][t] = c2.last[mac][t] //if product still produced in m,p and m,p+1 in child_1, it is the last
				} else {
					alternatives = inArr(child_1.mpInvInd[mac][t], child_1.mpInvInd[mac][t+1])
					if len(alternatives) != 0 {
						child_1.last[mac][t] = alternatives[rand.Intn(len(alternatives))]
					} else {
						child_1.last[mac][t] = -1
					} //if nothing produced on mp, last layer is -1
				}
				availableDuration = p.dPeriod
				keepFlag = false
				mpList = child_1.mpInvInd[mac][t]
				for prodInd := 0; prodInd < len(mpList); prodInd++ {
					if t != 0 {
						if mpList[prodInd] != child_1.last[mac][t-1] {
							availableDuration -= p.chgOver[mpList[prodInd]]
						}
					} else {
						if t > 1 {
							if ch.last[mac][t-2] == mpList[prodInd] { //if the mold is kept in the previous period as well (from t-2 to t-1), then there must be only the kept product on t-1 on that machine. Otherwise can't keep the same mold again
								if len(mpList) == 1 {
									keepFlag = true
								} else {
									availableDuration -= p.chgOver[mpList[prodInd]]
								}
							} else {
								keepFlag = true
							}
						} else {
							keepFlag = true
						} //the kept product mold from previous period is used, therefore it does not require changeover
					}
				}
				if availableDuration < 0.0 {
					for prodInd := 0; prodInd < len(mpList); prodInd++ {
						child_1.machineLayer[mpList[prodInd]][t] = -1
						child_1.lotsizeLayer[mpList[prodInd]][t] = 0
					}
					child_1.mpInvInd[mac][t] = make([]int, 0)
				} else {
					if !keepFlag {
						if t != 0 {
							child_1.last[mac][t-1] = -1 //if the product mold kept in previous period is not used, then don't keep
						}
					}
					for prodInd := 0; prodInd < len(mpList); prodInd++ {
						durationRatio += (float32(child_1.lotsizeLayer[mpList[prodInd]][t]) * float32(p.cycleTime[mpList[prodInd]]) / float32(p.socket[mpList[prodInd]]))
					}
					durationRatio = availableDuration / durationRatio
					for prodInd := 0; prodInd < len(mpList); prodInd++ { //increase or decrease the products in mp by the same percentage
						amt = child_1.lotsizeLayer[mpList[prodInd]][t]
						amt = int(float32(amt) * durationRatio)
						amt = amt - amt%int(p.socket[mpList[prodInd]])
						child_1.lotsizeLayer[mpList[prodInd]][t] = amt
					}
				}

				//do the same for second child

				if isin(ch.last[mac][t], child_2.mpInvInd[mac][t]) && isin(ch.last[mac][t], child_2.mpInvInd[mac][t+1]) {
					child_2.last[mac][t] = ch.last[mac][t] //if product still produced in m,p and m,p+1 in child_2, it is the last
				} else {
					alternatives = inArr(child_2.mpInvInd[mac][t], child_2.mpInvInd[mac][t+1])
					if len(alternatives) != 0 {
						child_2.last[mac][t] = alternatives[rand.Intn(len(alternatives))]
					} else {
						child_2.last[mac][t] = -1
					} //if nothing produced on mp, last layer is -1
				}
				availableDuration = p.dPeriod
				keepFlag = false
				mpList = child_2.mpInvInd[mac][t]
				for prodInd := 0; prodInd < len(mpList); prodInd++ {
					if t != 0 {
						if mpList[prodInd] != child_2.last[mac][t-1] {
							availableDuration -= p.chgOver[mpList[prodInd]]
						}
					} else {
						if t > 1 {
							if ch.last[mac][t-2] == mpList[prodInd] { //if the mold is kept in the previous period as well (from t-2 to t-1), then there must be only the kept product on t-1 on that machine. Otherwise can't keep the same mold again
								if len(mpList) == 1 {
									keepFlag = true
								} else {
									availableDuration -= p.chgOver[mpList[prodInd]]
								}
							} else {
								keepFlag = true
							}
						} else {
							keepFlag = true
						} //the kept product mold from previous period is used, therefore it does not require changeover
					}
				}
				if availableDuration < 0.0 {
					for prodInd := 0; prodInd < len(mpList); prodInd++ {
						child_2.machineLayer[mpList[prodInd]][t] = -1
						child_2.lotsizeLayer[mpList[prodInd]][t] = 0
					}
					child_2.mpInvInd[mac][t] = make([]int, 0)
				} else {
					if !keepFlag {
						if t != 0 {
							child_2.last[mac][t-1] = -1 //if the product mold kept in previous period is not used, then don't keep
						}
					}
					for prodInd := 0; prodInd < len(mpList); prodInd++ {
						durationRatio += (float32(child_2.lotsizeLayer[mpList[prodInd]][t]) * float32(p.cycleTime[mpList[prodInd]]) / float32(p.socket[mpList[prodInd]]))
					}
					durationRatio = availableDuration / durationRatio
					for prodInd := 0; prodInd < len(mpList); prodInd++ { //increase or decrease the products in mp by the same percentage
						amt = child_2.lotsizeLayer[mpList[prodInd]][t]
						amt = int(float32(amt) * durationRatio)
						amt = amt - amt%int(p.socket[mpList[prodInd]])
						child_2.lotsizeLayer[mpList[prodInd]][t] = amt
					}
				}

			} else {
				child_1.last[mac][t] = -1
				child_2.last[mac][t] = -1
			}
		}
	}

	if minAP != 0 {
		minAP -= 1
	}
	if maxAP != p.nPeriod-1 {
		maxAP += 1
	}
	for mac, _ := range affectedMac {
		for t := minAP; t <= maxAP; t++ {
			durationRatio = 0.0
			availableDuration = p.dPeriod
			mpList = child_1.mpInvInd[mac][t]
			for prodInd := 0; prodInd < len(mpList); prodInd++ {
				if t == 0 {
					availableDuration -= p.chgOver[mpList[prodInd]]
				} else {
					if child_1.last[mac][t-1] != mpList[prodInd] {
						availableDuration -= p.chgOver[mpList[prodInd]]
					}
				}
			}
			if availableDuration < 0.0 {
				for prodInd := 0; prodInd < len(mpList); prodInd++ {
					child_1.machineLayer[mpList[prodInd]][t] = -1
					child_1.lotsizeLayer[mpList[prodInd]][t] = 0
					if child_1.last[mac][t] == mpList[prodInd] {
						child_1.last[mac][t] = -1
					}
				}
				child_1.mpInvInd[mac][t] = make([]int, 0)
			} else {
				for prodInd := 0; prodInd < len(mpList); prodInd++ {
					durationRatio += (float32(child_1.lotsizeLayer[mpList[prodInd]][t]) * float32(p.cycleTime[mpList[prodInd]]) / float32(p.socket[mpList[prodInd]]))
					durationRatio = availableDuration / durationRatio
					for prodInd := 0; prodInd < len(mpList); prodInd++ { //increase or decrease the products in mp by the same percentage
						amt = child_1.lotsizeLayer[mpList[prodInd]][t]
						amt = int(float32(amt) * durationRatio)
						amt = amt - amt%int(p.socket[mpList[prodInd]])
						child_1.lotsizeLayer[mpList[prodInd]][t] = amt
					}
				}
			}

			//same for child_2

			durationRatio = 0.0
			availableDuration = p.dPeriod
			mpList = child_2.mpInvInd[mac][t]
			for prodInd := 0; prodInd < len(mpList); prodInd++ {
				if t == 0 {
					availableDuration -= p.chgOver[mpList[prodInd]]
				} else {
					if child_2.last[mac][t-1] != mpList[prodInd] {
						availableDuration -= p.chgOver[mpList[prodInd]]
					}
				}
			}
			if availableDuration < 0.0 {
				for prodInd := 0; prodInd < len(mpList); prodInd++ {
					child_2.machineLayer[mpList[prodInd]][t] = -1
					child_2.lotsizeLayer[mpList[prodInd]][t] = 0
					if child_2.last[mac][t] == mpList[prodInd] {
						child_2.last[mac][t] = -1
					}
				}
				child_2.mpInvInd[mac][t] = make([]int, 0)
			} else {
				for prodInd := 0; prodInd < len(mpList); prodInd++ {
					durationRatio += (float32(child_2.lotsizeLayer[mpList[prodInd]][t]) * float32(p.cycleTime[mpList[prodInd]]) / float32(p.socket[mpList[prodInd]]))
					durationRatio = availableDuration / durationRatio
					for prodInd := 0; prodInd < len(mpList); prodInd++ { //increase or decrease the products in mp by the same percentage
						amt = child_2.lotsizeLayer[mpList[prodInd]][t]
						amt = int(float32(amt) * durationRatio)
						amt = amt - amt%int(p.socket[mpList[prodInd]])
						child_2.lotsizeLayer[mpList[prodInd]][t] = amt
					}
				}
			}
		}
	}
	return []Chromosome{child_1, child_2}

}
