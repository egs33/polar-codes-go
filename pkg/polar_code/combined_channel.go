package polar_code

import (
	"math"
)

func encodeRecursive(information []int) []int {
	if len(information) == 1 {
		return information
	}
	halfLength := len(information) / 2
	temp1, temp2 := make([]int, halfLength), make([]int, halfLength)
	for i := 0; i < halfLength; i++ {
		temp1[i] = information[2*i] ^ information[2*i+1]
		temp2[i] = information[2*i+1]
	}
	ret := make([]int, len(information))
	encoded1, encoded2 := encodeRecursive(temp1), encodeRecursive(temp2)
	for i := 0; i < halfLength; i++ {
		ret[i] = encoded1[i]
		ret[halfLength+i] = encoded2[i]
	}
	return ret
}

func C(received []float64, decoded []int) float64 {
	return calcBitLogLikelihoodRatio(received, decoded)
}

func calcBitLogLikelihoodRatio(received []float64, decoded []int) float64 {
	if len(received) == 1 {
		return received[0]
	}
	splittedReceived := [][]float64{
		make([]float64, len(received)/2),
		make([]float64, len(received)/2),
	}
	for i := 0; i < len(received)/2; i++ {
		splittedReceived[0][i] = received[i]
		splittedReceived[1][i] = received[len(received)/2+i]
	}
	splittedDecoded := [][]int{make([]int, len(decoded)/2), make([]int, len(decoded)/2)}
	for i := 0; i < len(decoded)/2; i++ {
		splittedDecoded[0][i] = decoded[2*i] ^ decoded[2*i+1]
		splittedDecoded[1][i] = decoded[2*i+1]
	}

	llr1 := calcBitLogLikelihoodRatio(splittedReceived[0], splittedDecoded[0])
	llr2 := calcBitLogLikelihoodRatio(splittedReceived[1], splittedDecoded[1])

	if len(decoded)%2 == 0 {
		return 2 * math.Atanh(math.Tanh(llr1/2)*math.Tanh(llr2/2))
	}
	if decoded[len(decoded)-1] == 0 {
		return llr1 + llr2
	}
	return -llr1 + llr2
}
