package main

import (
	"casper/contract/utils"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/make-software/casper-go-sdk/v2/casper"
)

func main() {
	client := casper.NewRPCClient(casper.NewRPCHandler(utils.ENDPOINT, http.DefaultClient))
	blockHash := "1254037f76bb72c3e44585231341c111a0d2bcc69e082928579b8839bfb571c2"
	accountKey := "account-hash-7d18ff4673261bc5295fdc6640de169a5491ff9fe2d776f177d46a4dae1c3f9e"
	res, err := client.QueryGlobalStateByBlockHash(context.Background(), blockHash, accountKey, nil)
	if err != nil {
		panic(err)
	}

	globalStateRes, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	log.Println("globalStateRes info:", string(globalStateRes))
}
