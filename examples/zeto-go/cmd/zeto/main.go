package main

import (
	"context"
	// "encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/hyperledger/firefly-signer/pkg/abi"

	zetotypes "github.com/LF-Decentralized-Trust-labs/paladin/domains/zeto/pkg/types"
	"github.com/LF-Decentralized-Trust-labs/paladin/examples/common-go/pkg"
	"github.com/LF-Decentralized-Trust-labs/paladin/sdk/go/pkg/pldapi"
	"github.com/LF-Decentralized-Trust-labs/paladin/sdk/go/pkg/pldclient"
	"github.com/LF-Decentralized-Trust-labs/paladin/sdk/go/pkg/pldtypes"
	"github.com/LF-Decentralized-Trust-labs/paladin/sdk/go/pkg/query"
)

// const paladinPrefix = "node"
const tokenType = "Zeto_Anon"
var zetoConstructorABI = &abi.Entry{
	Type: abi.Constructor, Inputs: abi.ParameterArray{
		{Name: "tokenName", Type: "string"},
	},
}
var zetoCoinSchemaID *pldtypes.Bytes32
var zetoContractAddress *pldtypes.EthAddress

func with10DecimalsToHexUint256(x int64) *pldtypes.HexUint256 {
	bx := with10DecimalsToBigInt(x)
	return (*pldtypes.HexUint256)(bx)
}

func with10DecimalsToBigInt(x int64) *big.Int {
	return new(big.Int).Mul(
		big.NewInt(x),
		new(big.Int).Exp(big.NewInt(10), big.NewInt(10), big.NewInt(0)),
	)
}

func logWallet(ctx context.Context, identity string, node pldclient.PaladinClient, awaitedBalance *big.Int, step string) {
	fmt.Println(fmt.Sprintf("\n>>> %s: log wallet...", step))
	var addr pldtypes.HexBytes
	err := node.CallRPC(ctx, &addr, "ptx_resolveVerifier", identity, "domain:zeto:snark:babyjubjub", "iden3_pubkey_babyjubjub_compressed_0x")
	if err != nil {
		panic(fmt.Sprintf("ptx_resolveverifier failed: %s", err))
	}
	method := "pstate_queryContractStates"
	// if isNullifier {
	// 	method = "pstate_queryContractNullifiers"
	// }
	var coins []*zetotypes.ZetoCoinState
	err = node.CallRPC(ctx, &coins, method, "zeto", zetoContractAddress, zetoCoinSchemaID,
		query.NewQueryBuilder().Equal("owner", addr).Limit(100).Query(),
		"available")
	if err != nil {
		panic(fmt.Sprintf("query failed: %s", err))
	}

	balance := big.NewInt(0)
	summary := make([]string, len(coins))
	for ic, c := range coins {
		summary[ic] = fmt.Sprintf("%s...[%s]", c.ID.String()[0:8], c.Data.Amount.Int().Text(10))
		balance = new(big.Int).Add(balance, c.Data.Amount.Int())
	}
	fmt.Println(fmt.Sprintf("\n>>> %s: %s@%s balance=%s coins:%v awaitedBalance: %s", step, identity, node, balance, summary, awaitedBalance))
	if balance.Cmp(awaitedBalance) != 0 {
		panic("balance mismatch !")
	}
}


