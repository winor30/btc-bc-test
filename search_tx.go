package main

import (
	"bytes"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func toInputTxs(vins []btcjson.Vin) []string {
	var txids []string
	for _, vin := range vins {
		if vin.Txid == "" {
			continue
		}
		txids = append(txids, vin.Txid)
	}
	return txids
}

func getRawTx(txid string, client *rpcclient.Client) *btcjson.TxRawResult {
	txHash, err := chainhash.NewHashFromStr(txid)
	if err != nil {
		log.Fatalln("txid is invalid")
		os.Exit(1)
	}

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
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return decodedTxRaw
}

func searchPrevTxs(txid string, client *rpcclient.Client) []string {
	decodedTxRaw := getRawTx(txid, client)
	if decodedTxRaw == nil {
		return []string{}
	}
	return toInputTxs(decodedTxRaw.Vin)
}

func isCoinBaseTxs(cbtxids []string, client *rpcclient.Client) bool {
	var iscbtxs = true
	for _, cbtxid := range cbtxids {
		cbRawTx := getRawTx(cbtxid, client)
		iscbtxs = iscbtxs && cbRawTx.Vin[0].IsCoinBase()
	}
	return iscbtxs
}

func searchCoinBaseTxs(txid string, client *rpcclient.Client, depth *uint) []string {
	prevtxids := searchPrevTxs(txid, client)

	if len(prevtxids) == 0 {
		return []string{txid}
	}

	log.Println("depth is ", *depth)

	*depth++
	var allcbtxids []string
	for _, prevtxid := range prevtxids {
		cbtxids := searchCoinBaseTxs(prevtxid, client, depth)
		allcbtxids = append(allcbtxids, cbtxids...)
	}
	return allcbtxids
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
	if err != nil {
		log.Fatal(err)
	}
	defer client.Shutdown()

	if len(os.Args) < 2 {
		log.Fatalln("require txid")
		return
	}
	txid := os.Args[1]

	// search coin base transactions
	var depth uint
	cbtxids := searchCoinBaseTxs(txid, client, &depth)
	log.Println("Deepest depth for coin base txs: ", depth)

	// check coinbase
	log.Println("Are they coin base txs ? > ", isCoinBaseTxs(cbtxids, client))
	log.Println("Number of coin base txs > ", len(cbtxids))
}
