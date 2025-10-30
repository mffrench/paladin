all: abi copy-abi

.PHONY: abi
abi:
	cp ../../solidity/artifacts/contracts/tutorials/Storage.sol/Storage.json abis/Storage.json

.PHONY: copy-abi
copy-abi:
	cp ../common/contracts/abis/Storage.json abis/Storage.json

.PHONY: build
build:
	cd cmd/prvstore ; go build -o ../../bin/prs ; cd ../../

.PHONY: run
run:
	./bin/prs ../common-go/integration/testdata/config.json ./abis/Storage.json