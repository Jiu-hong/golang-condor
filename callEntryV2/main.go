package main

import (
	"context"
	"encoding/json"
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

	key1, err := key.NewKey("account-hash-5cc9364fc67616a74bd51122a439415f9b364098ea488ec833e69c4431a997e9")
	if err != nil {
		panic(err)
	}
	args := &types.Args{}
	args.AddArgument("owner", clvalue.NewCLKey(key1)).
		AddArgument("amount", *clvalue.NewCLUInt256(big.NewInt(1_000_000_000_000_000)))

	contractHash, err := key.NewHash("4c1c6de608510bf352dd4be09999f32d893267b5f7e6d4e493913b01402f7017")
	if err != nil {
		panic(err)
	}
	ep := "mint"
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
	jsonResult, err := json.Marshal(transactionRes)
	if err != nil {
		panic(err)
	}
	log.Println("Transaction info:", string(jsonResult))

}
