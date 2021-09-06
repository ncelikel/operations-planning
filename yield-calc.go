package main

func (ch *Chromosome) yieldAll(p Problem) []float32 {
	cost := []float32{0.0, 0.0, 0.0}
	curProduction := make([][]int, p.nPeriod)
	cumProduction := make([][]int, p.nPeriod)
	for period := 0; period < p.nPeriod; period++ {
		curProduction[period] = make([]int, p.nProduct)
		cumProduction[period] = make([]int, p.nProduct)
	}
	for product := 0; product < p.nProduct; product++ {
		curProduction[0][product] = ch.lotsizeLayer[product][0]
		cumProduction[0][product] = ch.lotsizeLayer[product][0]
		if curProduction[0][product] <= p.curDemand[0][product] {
			cost[0] += float32(p.curDemand[0][product]) - float32(curProduction[0][product])
		} else {
			cost[1] += float32(curProduction[0][product]) - float32(p.curDemand[0][product])
		}

	}
	for machine := 0; machine < p.nMachine; machine++ {
		for _, val := range ch.mpInvInd[machine][0] {
			cost[2] += p.chgOver[val]
		}
	}
	for period := 0; period < p.nPeriod; period++ {
		for product := 0; product < p.nProduct; product++ {
			curProduction[period][product] = ch.lotsizeLayer[product][period]
			cumProduction[period][product] = curProduction[period][product] + cumProduction[period-1][product]
			if cumProduction[period][product] <= p.cumDemand[period][product] {
				cost[0] += float32(p.cumDemand[period][product]) - float32(cumProduction[period][product])
			} else {
				cost[1] += float32(curProduction[period][product]) - float32(p.curDemand[period][product])
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
