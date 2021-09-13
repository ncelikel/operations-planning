package main

import (
	//"fmt"
	"math/rand"
	"sort"
)

func (p *Problem) randInit(freq float32, util float32) Chromosome { //create a random Chromosome (solution) from the defined problem parameters
	ch := Chromosome{}
	ch.machineLayer = make([][]int, p.nProduct)
	ch.lotsizeLayer = make([][]int, p.nProduct)
	ch.last = make([][]int, p.nMachine)
	ch.mInvInd = make([][][]int, p.nMachine)
	ch.mpInvInd = make([][][]int, p.nMachine)
	ch.utilization = make([][]float32, p.nMachine)
	ch.availability = make([][]float32, p.nMachine)
	for i := 0; i < p.nMachine; i++ {
		ch.mInvInd[i] = make([][]int, 0)
		ch.mpInvInd[i] = make([][]int, p.nPeriod)
		ch.last[i] = make([]int, p.nPeriod)
		ch.utilization[i] = make([]float32, p.nPeriod)
		ch.availability[i] = make([]float32, p.nPeriod)
		for j := 0; j < p.nPeriod; j++ {
			ch.last[i][j] = -1
			ch.mpInvInd[i][j] = make([]int, 0)
		}
	}
	machine := 0
	for i := 0; i < p.nProduct; i++ {
		ch.machineLayer[i] = make([]int, p.nPeriod)
		ch.lotsizeLayer[i] = make([]int, p.nPeriod)
	}
	for j := 0; j < p.nPeriod; j++ {
		for i := 0; i < p.nProduct; i++ {
			if rand.Float32() < freq { //freq is the probability of giving a machine id to any [product][period] location in the machine layer
				machine = p.mpMatch[i][rand.Intn(len(p.mpMatch[i]))]
				ch.machineLayer[i][j] = machine
				ch.mInvInd[machine] = append(ch.mInvInd[machine], []int{i, j})
				ch.mpInvInd[machine][j] = append(ch.mpInvInd[machine][j], i)
			} else {
				ch.machineLayer[i][j] = -1
			}
		}
	}

	var availableDuration float32
	u:=float32(0.9)
	var keepFlag bool
	var mpList []int
	for mach := 0; mach < p.nMachine; mach++ {
		if len(ch.mpInvInd[mach][0]) != 0 { //choose random last product in period 0
			ch.last[mach][0] = ch.mpInvInd[mach][0][rand.Intn(len(ch.mpInvInd[mach][0]))]
		}
		amt := make([]float64, 0)
		availableDuration = p.dPeriod
		mpList = ch.mpInvInd[mach][0]
		for prodInd := 0; prodInd < len(mpList); prodInd++ { //subtract changeover durations to find net available time
			availableDuration -= p.chgOver[mpList[prodInd]]
		}
		if availableDuration < 0.0 { //if no available time remains; cancel production on that [machine][period]
			for prodInd := 0; prodInd < len(mpList); prodInd++ {
				ch.machineLayer[mpList[prodInd]][0] = -1
				ch.lotsizeLayer[mpList[prodInd]][0] = 0
			}
			ch.mpInvInd[mach][0] = make([]int, 0)
			break
		}
		ch.utilization[mach][0] = util
		ch.availability[mach][0] = availableDuration
		for prodInd := 1; prodInd < len(mpList); prodInd++ { //divide the available duration into random intervals (production duration for each product in that [machine][period])
			amt = append(amt, rand.Float64())
		}
		sort.Float64s(amt)
		amt = append(amt, float64(1.0))
		amt = append([]float64{0.0}, amt...)
		for prodInd := 1; prodInd < len(mpList)+1; prodInd++ {
			ch.lotsizeLayer[mpList[prodInd-1]][0] = int(float32((amt[prodInd] - amt[prodInd-1])) * u * availableDuration * p.socket[mpList[prodInd-1]] / p.cycleTime[mpList[prodInd-1]])
		}

		for period := 1; period < p.nPeriod; period++ {
			if len(ch.mpInvInd[mach][period]) != 0 { //choose random last product in period 0
				ch.last[mach][period] = ch.mpInvInd[mach][period][rand.Intn(len(ch.mpInvInd[mach][period]))]
			}
			amt := make([]float64, 0)
			availableDuration = p.dPeriod
			keepFlag = false
			mpList = ch.mpInvInd[mach][period]
			for prodInd := 0; prodInd < len(mpList); prodInd++ {
				if mpList[prodInd] != ch.last[mach][period-1] {
					availableDuration -= p.chgOver[mpList[prodInd]]
				} else {
					if period != 1 {
						if ch.last[mach][period-2] == mpList[prodInd] { //if the mold is kept in the previous period as well (from t-2 to t-1), then there must be only the kept product on t-1 on that machine. Otherwise can't keep the same mold again
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
					ch.machineLayer[mpList[prodInd]][period] = -1
					ch.lotsizeLayer[mpList[prodInd]][period] = 0
				}
				ch.mpInvInd[mach][period] = make([]int, 0)
				break
			}
			if !keepFlag {
				ch.last[mach][period-1] = -1 //if the product mold kept in previous period is not used, then don't keep
			}
			ch.utilization[mach][period] = util
			ch.availability[mach][period] = availableDuration
			for prodInd := 1; prodInd < len(mpList); prodInd++ {
				amt = append(amt, rand.Float64())
			}
			sort.Float64s(amt)
			amt = append(amt, float64(1.0))
			amt = append([]float64{0.0}, amt...)
			for prodInd := 1; prodInd < len(mpList)+1; prodInd++ {
				ch.lotsizeLayer[mpList[prodInd-1]][period] = int(float32((amt[prodInd] - amt[prodInd-1])) * u * availableDuration * p.socket[mpList[prodInd-1]] / p.cycleTime[mpList[prodInd-1]])
			}
		}
	}

	return ch
}
