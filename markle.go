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
	txBit = reverseBit(txBit)
	txBit = reverse(txBit)
	bytes := bitToBytes(txBit)
	return chainhash.DoubleHashH(bytes)
}

func calcMarkleRoot(txs []chainhash.Hash, txTree *[][]chainhash.Hash) []chainhash.Hash {
	if len(txs) != 1 && len(txs)%2 == 1 {
		txs = append(txs, txs[len(txs)-1])
	}
	*txTree = append(*txTree, txs)

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
	return calcMarkleRoot(smtxs, txTree)
}

type HashPath struct {
	ID   int
	Hash chainhash.Hash
}

func (hash HashPath) String() string {

	return strconv.Itoa(hash.ID) + ": " + hash.Hash.String()
}

func getRelatedTxs(txid int, txTree *[][]chainhash.Hash) []HashPath {
	var marklePath []HashPath
	for _, txs := range *txTree {
		if txid%2 == 1 && len(txs) != 1 {
			marklePath = append(marklePath, HashPath{ID: txid - 1, Hash: txs[txid-1]})
		} else if txid%2 == 0 && len(txs) != 1 {
			marklePath = append(marklePath, HashPath{ID: txid + 1, Hash: txs[txid+1]})
		}

		txid = txid / 2
	}
	return marklePath
}

func validation(tx chainhash.Hash, marklePath []HashPath, markleRoot chainhash.Hash) bool {
	expectedRoot := tx
	for _, mhash := range marklePath {
		if mhash.ID%2 == 1 {
			expectedRoot = toHash(expectedRoot.String() + mhash.Hash.String())
		} else {
			expectedRoot = toHash(mhash.Hash.String() + expectedRoot.String())
		}
	}

	return markleRoot.String() == expectedRoot.String()
}

func main() {
	fmt.Println(`"ccc", "ddd", "aaa", "bbb", "eee"`)
	txs := []string{"ccc", "ddd", "aaa", "bbb", "eee"}

	var txHashs []chainhash.Hash
	for _, tx := range txs {
		txHashs = append(txHashs, toHash(tx))
	}
	var txTree [][]chainhash.Hash
	markleRoot := calcMarkleRoot(txHashs, &txTree)[0]
	fmt.Println("Markle Root is ", markleRoot)
	fmt.Println("Markle Path is ...")
	for i, txs := range txTree {
		fmt.Println(i, len(txs), txs)
	}

	// check valid transaction "aaa" that number of tx is 2
	// it will be true
	marklePath := getRelatedTxs(3, &txTree)
	fmt.Println(marklePath)
	fmt.Println("validation is ", validation(txHashs[3], marklePath, markleRoot))

	// check valid transaction "aaa" that number of tx is 3
	// it will be false
	marklePath = getRelatedTxs(1, &txTree)
	fmt.Println(marklePath)
	fmt.Println("validation is ", validation(txHashs[0], marklePath, markleRoot))
}
