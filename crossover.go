package main

import (
	"math/rand"
)

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

func (ch *Chromosome) blockCrossover(sizeCoeff float32, p Problem, c2 Chromosome, lsConservation float32) []Chromosome {
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
	affectedMac[-1] = true
	delete(affectedMac, -1)

	var child_1 Chromosome
	var child_2 Chromosome
	child_1.copyChromosome(*ch)
	child_2.copyChromosome(c2)
	prev_1 := make([][][]int, p.nMachine)
	pnew_1 := make([][][]int, p.nMachine)
	prev_2 := make([][][]int, p.nMachine)
	pnew_2 := make([][][]int, p.nMachine)
	for i := 0; i < p.nMachine; i++ {
		prev_1[i] = make([][]int, p.nPeriod)
		pnew_1[i] = make([][]int, p.nPeriod)
		prev_2[i] = make([][]int, p.nPeriod)
		pnew_2[i] = make([][]int, p.nPeriod)
		for j := 0; j < p.nPeriod; j++ {
			prev_1[i][j] = make([]int, len(child_1.mpInvInd[i][j]))
			copy(prev_1[i][j], child_1.mpInvInd[i][j])
			pnew_1[i][j] = []int{}
			prev_2[i][j] = make([]int, len(child_2.mpInvInd[i][j]))
			copy(prev_2[i][j], child_2.mpInvInd[i][j])
			pnew_2[i][j] = []int{}
		}
	}

	var product, period int
	var tmp []int
	
	for _, location := range blockIndices { //swap the mpInvInd, machineLayer, and lotsizeLayer entries
		product = location[0]
		period = location[1]

		if child_1.machineLayer[product][period] != -1 { //delete product from mpInvInd if any product was on that location
			tmp = make([]int, 0)
			for _, v := range child_1.mpInvInd[child_1.machineLayer[product][period]][period] {
				if v != product{
					tmp = append(tmp, v)
				}
			}
			child_1.mpInvInd[child_1.machineLayer[product][period]][period] = make([]int, len(tmp))
			prev_1[child_1.machineLayer[product][period]][period]=make([]int,len(tmp))

			copy(child_1.mpInvInd[child_1.machineLayer[product][period]][period], tmp)
			copy(prev_1[child_1.machineLayer[product][period]][period], tmp)
		}
		if c2.machineLayer[product][period] != -1 { //add new product to the machine if any introduced
			child_1.mpInvInd[c2.machineLayer[product][period]][period] = append(child_1.mpInvInd[c2.machineLayer[product][period]][period], product)
			pnew_1[c2.machineLayer[product][period]][period] = append(pnew_1[c2.machineLayer[product][period]][period], product)
		}

		if child_2.machineLayer[product][period] != -1 { //delete product from mpInvInd if any product produced on that location
			tmp = make([]int, len(child_2.mpInvInd[child_2.machineLayer[product][period]][period])-1)
			for _, v := range child_2.mpInvInd[child_2.machineLayer[product][period]][period] {
				if v != product{
					tmp = append(tmp, v)
				}
			}
			child_2.mpInvInd[child_2.machineLayer[product][period]][period] = make([]int, len(tmp))
			prev_2[child_2.machineLayer[product][period]][period]=make([]int,len(tmp))

			copy(child_2.mpInvInd[child_2.machineLayer[product][period]][period], tmp)
			copy(prev_2[child_2.machineLayer[product][period]][period], tmp)
		}
		if ch.machineLayer[product][period] != -1 {
			child_2.mpInvInd[ch.machineLayer[product][period]][period] = append(child_2.mpInvInd[ch.machineLayer[product][period]][period], product)
			pnew_2[ch.machineLayer[product][period]][period] = append(pnew_2[ch.machineLayer[product][period]][period], product)
		}
		//now swap machine and lot size layers; the required adjustments for these will be done below in the repair procedure
		child_1.machineLayer[product][period] = c2.machineLayer[product][period]
		child_2.machineLayer[product][period] = ch.machineLayer[product][period]
		child_1.lotsizeLayer[product][period] = c2.lotsizeLayer[product][period]
		child_2.lotsizeLayer[product][period] = ch.lotsizeLayer[product][period]
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
	var mpList []int
	var alternatives []int
	var durationRatio float32
	var amt int
	for mac, _ := range affectedMac { //restructure the last layer. if still compatible, keep. if not, choose another last for the mp. if unable to do so, -1
		for t := minAP; t <= maxAP; t++ {
			durationRatio = 0.0
			if t == p.nPeriod-1 {
				child_1.last[mac][t] = -1
			} else if isin(c2.last[mac][t], child_1.mpInvInd[mac][t]) && isin(c2.last[mac][t], child_1.mpInvInd[mac][t+1]) {
				child_1.last[mac][t] = c2.last[mac][t] //if product still produced in m,p and m,p+1 in child_1, it is the last
			} else {
				alternatives=make([]int,len(inArr(child_1.mpInvInd[mac][t], child_1.mpInvInd[mac][t+1])))
				copy(alternatives,inArr(child_1.mpInvInd[mac][t], child_1.mpInvInd[mac][t+1]))
				if len(alternatives) != 0 {
					child_1.last[mac][t] = alternatives[rand.Intn(len(alternatives))]
				} else {
					child_1.last[mac][t] = -1
				} //if nothing produced on mp, last layer is -1
			}
			availableDuration = p.dPeriod * child_1.utilization[mac][t]
			keepFlag = false
			mpList=make([]int, len(child_1.mpInvInd[mac][t]))
			copy(mpList, child_1.mpInvInd[mac][t])
			if(len(mpList)==0){
				continue;
			}
			for prodInd := 0; prodInd < len(mpList); prodInd++ {
				if t != 0 {
					if mpList[prodInd] != child_1.last[mac][t-1] {
						availableDuration -= p.chgOver[mpList[prodInd]]
					} else {
						if t != 1 {
							if child_1.last[mac][t-2] == mpList[prodInd] { //if the mold is kept in the previous period as well (from t-2 to t-1), then there must be only the kept product on t-1 on that machine. Otherwise can't keep the same mold again
								if len(child_1.mpInvInd[mac][t-1]) == 1 {
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
				} else {
					availableDuration -= p.chgOver[mpList[prodInd]]
				} //at t=0 can't save time from previous period
			}
			if availableDuration < 0.0 {
				for prodInd := 0; prodInd < len(mpList); prodInd++ {
					child_1.machineLayer[mpList[prodInd]][t] = -1
					child_1.lotsizeLayer[mpList[prodInd]][t] = 0
				}
				child_1.mpInvInd[mac][t] = make([]int, 0)
				break
			} else {
				if !keepFlag {
					if t != 0 {
						child_1.last[mac][t-1] = -1 //if the product mold kept in previous period is not used, then don't keep
					}
				}
				if (len(prev_1[mac][t]) == 0 || len(pnew_1[mac][t]) == 0) && len(prev_1[mac][t])+len(pnew_1[mac][t]) != 0 {
					for prodInd := 0; prodInd < len(mpList); prodInd++ {
						durationRatio += (float32(child_1.lotsizeLayer[mpList[prodInd]][t]) * float32(p.cycleTime[mpList[prodInd]]) / float32(p.socket[mpList[prodInd]]))
					}
					durationRatio = availableDuration / durationRatio //durationRatio is the available duration change in the mp
					for prodInd := 0; prodInd < len(mpList); prodInd++ {
						amt = child_1.lotsizeLayer[mpList[prodInd]][t]
						amt = int(float32(amt) * durationRatio)
						amt = amt - amt%int(p.socket[mpList[prodInd]])
						child_1.lotsizeLayer[mpList[prodInd]][t] = amt
					}
				} else if len(pnew_1[mac][t])+len(prev_1[mac][t]) != 0 {
					for _, prod := range prev_1[mac][t] {
						durationRatio += (float32(child_1.lotsizeLayer[prod][t]) * float32(p.cycleTime[prod]) / float32(p.socket[prod]))
					}
					durationRatio = availableDuration * lsConservation / durationRatio //durationRatio is the available duration change in the mp
					for _, prod := range prev_1[mac][t] {
						amt = child_1.lotsizeLayer[prod][t]
						amt = int(float32(amt) * durationRatio)
						amt = amt - amt%int(p.socket[prod])
						child_1.lotsizeLayer[prod][t] = amt
					}
					for _, prod := range pnew_1[mac][t] {
						durationRatio += (float32(child_1.lotsizeLayer[prod][t]) * float32(p.cycleTime[prod]) / float32(p.socket[prod]))
					}
					durationRatio = availableDuration * (1.0 - lsConservation) / durationRatio //durationRatio is the available duration change in the mp
					for _, prod := range pnew_1[mac][t] {
						amt = child_1.lotsizeLayer[prod][t]
						amt = int(float32(amt) * durationRatio)
						amt = amt - amt%int(p.socket[prod])
						child_1.lotsizeLayer[prod][t] = amt
					}
				}
			}
		}
	}
	for mac := range affectedMac { //restructure the last layer. if still compatible, keep. if not, choose another last for the mp. if unable to do so, -1
		for t := minAP; t <= maxAP; t++ { //adjust the lot size layer too, in order to fit max usable time in each mp
			//do the same for second child
			durationRatio = 0.0
			if t == p.nPeriod-1 {
				child_2.last[mac][t] = -1
			} else if isin(ch.last[mac][t], child_2.mpInvInd[mac][t]) && isin(ch.last[mac][t], child_2.mpInvInd[mac][t+1]) {
				child_2.last[mac][t] = ch.last[mac][t] //if product still produced in m,p and m,p+1 in child_2, it is the last
			} else {
				alternatives = inArr(child_2.mpInvInd[mac][t], child_2.mpInvInd[mac][t+1])
				if len(alternatives) != 0 {
					child_2.last[mac][t] = alternatives[rand.Intn(len(alternatives))]
				} else {
					child_2.last[mac][t] = -1
				} //if nothing produced on mp, last layer is -1
			}
			availableDuration = p.dPeriod * child_2.utilization[mac][t]
			keepFlag = false
			mpList=make([]int, len(child_2.mpInvInd[mac][t]))
			copy(mpList, child_2.mpInvInd[mac][t])
			if(len(mpList)==0){
				continue;
			}
			for prodInd := 0; prodInd < len(mpList); prodInd++ {
				if t != 0 {
					if mpList[prodInd] != child_2.last[mac][t-1] {
						availableDuration -= p.chgOver[mpList[prodInd]]
					} else {
						if t != 1 {
							if child_2.last[mac][t-2] == mpList[prodInd] { //if the mold is kept in the previous period as well (from t-2 to t-1), then there must be only the kept product on t-1 on that machine. Otherwise can't keep the same mold again
								if len(child_2.mpInvInd[mac][t-1]) == 1 {
									keepFlag = true
								} else {
									availableDuration -= p.chgOver[mpList[prodInd]]
								}
							} else {
								keepFlag = true
							}
						} else {
							keepFlag = true //no need for t-2 check in period 1
						} //the kept product mold from previous period is used, therefore it does not require changeover
					}
				} else {
					availableDuration -= p.chgOver[mpList[prodInd]]
				} //at t=0 can't save time from previous period
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
				if (len(prev_2[mac][t]) == 0 || len(pnew_2[mac][t]) == 0) && len(prev_2[mac][t])+len(pnew_2[mac][t]) != 0 {
					for prodInd := 0; prodInd < len(mpList); prodInd++ {
						durationRatio += (float32(child_2.lotsizeLayer[mpList[prodInd]][t]) * float32(p.cycleTime[mpList[prodInd]]) / float32(p.socket[mpList[prodInd]]))
					}
					durationRatio = availableDuration / durationRatio //durationRatio is the available duration change in the mp
					for prodInd := 0; prodInd < len(mpList); prodInd++ {
						amt = child_2.lotsizeLayer[mpList[prodInd]][t]
						amt = int(float32(amt) * durationRatio)
						amt = amt - amt%int(p.socket[mpList[prodInd]])
						child_2.lotsizeLayer[mpList[prodInd]][t] = amt
					}
				} else if len(prev_2[mac][t])+len(pnew_2[mac][t]) != 0 {
					for _, prod := range prev_2[mac][t] {
						durationRatio += (float32(child_2.lotsizeLayer[prod][t]) * float32(p.cycleTime[prod]) / float32(p.socket[prod]))
					}
					durationRatio = availableDuration * lsConservation / durationRatio //durationRatio is the available duration change in the mp
					for _, prod := range prev_2[mac][t] {
						amt = child_2.lotsizeLayer[prod][t]
						amt = int(float32(amt) * durationRatio)
						amt = amt - amt%int(p.socket[prod])
						child_2.lotsizeLayer[prod][t] = amt
					}
					for _, prod := range pnew_2[mac][t] {
						durationRatio += (float32(child_2.lotsizeLayer[prod][t]) * float32(p.cycleTime[prod]) / float32(p.socket[prod]))
					}
					durationRatio = availableDuration * (1.0 - lsConservation) / durationRatio //durationRatio is the available duration change in the mp
					for _, prod := range pnew_2[mac][t] {
						amt = child_2.lotsizeLayer[prod][t]
						amt = int(float32(amt) * durationRatio)
						amt = amt - amt%int(p.socket[prod])
						child_2.lotsizeLayer[prod][t] = amt
					}
				}
			}
		}
	}

	return []Chromosome{child_1, child_2}

}
