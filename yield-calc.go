package main

func (ch *Chromosome) yieldAll(p Problem) []float32 {
	cost := []float32{0.0, 0.0, 0.0}
	curProduction := make([][]int, p.nProduct)
	cumProduction := make([][]int, p.nProduct)
	for product := 0; product < p.nProduct; product++ {
		curProduction[product] = make([]int, p.nPeriod)
		cumProduction[product] = make([]int, p.nPeriod)
	}
	for product := 0; product < p.nProduct; product++ {
		curProduction[product][0] = ch.lotsizeLayer[product][0]
		cumProduction[product][0] = ch.lotsizeLayer[product][0]
		if curProduction[product][0] <= p.curDemand[product][0] {
			cost[0] += float32(p.curDemand[product][0]) - float32(curProduction[product][0])
		} else {
			cost[1] += float32(curProduction[product][0]) - float32(p.curDemand[product][0])
		}

	}
	for machine := 0; machine < p.nMachine; machine++ {
		for _, val := range ch.mpInvInd[machine][0] {
			cost[2] += p.chgOver[val]
		}
	}
	for period := 1; period < p.nPeriod; period++ {
		for product := 0; product < p.nProduct; product++ {
			curProduction[product][period] = ch.lotsizeLayer[product][period]
			cumProduction[product][period] = curProduction[product][period] + cumProduction[product][period-1]
			if cumProduction[product][period] <= p.cumDemand[product][period] {
				cost[0] += float32(p.cumDemand[product][period]) - float32(cumProduction[product][period])
			} else {
				cost[1] += float32(cumProduction[product][period]) - float32(p.cumDemand[product][period])
			}
		}
		for machine := 0; machine < p.nMachine; machine++ {
			for _, val := range ch.mpInvInd[machine][period] {
				if ch.last[machine][period] != val {
					cost[2] += p.chgOver[val]
				}
			}
		}
	}
	ch.objective = cost[0] + cost[1] + cost[2]
	ch.curProduction = curProduction
	ch.cumProduction = cumProduction
	return cost
}

/*func (ch *Chromosome) yieldChange(changedIndices[][][]int,newMac[][]int,newLS[][]int) []float32{
	return
}*/
//calculate difference in yield and objective by giving only the changed indices; in order to increase performance