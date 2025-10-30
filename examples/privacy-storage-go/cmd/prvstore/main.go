package main

import (
	"context"
	// "encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/firefly-signer/pkg/abi"

	"github.com/LF-Decentralized-Trust-labs/paladin/examples/common-go/pkg"
	nototypes "github.com/LF-Decentralized-Trust-labs/paladin/domains/noto/pkg/types"
	"github.com/LF-Decentralized-Trust-labs/paladin/sdk/go/pkg/pldtypes"
	// "github.com/LF-Decentralized-Trust-labs/paladin/sdk/go/pkg/pldapi"
	"github.com/LF-Decentralized-Trust-labs/paladin/sdk/go/pkg/pldclient"
)

type PublicStorageInput struct {
	Num int `json:"num"`
}

type PublicStorageOutput struct {
	Value string `json:"value"`
}

type penteConstructorParams struct {
	Group                nototypes.PentePrivateGroup `json:"group"`
	EVMVersion           string                      `json:"evmVersion"`
	EndorsementType      string                      `json:"endorsementType"`
	ExternalCallsEnabled bool                        `json:"externalCallsEnabled"`
}

type penteDeployParams struct {
	Group    nototypes.PentePrivateGroup `json:"group"`
	Bytecode pldtypes.HexBytes           `json:"bytecode"`
	Inputs   any                         `json:"inputs"`
}

var pentePrivGroupComps = abi.ParameterArray{
	{Name: "salt", Type: "bytes32"},
	{Name: "members", Type: "string[]"},
}

var penteGroupABI = &abi.Parameter{
	Name: "group", Type: "tuple", Components: pentePrivGroupComps,
}

var penteConstructorABI = &abi.Entry{
	Type: abi.Constructor, Inputs: abi.ParameterArray{
		penteGroupABI,
		{Name: "evmVersion", Type: "string"},
		{Name: "endorsementType", Type: "string"},
		{Name: "externalCallsEnabled", Type: "bool"},
	},
}


func main() {
	nodeConnections := examples_commongo.GetNodeConnections()
	if len(nodeConnections) < 3 {
		panic(fmt.Errorf("should be at least 3 node connections"))
	}

	fmt.Println("\n>>> STEP 0 : Initializing Paladin Clients from environment configuration...")
	paladinClientNode1, err := pldclient.New().HTTP(context.Background(), &nodeConnections[0].ClientOptions)
	if err != nil {
		panic(fmt.Errorf(fmt.Sprintf("Unable to instanciate paladin client 1: %s", err)))
	}
	verifierNode1 := "member@node1"
	// paladinClientNode2, err := pldclient.New().HTTP(context.Background(), &nodeConnections[1].ClientOptions)
	// if err != nil {
	// 	panic(fmt.Errorf(fmt.Sprintf("Unable to instanciate paladin client 2: %s", err)))
	// }
	verifierNode2 := "member@node2"
	// paladinClientNode3, err := pldclient.New().HTTP(context.Background(), &nodeConnections[2].ClientOptions)
	// if err != nil {
	// 	panic(fmt.Errorf(fmt.Sprintf("Unable to instanciate paladin client 3: %s", err)))
	// }
	// verifierNode3 := "outsider@node3"

	fmt.Println("\n>>> STEP 1: Deploying pente new privacy group...")
	ENDORSEMENT_TYPE__GROUP_SCOPED_IDENTITIES := "group_scoped_identities"
	penteGroupNodes1and2 := nototypes.PentePrivateGroup{
		Salt:    pldtypes.RandBytes32(),                 // unique salt must be shared privately to retain anonymity
		Members: []string{verifierNode1, verifierNode2}, // these will be salted to establish the endorsement key identifiers
	}	
	txP := paladinClientNode1.TxBuilder(context.Background()).
		Private().Domain("pente").Constructor().
		Inputs(&penteConstructorParams{
			Group:                penteGroupNodes1and2,
			EVMVersion:           "shanghai",
			EndorsementType:      ENDORSEMENT_TYPE__GROUP_SCOPED_IDENTITIES,
			ExternalCallsEnabled: true,
		}).
		ABI(abi.ABI{penteConstructorABI}).From(verifierNode1).Send()
	resP := txP.Wait(1500 * time.Millisecond)
	if resP.Error() != nil {
		panic(fmt.Errorf("Unable to deploy pente privacy group: %s", resP.Error()))
	}
	receiptP := resP.Receipt()
	if receiptP == nil {
		panic(fmt.Errorf("Unable to deploy pente privacy group: no receipt ?"))
	}
	if receiptP.ContractAddress == nil || len(receiptP.ContractAddress) == 0 {
		fmt.Println(fmt.Sprintf("Unable to deploy pente privacy group: no contract address in receipt ? -> %#v", receiptP))
		panic(fmt.Errorf("Unable to deploy pente privacy group: no contract address in receipt ?"))
	}
	penteContract := receiptP.ContractAddress
	fmt.Println(fmt.Sprintf(">>> STEP 1: pente privacy group deployed [%s]...", penteContract))

	fmt.Println("\n>>> STEP 2: Deploying contract storage into pente privacy group...")
	abi := examples_commongo.GetSMCABI()
	byc := examples_commongo.GetSMCBytecode()
	txD := paladinClientNode1.TxBuilder(context.Background()).
		Private().ABI(abi).Domain("pente").To(penteContract).
		Function("").
		Inputs(&penteDeployParams{
			Group:    penteGroupNodes1and2,
			Bytecode: byc,
		}).From(verifierNode1).Send()
	resD := txD.Wait(1500 * time.Millisecond)
	if resD.Error() != nil {
		panic(fmt.Errorf("Unable to deploy contract storage into pente privacy group: %s", resP.Error()))
	}
	receiptD := resP.Receipt()
	if receiptD == nil {
		panic(fmt.Errorf("Unable to deploy contract storage into pente privacy group: no receipt ?"))
	}
	if receiptD.ContractAddress == nil || len(receiptD.ContractAddress) == 0 {
		fmt.Println(fmt.Sprintf("Unable to deploy contract storage into pente privacy group: no contract address in receipt ? -> %#v", receiptP))
		panic(fmt.Errorf("Unable to deploy contract storage into pente privacy group: no contract address in receipt ?"))
	}
	contract := receiptD.ContractAddress
	fmt.Println(fmt.Sprintf(">>> STEP 2: contract storage deployed into pente privacy group [%s]...", contract))
}
