package main

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"strconv"
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

func reverseBit(bit string) string {
	newBit := ""
	for _, c := range bit {
		if c == '0' {
			newBit = newBit + "1"
		} else if c == '1' {
			newBit = newBit + "0"
		}
	}
	return newBit
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func bitToBytes(bit string) []byte {
	var bytes []byte
	runes := []rune(bit)
	for i := 0; i < len(runes)/8; i++ {
		s := string(runes[i*8 : (i+1)*8])
		num, _ := strconv.ParseUint(s, 2, 8)
		bytes = append(bytes, byte(num))
	}
	return bytes
}

func toHash(txHex string) chainhash.Hash {
	txByte, _ := hex.DecodeString(txHex)
	txBit := toBitString(txByte)
	// fmt.Println(txHex, len(txHex))
	// fmt.Println(txByte, len(txByte))
	// fmt.Println(txBit, len(txBit))

	txBit = reverseBit(txBit)
	txBit = reverse(txBit)
	// fmt.Println(txBit)

	bytes := bitToBytes(txBit)

	return chainhash.DoubleHashH(bytes)
}

func main() {
	txHex1 := "4eb8629ffb3bdf1035951d6df78fdb0bf5770a1b6b5744995ad593a52b8c2dc3"
	txHex2 := "ba219c854e7a05df49d45e9c97bd006f3ef9daefe93a43a82c82cf0eef02569d"
	hash1 := toHash(txHex1)
	hash2 := toHash(txHex2)
	fmt.Println(hash1)
	fmt.Println(hash2)

	hashSum := toHash(txHex1 + txHex2)
	fmt.Println(hashSum)
}
