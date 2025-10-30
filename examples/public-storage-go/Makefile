all: abi copy-abi

.PHONY: abi
abi:
	cp ../../solidity/artifacts/contracts/tutorials/Storage.sol/Storage.json abis/Storage.json

.PHONY: copy-abi
copy-abi:
	cp ../common/contracts/abis/Storage.json abis/Storage.json

.PHONY: build
build:
	cd cmd/pubstore ; go build -o ../../bin/pbs ; cd ../../

.PHONY: run
run:
	./bin/pbs ../common-go/integration/testdata/config.json ./abis/Storage.json