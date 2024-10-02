package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"casper/contract/utils"

	"github.com/make-software/casper-go-sdk/v2/rpc"
)

func main() {
	TransactionHash := "7a212286fd55573fd404de7ef2a663b982cd31124c4fbfdc4bbd3d205b86f7f4"
	rpcClient := rpc.NewClient(rpc.NewHttpHandler(utils.ENDPOINT, http.DefaultClient))
	transaction, err := rpcClient.GetTransactionByTransactionHash(context.Background(), TransactionHash)
	if err != nil {
		panic(err)
	}

	jsonResult, err := json.Marshal(transaction)
	if err != nil {
		panic(err)
	}
	log.Println("Transaction info:", string(jsonResult))
}
