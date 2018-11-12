package main

import (
	"bytes"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func checkError(err error) {
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}

func checkErrorMsg(err error, msg string) {
	if err != nil {
		log.Fatalln(msg)
		os.Exit(1)
	}
}

func getTxHashes(height int64, client *rpcclient.Client) []chainhash.Hash {
	bcHash, err := client.GetBlockHash(height)
	checkError(err)

	bc, err := client.GetBlock(bcHash)
	checkError(err)

	txhashes, err := bc.TxHashes()
	checkError(err)
	return txhashes
}

func getRawTx(txid string, client *rpcclient.Client) *btcjson.TxRawResult {
	txHash, err := chainhash.NewHashFromStr(txid)
	checkErrorMsg(err, "txid is invalid")

	rawTx, err := client.GetRawTransaction(txHash)
	if err != nil {
		// if txid cannot find previous txid, it might be coin base tx
		if err.Error() == "-5: No such mempool or blockchain transaction. Use gettransaction for wallet transactions." {
			return nil
		}
		log.Fatal(err.Error())
		os.Exit(1)
	}

	bufWriter := new(bytes.Buffer)
	if err = rawTx.MsgTx().Serialize(bufWriter); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	decodedTxRaw, err := client.DecodeRawTransaction(bufWriter.Bytes())
	checkError(err)

	return decodedTxRaw
}

func getRawTxs(txhashes []chainhash.Hash, client *rpcclient.Client) []btcjson.TxRawResult {
	var txraws []btcjson.TxRawResult
	for _, txhash := range txhashes {
		txraws = append(txraws, *getRawTx(txhash.String(), client))
	}
	return txraws
}

func calcCirculation(height uint64, client *rpcclient.Client) float64 {
	txhashes := getTxHashes(int64(height), client)

	txs := getRawTxs(txhashes, client)

	var sum float64
	for _, tx := range txs {
		for _, vout := range tx.Vout {
			sum += vout.Value
		}
	}
	return sum
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
		return
	}

	host := os.Getenv("host")
	user := os.Getenv("user")
	pass := os.Getenv("pass")

	// Connect to local bitcoin core RPC server using HTTP POST mode.
	connCfg := &rpcclient.ConnConfig{
		Host:         host,
		User:         user,
		Pass:         pass,
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg, nil)
	checkError(err)
	defer client.Shutdown()

	if len(os.Args) < 3 {
		log.Fatalln("require start and end block height")
		return
	}

	start, err := strconv.Atoi(os.Args[1])
	checkError(err)
	end, err := strconv.Atoi(os.Args[2])
	checkError(err)

	var velocity float64
	for i := start; i < end; i++ {
		velocity += calcCirculation(uint64(i), client)
	}

	velocity = velocity / float64(end-start)

	log.Println("total bitcoin circulation is ", velocity)
}
