package util

import (
	"math/rand"
)

func IsSameBitSequence(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func GenerateBitSlice(length int) []int {
	ret := make([]int, length)
	for i := 0; i < length; i++ {
		ret[i] = rand.Intn(2)
	}
	return ret
}
