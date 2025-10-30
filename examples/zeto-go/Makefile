all: abi copy-abi

.PHONY: abi
.PHONY: copy-abi
copy-abi:
	cp ../common/contracts/zeto-abis/SampleERC20.json abis/SampleERC20.json

.PHONY: build
build:
	cd cmd/zeto ; go build -o ../../bin/zeto ; cd ../../

.PHONY: run
run:
	./bin/zeto ../common-go/integration/testdata/config.json ./abis/SampleERC20.json
