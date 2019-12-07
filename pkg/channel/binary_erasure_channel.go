package channel

import (
	"math"
	"math/rand"
	"time"
)

type BinaryErasureChannel struct {
	erasureProbability float64
}

func (bec BinaryErasureChannel) CalcErrorProbabilityOfCombinedChannels(length int) []struct {
	Index int
	Prob  float64
} {
	memo = make(map[int]map[float64]float64)
	channelErrorProbs := make([]struct {
		Index int
		Prob  float64
	}, length)
	for i := 0; i < length; i++ {
		channelErrorProbs[i] = struct {
			Index int
			Prob  float64
		}{Index: i + 1, Prob: bec.CalcErrorProbabilityOfCombinedChannel(length, i+1)}
	}
	return channelErrorProbs
}

func (bec BinaryErasureChannel) CalcErrorProbabilityOfCombinedChannel(length int, index int) float64 {
	return calcErrorProbabilityViaDensityEvolution(length, index,
		map[float64]float64{math.Inf(1): 1 - bec.erasureProbability, 0: bec.erasureProbability})
}

func NewBinaryErasureChannel(erasureProbability float64) BinaryErasureChannel {
	rand.Seed(time.Now().UnixNano())
	return BinaryErasureChannel{erasureProbability: erasureProbability}
}

func (bec BinaryErasureChannel) Channel(input []int) []float64 {
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
