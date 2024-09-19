package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/make-software/casper-go-sdk/v2/casper"
	"github.com/make-software/ces-go-parser/v2"
)

func main3() {
	testnetNodeAddress := "http://3.14.48.188:7777/rpc"
	rpcClient := casper.NewRPCClient(casper.NewRPCHandler(testnetNodeAddress, http.DefaultClient))

	ctx := context.Background()
	deployResult, err := rpcClient.GetDeploy(ctx, "83c5e2e10133946997a7491ceaf94ff269005a7e9c8ccb069aa9ce2b8b686ae1")
	if err != nil {
		panic(err)
	}

	contractHash, err := casper.NewHash("d95b1ccaf999c94e80b8a59ea8607c62eceb7783bf558661702e3fc7ea43dfee")
	if err != nil {
		panic(err)
	}

	parser, err := ces.NewParser(rpcClient, []casper.Hash{contractHash})
	if err != nil {
		fmt.Println("there")
		panic(err)
	}

	parseResults, err := parser.ParseExecutionResults(deployResult.ExecutionResults.ExecutionResult)
	if err != nil {
		fmt.Println("here")
		panic(err)
	}
	for _, result := range parseResults {
		if result.Error != nil {
			panic(err)
		}
		fmt.Println(result.Event)
	}
}
