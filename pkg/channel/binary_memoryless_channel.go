package channel

import (
	"math"
)

/*
Binary Input Stationary Memoryless Channel

input is {0, 1}, and output is Log likelihood ratio (ln(W(y|0)/W(y|1)), y is channel output).
*/
type BinaryMemorylessChannel interface {
	Channel([]int) []float64

	/* Evaluate error probability via density evolution
	 length is code length
	 index start with 1

	WARNING: This method is not thread safe.
	*/
	CalcErrorProbabilityOfCombinedChannels(length int) []struct {
		Index int
		Prob  float64
	}
}

func calcErrorProbabilityViaDensityEvolution(length int, index int, base map[float64]float64) float64 {
	evolvedProbability := densityEvolutionDiscreteProbability(length, index, base)
	sum := 0.0
	for llr, prob := range evolvedProbability {
		switch {
		case llr == 0:
			sum += prob / 2
		case llr < 0:
			sum += prob
		}
	}

	return sum
}

const padding = 1000000

var memo = make(map[int]map[float64]float64)

/*
approximate log likelihood ratio for evaluate density evolution probability.
*/
func approximateLLR(llr float64, prob float64) (float64, bool) {
	if prob < 0.0000001 {
		return 0, false
	}
	v := math.Floor(llr*20) / 20
	if math.IsNaN(v) {
		return 0, false
	}
	return v, true
}

func densityEvolutionDiscreteProbability(length int, index int, base map[float64]float64) map[float64]float64 {
	if length == 1 {
		return base
	}

	if value, ok := memo[length*padding+index]; ok {
		return value
	}

	if index%2 == 0 {
		child := densityEvolutionDiscreteProbability(length/2, index/2, base)
		ret := make(map[float64]float64)
		for llr1, prob1 := range child {
			for llr2, prob2 := range child {
				llr, isTarget := approximateLLR(llr1+llr2, prob1*prob2)
				if !isTarget {
					continue
				}
				ret[llr] += prob1 * prob2
			}
		}
		memo[length*padding+index] = ret
		return ret
	}
	child := densityEvolutionDiscreteProbability(length/2, (index+1)/2, base)
	ret := make(map[float64]float64)
	for llr1, prob1 := range child {
		for llr2, prob2 := range child {
			llr := 2 * math.Atanh(math.Tanh(llr1/2)*math.Tanh(llr2/2))
			llr, isTarget := approximateLLR(llr, prob1*prob2)
			if !isTarget {
				continue
			}
			ret[llr] += prob1 * prob2
		}
	}
	memo[length*padding+index] = ret
	return ret
}
