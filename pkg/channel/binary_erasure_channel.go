package channel

import (
	"math"
	"math/rand"
	"time"
)

type BinaryErasureChannel struct {
	erasureProbability float64
}

func (bec BinaryErasureChannel) CalcErrorProbabilityOfCombinedChannel(length int, index int) float64 {
	return calcErrorProbabilityViaDensityEvolution(length, index,
		map[float64]float64{math.Inf(1): 1 - bec.erasureProbability, 0: bec.erasureProbability})
}

func NewBinaryErasureChannel(erasureProbability float64) BinaryErasureChannel {
	return BinaryErasureChannel{erasureProbability: erasureProbability}
}

func (bec BinaryErasureChannel) Channel(input []int) []float64 {
	rand.Seed(time.Now().UnixNano())
	output := make([]float64, len(input))

	for index, bit := range input {
		if rand.Float64() < bec.erasureProbability {
			output[index] = 0
			continue
		}
		if bit == 0 {
			output[index] = math.Inf(1)
		} else {
			output[index] = math.Inf(-1)
		}
	}

	return output
}
