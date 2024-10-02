package main

import (
	"context"
	"encoding/hex"
	"log"
	"math/big"
	"net/http"
	"time"

	"casper/contract/utils"

	"github.com/make-software/casper-go-sdk/v2/casper"
	"github.com/make-software/casper-go-sdk/v2/rpc"
	"github.com/make-software/casper-go-sdk/v2/types"
	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
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
		TTL:       utils.TTL,
		InitiatorAddr: types.InitiatorAddr{
			PublicKey: &pubKey,
		},
		PricingMode: types.PricingMode{
			Fixed: &types.FixedMode{
				GasPriceTolerance: 3,
			},
		},
	}

	contractPath := "/home/ubuntu/caspereco/cep18/v2/cep18.wasm"
	moduleBytes := utils.GetmoduleBytes(contractPath)

	args := &types.Args{}
	args.AddArgument("name", *clvalue.NewCLString("Test")).
		AddArgument("symbol", *clvalue.NewCLString("test")).
		AddArgument("decimals", *clvalue.NewCLUint8(9)).
		AddArgument("total_supply", *clvalue.NewCLUInt256(big.NewInt(1_000_000_000_000_000))).
		AddArgument("events_mode", *clvalue.NewCLUint8(2)).
		AddArgument("enable_mint_burn", *clvalue.NewCLUint8(1))

	body := types.TransactionV1Body{
		Args: args,
		Target: types.TransactionTarget{
			Session: &types.SessionTarget{
				ModuleBytes: hex.EncodeToString(moduleBytes),
				Runtime:     types.TransactionRuntimeVmCasperV1,
			},
		},
		TransactionEntryPoint: types.TransactionEntryPoint{
			Call: &struct{}{},
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
}
