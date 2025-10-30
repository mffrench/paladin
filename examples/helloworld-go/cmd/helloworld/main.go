package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/LF-Decentralized-Trust-labs/paladin/examples/common-go/pkg"
	"github.com/LF-Decentralized-Trust-labs/paladin/sdk/go/pkg/pldapi"
	"github.com/LF-Decentralized-Trust-labs/paladin/sdk/go/pkg/pldclient"
)

type HelloWorldInput struct {
	Name string `json:"name"`
}

func main() {
	nodeConnections := examples_commongo.GetNodeConnections()
	if len(nodeConnections) == 0 {
		panic(fmt.Errorf("no node connections"))
	}

	fmt.Println("\n>>> STEP 0 : Initializing Paladin Client from environment configuration...")
	paladinClientNode1, err := pldclient.New().HTTP(context.Background(), &nodeConnections[0].ClientOptions)
	if err != nil {
		panic(fmt.Errorf(fmt.Sprintf("Unable to instanciate paladinc client: %s", err)))
	}

	fmt.Println("\n>>> STEP 1: Deploying the HelloWorld contract...")
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
	fmt.Println(fmt.Sprintf(">>> STEP 1: Deployed the HelloWorld contract [%s]", contractAddress))

	fmt.Println("\n>>> STEP 2: Calling the sayHello function...")
	txH := paladinClientNode1.TxBuilder(context.Background()).
		Public().ABI(abi).From("owner@node1").Function("sayHello").
		To(contractAddress).Inputs(HelloWorldInput{
		Name: "John",
	}).Send()
	resH := txH.Wait(1500 * time.Millisecond)
	if resH.Error() != nil {
		panic(fmt.Errorf("Unable to call hello world SMC func: sayHello. Reason: %s", resH.Error()))
	}
	receiptH := resH.Receipt()
	if receiptH == nil {
		panic(fmt.Errorf("Unable to call hello world func: sayHello. no receipt ?"))
	}
	if receiptH.TransactionHash == nil || len(receiptH.TransactionHash) == 0 {
		fmt.Println(fmt.Sprintf("Unable to call hello world func: sayHello. no transaction hash in receipt ? -> %#v", receiptH))
		panic(fmt.Errorf("Unable to call hello world func: sayHello. no transaction hash in receipt ?"))
	}
	if !receiptH.Success {
		panic("Unable to call hello world func: sayHello. Trasnsaction failed!")
	}
	fmt.Println(fmt.Sprintf(">>> STEP 2: sayHello function executed successfully [%s]", receiptH.TransactionHash))

	fmt.Println("\n>>> STEP 3: Retrieving and verifying emitted events...")
	events, err := paladinClientNode1.BlockIndex().DecodeTransactionEvents(context.Background(), *receiptH.TransactionHash, abi, "pretty=true")
	if err != nil {
		panic(fmt.Errorf("Unable to get event from transaction %s: %s", receiptH.TransactionHash, err))
	}
	if len(events) < 1 {
		panic(fmt.Errorf("Unable to get event from transaction %s: empty events array", receiptH.TransactionHash))
	}
	var eventData pldapi.TransactionActivityRecord
	err = json.Unmarshal(events[0].Data, &eventData)
	if err != nil {
		panic(fmt.Sprintf("Unable to decode event: %s", err))
	}
	if eventData.Message != "Welcome to Paladin, John" {
		panic(fmt.Sprintf("Event data does not match the expected output!"))
	}
	fmt.Println(fmt.Sprintf(">>> STEP 3: received event message: %s", eventData.Message))
}
