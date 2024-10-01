package main

import (
	"casper/contract/utils"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"

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

	header := types.DefaultDeployHeader()
	header.ChainName = utils.NETWORKNAME
	header.Account = keys.PublicKey()

	header.Timestamp = types.Timestamp(time.Now())
	payment := types.StandardPayment(big.NewInt(4000000000))

	sessionArgs := &types.Args{}
	key1, err := key.NewKey("account-hash-bf06bdb1616050cea5862333d1f4787718f1011c95574ba92378419eefeeee59")
	if err != nil {
		panic(err)
	}
	sessionArgs.AddArgument("amount", *clvalue.NewCLUInt256(big.NewInt(2500000000))).
		AddArgument("owner", clvalue.NewCLKey(key1))

	contractHash, err := key.NewContract("8c25484987e1cdbf2bcf73aa438bbc873d11fed0c9a097651a19bbec504e660a")
	if err != nil {
		panic(err)
	}
	varVal := json.Number("1")
	session := types.ExecutableDeployItem{
		StoredVersionedContractByHash: &types.StoredVersionedContractByHash{
			Hash:       contractHash,
			EntryPoint: "mint",
			Version:    &varVal,
			Args:       sessionArgs,
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
