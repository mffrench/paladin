package main

import (
	"context"
	"encoding/json"
	// "encoding/json"
	"fmt"
	"time"

	"github.com/LF-Decentralized-Trust-labs/paladin/examples/common-go/pkg"
	"github.com/LF-Decentralized-Trust-labs/paladin/sdk/go/pkg/pldapi"
	"github.com/LF-Decentralized-Trust-labs/paladin/sdk/go/pkg/pldclient"
)

type PublicStorageInput struct {
	Num int `json:"num"`
}

type PublicStorageOutput struct {
	Value string `json:"value"`
}

func main() {
	nodeConnections := examples_commongo.GetNodeConnections()
	if len(nodeConnections) == 0 {
		panic(fmt.Errorf("no node connections"))
	}

	fmt.Println("\n>>> STEP 0 : Initializing Paladin Client from environment configuration...")
	paladinClientNode1, err := pldclient.New().HTTP(context.Background(), &nodeConnections[0].ClientOptions)
	if err != nil {
		panic(fmt.Errorf(fmt.Sprintf("Unable to instanciate paladin client: %s", err)))
	}

	fmt.Println("\n>>> STEP 1: Deploying the PublicStorage contract...")
	abi := examples_commongo.GetSMCABI()
	byc := examples_commongo.GetSMCBytecode()
	// TODO: check how to get verifier value owner@node1 from golang
	txD := paladinClientNode1.TxBuilder(context.Background()).
		Public().ABI(abi).Bytecode(byc).From("owner@node1").Send()
	resD := txD.Wait(1500 * time.Millisecond)
	if resD.Error() != nil {
		panic(fmt.Errorf("Unable to deploy smart contract hello world: %s", resD.Error()))
	}
	receiptD := resD.Receipt()
	if receiptD == nil {
		panic(fmt.Errorf("Unable to deploy smart contract hello world: no receipt ?"))
	}
	if receiptD.ContractAddress == nil || len(receiptD.ContractAddress) == 0 {
		fmt.Println(fmt.Sprintf("Unable to deploy smart contract hello world: no contract address in receipt ? -> %#v", receiptD))
		panic(fmt.Errorf("Unable to deploy smart contract hello world: no contract address in receipt ?"))
	}
	contractAddress := receiptD.ContractAddress
	fmt.Println(fmt.Sprintf(">>> STEP 1: Deployed the PublicStorage contract [%s]", contractAddress))

	fmt.Println("\n>>> STEP 2: Calling the store function...")
	txH := paladinClientNode1.TxBuilder(context.Background()).
		Public().ABI(abi).From("owner@node1").Function("store").
		To(contractAddress).Inputs(PublicStorageInput{
			Num: 123,
		}).Send()
	resH := txH.Wait(1500 * time.Millisecond)
	if resH.Error() != nil {
		panic(fmt.Errorf("Unable to call public storage SMC func: store. Reason: %s", resH.Error()))
	}
	receiptH := resH.Receipt()
	if receiptH == nil {
		panic(fmt.Errorf("Unable to call public storage SMC func: store. no receipt ?"))
	}
	if receiptH.TransactionHash == nil || len(receiptH.TransactionHash) == 0 {
		fmt.Println(fmt.Sprintf("Unable to call public storage SMC func: store. no transaction hash in receipt ? -> %#v", receiptH))
		panic(fmt.Errorf("Unable to call public storage SMC func: store. no transaction hash in receipt ?"))
	}
	if !receiptH.Success {
		panic("Unable to call public storage SMC func: store. Transaction failed!")
	}
	fmt.Println(fmt.Sprintf(">>> STEP 2: store function executed successfully [%s]", receiptH.TransactionHash))

	fmt.Println("\n>>> STEP 3: Retrieving the stored value...")
	queryBuilder := paladinClientNode1.TxBuilder(context.Background()).
		Public().ABI(abi).From("owner@node1").Function("retrieve").
		To(contractAddress)
	callTX := &pldapi.TransactionCall{
		TransactionInput: *queryBuilder.BuildTX().TX(),
		DataFormat:       "mode=array",
		PublicCallOptions: pldapi.PublicCallOptions{
			Block: "latest",
		},
	}
	rjson, err := paladinClientNode1.PTX().Call(context.Background(), callTX)
	if err != nil {
		panic(fmt.Errorf("Unable to call public storage SMC func: retrieve: %s", err))
	}
	var value PublicStorageOutput
	err = json.Unmarshal(rjson, &value)
	if err != nil {
		panic(fmt.Errorf("Unable to unmarshal public storage SMC func retrieve result: %s", err))
	}
	if value.Value != "123" {
		panic(fmt.Errorf("value is incorrect !"))
	}
	fmt.Println(">>> STEP 3: successfully retrieved stored value")
}
