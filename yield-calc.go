package main

func (ch *Chromosome) yieldAll(p Problem) []float32 {
	cost := []float32{0.0, 0.0, 0.0}
	ch.curProduction = make([][]int, p.nProduct)
	ch.cumProduction = make([][]int, p.nProduct)
	for product := 0; product < p.nProduct; product++ {
		ch.curProduction[product] = make([]int, p.nPeriod)
		ch.cumProduction[product] = make([]int, p.nPeriod)
	}
	for product := 0; product < p.nProduct; product++ {
		ch.curProduction[product][0] = ch.lotsizeLayer[product][0]
		ch.cumProduction[product][0] = ch.lotsizeLayer[product][0]
		if ch.curProduction[product][0] <= p.curDemand[product][0] {
			cost[0] += float32(p.curDemand[product][0]) - float32(ch.curProduction[product][0])
		} else {
			cost[1] += float32(ch.curProduction[product][0]) - float32(p.curDemand[product][0])
		}

	}
	for machine := 0; machine < p.nMachine; machine++ {
		for _, val := range ch.mpInvInd[machine][0] {
			cost[2] += p.chgOver[val]
		}
	}
	for period := 1; period < p.nPeriod; period++ {
		for product := 0; product < p.nProduct; product++ {
			ch.curProduction[product][period] = ch.lotsizeLayer[product][period]
			ch.cumProduction[product][period] = ch.curProduction[product][period] + ch.cumProduction[product][period-1]
			if ch.cumProduction[product][period] <= p.cumDemand[product][period] {
				cost[0] += float32(p.cumDemand[product][period]) - float32(ch.cumProduction[product][period])
			} else {
				cost[1] += float32(ch.cumProduction[product][period]) - float32(p.cumDemand[product][period])
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
	return cost
}

/*func (ch *Chromosome) yieldChange(changedIndices[][][]int,newMac[][]int,newLS[][]int) []float32{
	return
}*/
//calculate difference in yield and objective by giving only the changed indices; in order to increase performance