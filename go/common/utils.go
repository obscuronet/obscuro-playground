package common

import (
	"math/rand"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

func MaxInt(x, y uint32) uint32 {
	if x < y {
		return y
	}
	return x
}

// ExtractPotentialAddress - given a 32 byte hash , it checks whether it can be an address and extracts that
func ExtractPotentialAddress(hash gethcommon.Hash) *gethcommon.Address {
	bitlen := hash.Big().BitLen()
	// Addresses have 20 bytes. If the field has more, it means it is clearly not an address
	// Discovering addresses with more than 20 leading 0s is very unlikely, so we assume that
	// any topic that has less than 80 bits of data to not be an address for sure
	if bitlen < 80 || bitlen > 160 {
		return nil
	}
	a := gethcommon.BytesToAddress(hash.Bytes())
	return &a
}

// Generates a random string n characters long.
func RandomStr(n int) string {
	randGen := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec

	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	suffix := make([]rune, n)
	for i := range suffix {
		suffix[i] = letters[randGen.Intn(len(letters))]
	}
	return string(suffix)
}
