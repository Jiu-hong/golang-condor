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
	deployHash := "e2e370e502b759c21d68313c0f101a75a06d6780a896f5b1eb1d6152bdbfd733"
	rpcClient := rpc.NewClient(rpc.NewHttpHandler(utils.ENDPOINT, http.DefaultClient))
	deploy, err := rpcClient.GetDeploy(context.Background(), deployHash)
	if err != nil {
		panic(err)
	}
	jsonResult, err := json.Marshal(deploy)
	if err != nil {
		panic(err)
	}
	log.Println("Deploy info:", string(jsonResult))
}
