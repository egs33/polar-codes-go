package main

import (
	"github.com/egs33/polar-codes-go/internal/util"
	"github.com/egs33/polar-codes-go/pkg/channel"
	"github.com/egs33/polar-codes-go/pkg/polar_code"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	bsc := channel.NewBinarySymmetricChannel(0.11)
	code := polar_code.NewPolarCode(64, 32, bsc)

	information := util.GenerateBitSlice(32)
	codeword := code.Encode(information)
	received := bsc.Channel(codeword)
	code.SuccessiveCancellationDecode(received)
}
