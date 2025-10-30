package examples_commongo

import (
	"encoding/json"
	"io"
	"os"

	"github.com/hyperledger/firefly-signer/pkg/abi"
	"github.com/LF-Decentralized-Trust-labs/paladin/sdk/go/pkg/pldtypes"
	"github.com/LF-Decentralized-Trust-labs/paladin/sdk/go/pkg/solutils"
)

type SolidityBuild struct {
	ABI      abi.ABI           `json:"abi"`
	Bytecode pldtypes.HexBytes `json:"bytecode"`
}

func ParseJSONFileToABI(filepath string) (abi.ABI, error) {
	jsonFile, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var abi abi.ABI
	err = json.Unmarshal([]byte(byteValue), &abi)
	if err != nil {
		return nil, err
	}
	return abi, nil
}

func ParseJSONFileToSolidityBuild(filepath string) (*solutils.SolidityBuild, error) {
	jsonFile, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var sbuild solutils.SolidityBuild
	err = json.Unmarshal([]byte(byteValue), &sbuild)
	if err != nil {
		return nil, err
	}
	return &sbuild, nil	
}
