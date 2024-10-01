package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"

	"casper/contract/utils"

	"github.com/make-software/casper-go-sdk/v2/casper"
	"github.com/make-software/casper-go-sdk/v2/rpc"
	"github.com/make-software/casper-go-sdk/v2/types"
	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
	"github.com/make-software/casper-go-sdk/v2/types/key"
)

func main() {
	keys, err := casper.NewED25519PrivateKeyFromPEMFile(utils.KEYPATH)
	if err != nil {
		panic(err)
	}
	pubKey := keys.PublicKey()

	header := types.TransactionV1Header{
		ChainName: utils.NETWORKNAME,
		Timestamp: types.Timestamp(time.Now().UTC()),
		TTL:       180000000000,
		InitiatorAddr: types.InitiatorAddr{
			PublicKey: &pubKey,
		},
		PricingMode: types.PricingMode{
			Fixed: &types.FixedMode{
				GasPriceTolerance: 3,
			},
		},
	}

	args := &types.Args{}
	args.AddArgument("name", *clvalue.NewCLString("Test")).
		AddArgument("symbol", *clvalue.NewCLString("test")).
		AddArgument("decimals", *clvalue.NewCLUint8(9)).
		AddArgument("total_supply", *clvalue.NewCLUInt256(big.NewInt(1_000_000_000_000_000))).
		AddArgument("events_mode", *clvalue.NewCLUint8(2)).
		AddArgument("enable_mint_burn", *clvalue.NewCLUint8(1))
	contractHash, err := key.NewHash("4c1c6de608510bf352dd4be09999f32d893267b5f7e6d4e493913b01402f7017")
	if err != nil {
		panic(err)
	}
	ep := "transfer"
	body := types.TransactionV1Body{
		Args: args,
		Target: types.TransactionTarget{
			Stored: &types.StoredTarget{
				ID:      types.TransactionInvocationTarget{ByHash: &contractHash},
				Runtime: types.TransactionRuntimeVmCasperV1,
			},
		},
		TransactionEntryPoint: types.TransactionEntryPoint{
			Custom: &ep,
		},
		TransactionScheduling: types.TransactionScheduling{
			Standard: &struct{}{},
		},
		TransactionCategory: 2,
	}

	transaction, err := types.MakeTransactionV1(header, body)
	if err != nil {
		panic(err)
	}
	err = transaction.Sign(keys)
	if err != nil {
		panic(err)
	}

	b, err := json.Marshal(transaction)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	fmt.Println(string(b))
	// log.Println("transaction", transaction)

	rpcClient := rpc.NewClient(rpc.NewHttpHandler(utils.ENDPOINT, http.DefaultClient))
	res, err := rpcClient.PutTransactionV1(context.Background(), *transaction)
	if err != nil {
		panic(err)
	}

	log.Println("TransactionV1 submitted:", res.TransactionHash.TransactionV1)

	time.Sleep(time.Second * 10)
	transactionRes, err := rpcClient.GetTransactionByTransactionHash(context.Background(), res.TransactionHash.TransactionV1.ToHex())
	if err != nil {
		panic(err)
	}
	// log.Println("transactionRes: ", transactionRes.)

	log.Println("transactionRes.Transaction: ", transactionRes.Transaction)
	log.Println("transactionRes.ExecutionInfo: ", transactionRes.ExecutionInfo)

}
