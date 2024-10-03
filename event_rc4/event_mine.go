package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/make-software/casper-go-sdk/v2/casper"
	"github.com/make-software/ces-go-parser/v2"
)

func main() {
	testnetNodeAddress := "http://3.14.48.188:7777/rpc"
	rpcClient := casper.NewRPCClient(casper.NewRPCHandler(testnetNodeAddress, http.DefaultClient))

	ctx := context.Background()
	transactionResult, err := rpcClient.GetTransactionByTransactionHash(ctx, "f8f41e0af65d775898816ecc103f301fa21e42f15bec3f995214e3d7be71e472")
	if err != nil {
		panic(err)
	}
	contractHash, err := casper.NewHash("84c52e578dffd9bf39949d0d0e38c5eaa09c8299a7b2956a5bbc3c51780598c4")
	if err != nil {
		panic(err)
	}

	// for rc4 only since
	// AddressableEntity functionality is enabled only for rc4,
	// and will not be used in the future versions for some time.
	parser, err := ces.NewParserWithVersion(rpcClient, []casper.Hash{contractHash}, ces.Casper2xRC4)
	if err != nil {
		panic(err)
	}

	parseResults, err := parser.ParseExecutionResults(transactionResult.ExecutionInfo.ExecutionResult)
	if err != nil {
		panic(err)
	}
	for _, result := range parseResults {
		if result.Error != nil {
			panic(result.Error)
		}
		fmt.Println(result.Event)
	}
}
