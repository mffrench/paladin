package examples_commongo

import (
	"fmt"
	"os"

	"github.com/hyperledger/firefly-signer/pkg/abi"

	"github.com/LF-Decentralized-Trust-labs/paladin/config/pkg/pldconf"
)

type NodeClientJSON struct {
	URL         string `json:"url"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Bearer      string `json:"bearer"`
	TlsInsecure bool   `json:"tlsInsecure"`
}

type NodeConnectionJSON struct {
	Name   string         `json:"name"`
	ID     string         `json:"id"`
	Client NodeClientJSON `json:"client"`
}

type NodeConnection struct {
	Name          string
	ID            string
	ClientOptions pldconf.HTTPClientConfig
}

func convertJSONToNodeConnection(jsonNode NodeConnectionJSON) *NodeConnection {
	clientOpts := pldconf.HTTPClientConfig{
		URL: jsonNode.Client.URL,
	}

	if len(jsonNode.Client.Username) > 0 && len(jsonNode.Client.Password) > 0 {
		clientOpts.Auth.Username = jsonNode.Client.Username
		clientOpts.Auth.Password = jsonNode.Client.Password
	} /* else if jsonNode.Client.Bearer {
		no bearer in golang sdk yet ?
	}*/

	if jsonNode.Client.TlsInsecure {
		clientOpts.TLS.InsecureSkipHostVerify = true
	}

	return &NodeConnection{
		Name:          jsonNode.Name,
		ID:            jsonNode.ID,
		ClientOptions: clientOpts,
	}
}

func loadFromConfigFile(configPath string) []*NodeConnection {
	config := []*NodeConnectionJSON{}
	ParseJSONFileToStruct(configPath, &config)

	connections := []*NodeConnection{}
	for _, c := range config {
		// fmt.Println(fmt.Sprintf("connection config: %#v", c))
		connections = append(connections, convertJSONToNodeConnection(*c))
	}
	return connections
}

func getConfigPathFromArgs() (string, error) {
	args := os.Args
	if len(args) < 2 {
		return "", fmt.Errorf("no config path in args")
	}

	configPath := args[1]
	return configPath, nil
}

func getABIPathFromArgs() (string, error) {
	args := os.Args
	if len(args) < 3 {
		return "", fmt.Errorf("no smart contract ABI path in args")
	}

	abiPath := args[2]
	return abiPath, nil
}

// func GetCachePath() string {
// 	return ""
// }

// func FindLatestContractDataFile(dataDir string) (string, error) {
// 	return "", nil
// }

var nodeConnections []*NodeConnection

func GetNodeConnections(opts ...OptionsBuilder) []*NodeConnection {
	if nodeConnections != nil {
		return nodeConnections
	}

	options := BuildOptions(opts...)
	var pathToUse string
	var err error
	if _, ok := options.KV[ConfigPathOpt]; ok {
		pathToUse = options.KV[ConfigPathOpt].(string)
	}
	pathToUse, err = getConfigPathFromArgs()
	if err != nil {
		fmt.Println(fmt.Sprintf("Error while defining config path: %s", err))
		panic(err)
	}

	nodeConnections = loadFromConfigFile(pathToUse)
	return nodeConnections
}

var smcAbi abi.ABI

func GetSMCABI(opts ...OptionsBuilder) abi.ABI {
	if smcAbi != nil {
		return smcAbi
	}

	options := BuildOptions(opts...)
	var pathToUse string
	var err error
	if _, ok := options.KV[SMCABITPathOpt]; ok {
		pathToUse = options.KV[SMCABITPathOpt].(string)
	} else {
		pathToUse, err = getABIPathFromArgs()
		if err != nil {
			fmt.Println(fmt.Sprintf("Error while defining smc abi path: %s", err))
			panic(err)
		}

	}

	sbuild, err := ParseJSONFileToSolidityBuild(pathToUse)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error while build solidity build from file %s", pathToUse))
		panic(err)
	}

	smcAbi = sbuild.ABI
	return smcAbi
}

var smcByteCode []byte

func GetSMCBytecode(opts ...OptionsBuilder) []byte {
	if smcByteCode != nil {
		return smcByteCode
	}

	options := BuildOptions(opts...)
	var pathToUse string
	var err error
	if _, ok := options.KV[SMCABITPathOpt]; ok {
		pathToUse = options.KV[SMCABITPathOpt].(string)
	} else {
		pathToUse, err = getABIPathFromArgs()
		if err != nil {
			fmt.Println(fmt.Sprintf("Error while defining smc abi path: %s", err))
			panic(err)
		}
	}
	
	sbuild, err := ParseJSONFileToSolidityBuild(pathToUse)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error while build solidity build from file %s", pathToUse))
		panic(err)
	}

	smcByteCode = sbuild.Bytecode
	return smcByteCode
}
