package polar_code

import (
	"github.com/egs33/polar-codes-go/internal/set"
	"github.com/egs33/polar-codes-go/pkg/channel"
	"math"
	"math/rand"
	"sort"
	"time"
)

type PolarCode struct {
	length          int
	dimension       int
	informationBits []int
	/* informationBitSet has same contents of informationBits.
	   This field is for performance.
	*/
	informationBitSet set.IntSet
}

func NewPolarCode(length int, dimension int, channel channel.BinaryMemorylessChannel) PolarCode {
	rand.Seed(time.Now().UnixNano())
	polarCode := PolarCode{
		length:            length,
		dimension:         dimension,
		informationBits:   make([]int, 0),
		informationBitSet: set.NewIntSet(),
	}
	channelErrorProbs := channel.CalcErrorProbabilityOfCombinedChannels(length)

	sort.Slice(channelErrorProbs, func(i, j int) bool {
		return math.Abs(channelErrorProbs[i].Prob) < math.Abs(channelErrorProbs[j].Prob)
	})
	for i := 0; i < dimension; i++ {
		polarCode.informationBitSet.Add(channelErrorProbs[i].Index - 1)
	}
	informationBits := polarCode.informationBitSet.Values()
	sort.Ints(informationBits)
	polarCode.informationBits = informationBits

	return polarCode
}

func (code PolarCode) Encode(information []int) []int {
	data := make([]int, code.length)
	for i, bitIindex := range code.informationBits {
		data[bitIindex] = information[i]
	}
	return encodeRecursive(data)
}

func (code PolarCode) SuccessiveCancellationDecode(received []float64) []int {
	decoded := make([]int, 0)
	informationDecoded := make([]int, 0)
	for i := 0; i < code.length; i++ {
		if !code.informationBitSet.Contains(i) {
			decoded = append(decoded, 0)
			continue
		}
		llr := calcBitLogLikelihoodRatio(received, decoded)
		var bit int
		switch {
		case llr > 0:
			bit = 0
		case llr < 0:
			bit = 1
		case llr == 0:
			bit = rand.Intn(2)
		}
		decoded = append(decoded, bit)
		informationDecoded = append(informationDecoded, bit)
	}
	return informationDecoded
}

type listDecodePath struct {
	pathMetrics float64
	bits        []int
}

func (code PolarCode) SuccessiveCancellationListDecode(received []float64, listSize int) []int {
	decodedList := code.sclDecodeInner(received, listSize)

	decoded := make([]int, 0)
	for index, bit := range decodedList[0].bits {
		if code.informationBitSet.Contains(index) {
			decoded = append(decoded, bit)
		}
	}
	return decoded
}

func (code PolarCode) sclDecodeInner(received []float64, listSize int) []listDecodePath {
	decodedList := []listDecodePath{{pathMetrics: 0, bits: make([]int, 0)}}
	for i := 0; i < code.length; i++ {
		if !code.informationBitSet.Contains(i) {
			for k := range decodedList {
				decodedList[k].bits = append(decodedList[k].bits, 0)
			}
			continue
		}
		nextDecodedList := make([]listDecodePath, 0)

		for _, decoded := range decodedList {
			llr := calcBitLogLikelihoodRatio(received, decoded.bits)
			copy0 := make([]int, len(decoded.bits))
			copy1 := make([]int, len(decoded.bits))
			copy(copy0, decoded.bits)
			copy(copy1, decoded.bits)
			nextDecodedList = append(nextDecodedList,
				listDecodePath{
					pathMetrics: decoded.pathMetrics + math.Log(1+math.Exp(-llr)),
					bits:        append(copy0, 0),
				}, listDecodePath{
					pathMetrics: decoded.pathMetrics + math.Log(1+math.Exp(llr)),
					bits:        append(copy1, 1),
				})
		}
		if len(nextDecodedList) <= listSize {
			decodedList = nextDecodedList
			continue
		}
		sort.Slice(nextDecodedList, func(i, j int) bool {
			return nextDecodedList[i].pathMetrics < nextDecodedList[j].pathMetrics
		})
		decodedList = nextDecodedList[0:listSize]
	}
	return decodedList
}

func (code PolarCode) SuccessiveCancellationNonUniqueListDecode(received []float64, listSize int) [][]int {
	decodedList := code.sclDecodeInner(received, listSize)

	allDecoded := make([][]int, len(decodedList))
	for i, decoded := range decodedList {
		informationBits := make([]int, 0)
		for index, bit := range decoded.bits {
			if code.informationBitSet.Contains(index) {
				informationBits = append(informationBits, bit)
			}
		}
		allDecoded[i] = informationBits
	}

	return allDecoded
}
