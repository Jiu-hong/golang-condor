package main

import (
	"casper/contract/utils"
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/make-software/casper-go-sdk/v2/casper"
)

func main() {
	client := casper.NewRPCClient(casper.NewRPCHandler(utils.ENDPOINT, http.DefaultClient))

	// latest state-root-hash
	latestHashResult, err := client.GetStateRootHashLatest(context.Background())
	if err != nil {
		panic(err)
	}
	stateRootHash := latestHashResult.StateRootHash.String()

	// // public key to be checked CEP18 balance
	// pubKey, err := casper.NewPublicKey("01e05e824a8d8bb1233d7e375ae4b9aefec8c7436eb175af5ae8fb6eefe0660dda")
	// if err != nil {
	// 	panic(err)
	// }
	// accountHash := pubKey.AccountHash()

	// account hash to be checked CEP18 balance
	accountHash, err := casper.NewAccountHash("account-hash-5cc9364fc67616a74bd51122a439415f9b364098ea488ec833e69c4431a997e9")
	if err != nil {
		panic(err)
	}
	// === create dictionary_item_key ===
	prefixByte := []byte{0} // 0 for account ; 1 for contract package
	item_key_bytes := append(prefixByte, accountHash.Bytes()...)
	item_key := b64.StdEncoding.EncodeToString(item_key_bytes)

	// balance uref from the contract named_key
	balances_uref := "uref-fb6d7dd568bb45bd7433498c37fabf0883f8e5700c08a6541530d3425f66f17f-007"

	// get balance
	res, err := client.GetDictionaryItem(context.Background(), &stateRootHash, balances_uref, item_key)
	if err != nil {
		panic(err)
	}

	globalStateRes, err := json.Marshal(res.StoredValue.CLValue)
	if err != nil {
		panic(err)
	}
	fmt.Println("balance info:", string(globalStateRes))
}
