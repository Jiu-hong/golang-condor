package main

import (
	"casper/contract/utils"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"time"

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

	header := types.DefaultDeployHeader()
	header.ChainName = utils.NETWORKNAME
	header.Account = keys.PublicKey()

	header.Timestamp = types.Timestamp(time.Now())
	payment := types.StandardPayment(big.NewInt(4000000000))

	moduleBytes, err := os.ReadFile(utils.ContractPath)
	if err != nil {
		panic(err)
	}

	args := &casper.Args{}
	args.AddArgument("name", *clvalue.NewCLString("Test")).
		AddArgument("symbol", *clvalue.NewCLString("test")).
		AddArgument("decimals", *clvalue.NewCLUint8(9)).
		AddArgument("total_supply", *clvalue.NewCLUInt256(big.NewInt(1_000_000_000_000_000))).
		AddArgument("events_mode", *clvalue.NewCLUint8(2)).
		AddArgument("enable_mint_burn", *clvalue.NewCLUint8(1))
	session := types.ExecutableDeployItem{
		ModuleBytes: &types.ModuleBytes{
			ModuleBytes: hex.EncodeToString(moduleBytes),
			Args:        args,
		},
	}

	deploy, err := types.MakeDeploy(header, payment, session)
	if err != nil {
		panic(err)
	}
	err = deploy.Sign(keys)
	if err != nil {
		panic(err)
	}

	rpcClient := rpc.NewClient(rpc.NewHttpHandler(utils.ENDPOINT, http.DefaultClient))
	_, err = rpcClient.PutDeploy(context.Background(), *deploy)
	if err != nil {
		panic(err)
	}
	fmt.Println("deploy hash", deploy.Hash.ToHex())
}
