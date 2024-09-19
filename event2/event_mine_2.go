package main

import (
	"casper/contract/helper"
	"fmt"
	"os"
	"regexp"

	"github.com/make-software/casper-go-sdk/v2/sse"
	"gopkg.in/yaml.v3"
)

type response2 struct {
	Partner_id       int    `json:"partner_id"`
	Amount           int    `json:"amount"`
	Dest_chain_id    string `json:"dest_chain_id"`
	Dest_amount      int    `json:"dest_amount"`
	Deposit_id       int    `json:"deposit_id"`
	Src_token        string `json:"src_token"`
	Recipient        string `json:"recipient"`
	Refund_recipient string `json:"refund_recipient"`
	Dest_token       string `json:"dest_token"`
}

func main() {
	helper.Title()

	data, err := os.ReadFile("./data/sse.json")

	if err != nil {
		panic(err)
	}
	rawEvent := sse.RawEvent{
		Data: data,
	}

	res, err := rawEvent.ParseAsTransactionProcessedEvent()
	if err != nil {
		panic(err)
	}

	message := res.TransactionProcessedPayload.Messages
	// fmt.Print(*message[0].Message.String)
	fmt.Println()

	str := *message[0].Message.String
	// fmt.Println(str)

	// regular expression pattern
	r, _ := regexp.Compile("{(.*?)}")

	// find the content between curly brackets
	mine := r.FindString(str)

	res2 := &response2{}
	blob := []byte(mine)

	// convert the yaml string to struct type
	if err := yaml.Unmarshal(blob, &res2); err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Printf("Partner_id: %+v", res2.Partner_id)
		fmt.Println()
		fmt.Printf("Amount: %+v", res2.Amount)
		fmt.Println()
		fmt.Printf("Dest_chain_id: %+v", res2.Dest_chain_id)
		fmt.Println()
		fmt.Printf("Dest_amount: %+v", res2.Dest_amount)
		fmt.Println()
		fmt.Printf("Deposit_id: %+v", res2.Deposit_id)
		fmt.Println()
		fmt.Printf("Recipient: %+v", res2.Recipient)
		fmt.Println()
		fmt.Printf("Src_token: %+v", res2.Src_token)
		fmt.Println()
		fmt.Printf("Refund_recipient: %+v", res2.Refund_recipient)
		fmt.Println()
		fmt.Printf("Dest_token: %+v", res2.Dest_token)
	}

}

// type MessagePayload struct {
// 	String *string      `json:"String"`
// 	Bytes  *ModuleBytes `json:"Bytes"`
// }
// type Message struct {
// 	// The payload of the message.
// 	Message MessagePayload `json:"message"`
// 	// The name of the topic on which the message was emitted on.
// 	TopicName string `json:"topic_name"`
// 	// The hash of the name of the topic.
// 	TopicNameHash key.Hash `json:"topic_name_hash"`
// 	// The identity of the entity that produced the message.
// 	EntityHash key.EntityAddr `json:"entity_hash"`
// 	// Message index in the block.
// 	BlockIndex uint64 `json:"block_index"`
// 	// Message index in the topic.
// 	TopicIndex uint32 `json:"topic_index"`
// }

// TransactionProcessedPayload struct {
// 	BlockHash       key.Hash              `json:"block_hash"`
// 	TransactionHash types.TransactionHash `json:"transaction_hash"`
// 	InitiatorAddr   types.InitiatorAddr   `json:"initiator_addr"`
// 	Timestamp       time.Time             `json:"timestamp"`
// 	TTL             string                `json:"ttl"`
// 	ExecutionResult types.ExecutionResult `json:"execution_result"`
// 	Messages        []types.Message       `json:"messages"`
// }

// FundsDepositedEventLog(FundsDepositedEventLog { partner_id: 1, amount: 1, dest_chain_id: "ethereum", dest_amount: 90000, deposit_id: 20, src_token: "ec301e17c49ee4d18fc2d3f3766fce82389edac756b2f85aef31a8422414289a", recipient: "Key::Account(a4628515772103ba1174bcf28cc99a8785d291c5192e142f53bd23ecfa55556a)", refund_recipient: "Key::Account(0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20)", dest_token: "Key::AddressableEntity(contract-a26eba1eed80d8248d15a619f2395b7b75e79baff4873a0486052dec1fa1b4c1)" })
