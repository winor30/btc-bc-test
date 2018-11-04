package main

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

func fillZero(bit string) string {
	if len(bit) < 8 {
		bit = "0" + bit
		bit = fillZero(bit)
	}
	return bit
}

func toBitString(bytes []byte) string {
	bits := ""
	for _, b := range bytes {
		bit8 := fmt.Sprintf("%b", b)
		bit8 = fillZero(bit8)
		bits += bit8
	}
	return bits
}

func toHash(tx string) chainhash.Hash {
	return chainhash.DoubleHashH(nil)
}

func main() {
	txStr := "4eb8629ffb3bdf1035951d6df78fdb0bf5770a1b6b5744995ad593a52b8c2dc3"
	fmt.Println(txStr, len(txStr))
	tx, _ := hex.DecodeString(txStr)
	fmt.Println(tx, len(tx))
	txBit := toBitString(tx)
	fmt.Println(txBit, len(txBit))
}
