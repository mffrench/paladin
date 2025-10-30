all: abi copy-abi

.PHONY: abi
abi:
	cp ../../solidity/artifacts/contracts/tutorials/HelloWorld.sol/HelloWorld.json abis/HelloWorld.json

.PHONY: copy-abi
copy-abi:
	cp ../common/contracts/abis/HelloWorld.json abis/HelloWorld.json

.PHONY: build
build:
	cd cmd/helloworld ; go build -o ../../bin/hw ; cd ../../

.PHONY: run
run:
	./bin/hw ../common-go/integration/testdata/config.json ./abis/HelloWorld.json