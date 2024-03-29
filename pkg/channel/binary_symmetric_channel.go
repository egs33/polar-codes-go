package channel

import (
	"math"
	"math/rand"
	"time"
)

type BinarySymmetricChannel struct {
	crossoverProbability float64
}

func (bsc BinarySymmetricChannel) CalcErrorProbabilityOfCombinedChannels(length int) []struct {
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
		}{Index: i + 1, Prob: bsc.CalcErrorProbabilityOfCombinedChannel(length, i+1)}
	}
	return channelErrorProbs
}

func (bsc BinarySymmetricChannel) CalcErrorProbabilityOfCombinedChannel(length int, index int) float64 {
	llr0 := math.Log2((1 - bsc.crossoverProbability) / bsc.crossoverProbability)
	return calcErrorProbabilityViaDensityEvolution(length, index,
		map[float64]float64{llr0: 1 - bsc.crossoverProbability, -llr0: bsc.crossoverProbability})
}

func NewBinarySymmetricChannel(crossoverProbability float64) BinarySymmetricChannel {
	rand.Seed(time.Now().UnixNano())
	return BinarySymmetricChannel{crossoverProbability: crossoverProbability}
}

func (bsc BinarySymmetricChannel) Channel(input []int) []float64 {
	output := make([]float64, len(input))
	zeroLLR := math.Log((1 - bsc.crossoverProbability) / bsc.crossoverProbability)

	for index, bit := range input {
		if rand.Float64() < bsc.crossoverProbability {
			bit = 1 - bit
		}
		if bit == 0 {
			output[index] = zeroLLR
		} else {
			output[index] = -zeroLLR
		}
	}

	return output
}
