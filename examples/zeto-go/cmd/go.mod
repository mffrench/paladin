module github.com/LF-Decentralized-Trust-labs/paladin/examples/zeto

go 1.24.0

toolchain go1.24.3

replace github.com/LF-Decentralized-Trust-labs/paladin/examples/common-go => ../../common-go

replace github.com/LF-Decentralized-Trust-labs/paladin/common/go => ../../../common/go

replace github.com/LF-Decentralized-Trust-labs/paladin/config => ../../../config

replace github.com/LF-Decentralized-Trust-labs/paladin/sdk/go => ../../../sdk/go

require (
	github.com/LF-Decentralized-Trust-labs/paladin/examples/common-go v0.0.0-00010101000000-000000000000
	github.com/LF-Decentralized-Trust-labs/paladin/sdk/go v0.0.0-00010101000000-000000000000
)

require (
	github.com/LF-Decentralized-Trust-labs/paladin/common/go v0.0.0-00010101000000-000000000000 // indirect
	github.com/LF-Decentralized-Trust-labs/paladin/config v0.0.0-00010101000000-000000000000 // indirect
	github.com/aidarkhanov/nanoid v1.0.8 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/go-resty/resty/v2 v2.14.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/hyperledger/firefly-common v1.5.5 // indirect
	github.com/hyperledger/firefly-signer v1.1.22 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/term v0.30.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
