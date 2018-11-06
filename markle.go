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

func toSumHash(txs []chainhash.Hash) []chainhash.Hash {
	fmt.Println(txs)
	if len(txs) == 1 {
		return txs
	}

	var smtxs []chainhash.Hash
	for i := 0; i < len(txs)/2+1; i++ {
		if 1+2*i < len(txs) {
			smtxs = append(smtxs, toHash(txs[2*i].String()+txs[1+2*i].String()))
		} else if 1+2*i == len(txs) {
			smtxs = append(smtxs, toHash(txs[2*i].String()+txs[2*i].String()))
		}
	}
	return toSumHash(smtxs)
}

func main() {
	fmt.Println(`"aaa", "bbb", "ccc", "ddd", "eee", "fff", "ggg", "hhh"`)
	txs := []string{"aaa", "bbb", "ccc", "ddd", "eee", "fff", "ggg", "hhh"}
	var txHashs []chainhash.Hash
	for _, tx := range txs {
		txHashs = append(txHashs, toHash(tx))
	}
	fmt.Println(toSumHash(txHashs))

	fmt.Println(`"ccc", "ddd", "aaa", "bbb", "eee"`)
	txs = []string{"ccc", "ddd", "aaa", "bbb", "eee"}
	txHashs = []chainhash.Hash{}
	for _, tx := range txs {
		txHashs = append(txHashs, toHash(tx))
	}

	fmt.Println(toSumHash(txHashs))
}
