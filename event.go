package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/make-software/casper-go-sdk/v2/casper"
	"github.com/make-software/ces-go-parser/v2"
)

func main2() {
	contractHash, err := casper.NewHash("5f357dd77968a7408d4f5ccd6f31e72d7fdb98c0ad25bb1d093e1e9af2138d32")
	fmt.Println("contractHash", contractHash)
	if err != nil {
		panic(err)
	}
	handler := casper.NewRPCHandler("http://3.14.48.188:7777/rpc", http.DefaultClient)
	client := casper.NewRPCClient(handler)

	blockHeight := uint64(3453286)
	fmt.Println("blockHeight", blockHeight)
	block, err := client.GetBlockByHeight(context.Background(), blockHeight)
	if err != nil {
		log.Fatalf("Failed to get block: %v", err)
	}

	for _, transaction := range block.Block.Transactions {
		deployInfo, err := client.GetDeploy(context.Background(), transaction.Hash.String())
		if err != nil {
			log.Fatalf("Failed to get the deploy: %v", err)
		}
		fmt.Println("deployInfo.Deploy.Hash", deployInfo.Deploy.Hash)
		// fmt.Println(deployInfo.Deploy.Session.StoredVersionedContractByHash.Hash.String())
		// fmt.Println(deployInfo.Deploy.Session)
		executionResult := deployInfo.ExecutionResults.ExecutionResult

		cesParser, err := ces.NewParser(client, []casper.Hash{contractHash})
		if err != nil {
			panic(err)
		}
		parseResult, err := cesParser.ParseExecutionResults(executionResult)
		if err != nil {
			panic(err)
		}
		for _, result := range parseResult {
			if result.Error == nil {
				eventName := result.Event.Name
				eventData := result.Event.Data
				fmt.Println(eventName, eventData)
			} else {
				fmt.Println(result.Error)
			}
		}
	}
}