func main() {
	ctx := context.Background()
	
	nodeConnections := examples_commongo.GetNodeConnections()
	if len(nodeConnections) < 3 {
		panic(fmt.Errorf("should be at least 3 node connections"))
	}

	fmt.Println("\n> Privacy-preserving CBDC token, using private minting...")
	
	fmt.Println("\n>>> STEP 0 : Initializing Paladin Clients from environment configuration...")
	paladinClientNode1, err := pldclient.New().HTTP(context.Background(), &nodeConnections[0].ClientOptions)
	if err != nil {
		panic(fmt.Errorf(fmt.Sprintf("Unable to instanciate paladin client 1: %s", err)))
	}
	bank1 := "bank1@node1"
	paladinClientNode2, err := pldclient.New().HTTP(context.Background(), &nodeConnections[1].ClientOptions)
	if err != nil {
		panic(fmt.Errorf(fmt.Sprintf("Unable to instanciate paladin client 2: %s", err)))
	}
	bank2 := "bank2@node2"
	paladinClientNode3, err := pldclient.New().HTTP(context.Background(), &nodeConnections[2].ClientOptions)
	if err != nil {
		panic(fmt.Errorf(fmt.Sprintf("Unable to instanciate paladin client 3: %s", err)))
	}
	cbdcIssuer := "centralbank@node3"

	fmt.Println("\n>>> STEP 1 : Deploy Palatin Zeto contract...")
	deploy := paladinClientNode3.ForABI(ctx, abi.ABI{zetoConstructorABI}).				
			Private().
			Domain("zeto").
			Constructor().
			From(cbdcIssuer).
			Inputs(&zetotypes.InitializerParams{
				TokenName: tokenType,
			}).
			Send().
			Wait(5 * time.Second)
	if deploy.Error() != nil {
		panic(fmt.Sprintf("deploy failed: %s", err))
	}
	if deploy.Receipt() == nil {
		panic("not address for deployed smart contract ?!")
	}
	zetoContractAddress = deploy.Receipt().ContractAddress
	fmt.Println(fmt.Sprintf("\n>>> Zeto deployed at %s", zetoContractAddress))

	var schemas []*pldapi.Schema
	fmt.Println("\n>>> STEP 2 : get Zeto contract schema ID...")
	err = paladinClientNode3.CallRPC(ctx, &schemas, "pstate_listSchemas", "zeto")
	if err != nil {
		panic(fmt.Sprintf("pstate_listSchemas failed: %s", err))
	}
	for _, s := range schemas {
		if s.Signature == "type=ZetoCoin(uint256 salt,bytes32 owner,uint256 amount,bool locked),labels=[owner,locked]" {
			zetoCoinSchemaID = &s.ID
		}
	}
	if zetoCoinSchemaID == nil {
		panic("no zeto coing schema ID for deployed smart contract ?!")
	}

	fmt.Println("\n>>> STEP 3 : issue tokens to bank1 and bank2...")
	paladinClientNode3.ForABI(ctx, zetotypes.ZetoFungibleABI).
		Private().
		Domain("zeto").
		Function("mint").
		To(zetoContractAddress).
		From(cbdcIssuer).
		Inputs(&zetotypes.FungibleMintParams{
			Mints: []*zetotypes.FungibleTransferParamEntry{
				{
					To:     bank1,
					Amount: with10DecimalsToHexUint256(100000),
				},
				{
					To:     bank2,
					Amount: with10DecimalsToHexUint256(200000),
				},
			},
		}).
		Send().
		Wait(5 * time.Second)		

	logWallet(ctx, bank1, paladinClientNode1, with10DecimalsToBigInt(100000), "STEP 3")
	logWallet(ctx, bank2, paladinClientNode2, with10DecimalsToBigInt(200000), "STEP 3")

	fmt.Println("\n>>> STEP 4 : transfer tokens from bank1 to bank2...")
	paladinClientNode1.ForABI(ctx, zetotypes.ZetoFungibleABI).
		Private().
		Domain("zeto").
		Function("transfer").
		To(zetoContractAddress).
		From(bank1).
		Inputs(&zetotypes.FungibleTransferParams{
			Transfers: []*zetotypes.FungibleTransferParamEntry{
				{
					To:     bank2,
					Amount: with10DecimalsToHexUint256(1000),
				},
			},
		}).
		Send().
		Wait(5 * time.Second)


	logWallet(ctx, bank1, paladinClientNode1, with10DecimalsToBigInt(99000), "STEP 4")
	logWallet(ctx, bank2, paladinClientNode2, with10DecimalsToBigInt(201000), "STEP 4")

	// TODO: redeem
}
