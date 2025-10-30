#!/bin/sh

cd ./testinfra/ ; go mod tidy ; cd -
cd ./toolkit/go/ ; go mod tidy ; cd -
cd ./registries/static/ ; go mod tidy ; cd -
cd ./registries/evm/ ; go mod tidy ; cd -
cd ./signingmodules/example/ ; go mod tidy ; cd -
cd ./core/go/ ; go mod tidy ; cd -
cd ./common/go/ ; go mod tidy ; cd -
cd ./perf/ ; go mod tidy ; cd -
cd ./operator/ ; go mod tidy ; cd -
cd ./domains/integration-test/ ; go mod tidy ; cd -
cd ./domains/zeto/ ; go mod tidy ; cd -
cd ./domains/noto/ ; go mod tidy ; cd -
cd ./config/ ; go mod tidy ; cd -
cd ./transports/grpc/ ; go mod tidy ; cd -
cd ./sdk/go/ ; go mod tidy ; cd -
cd ./examples/common-go/ ; go mod tidy ; cd -
cd ./examples/helloworld-go/cmd ; go mod tidy ; cd -
cd ./examples/public-storage-go/cmd ; go mod tidy ; cd -
cd ./examples/privacy-storage-go/cmd ; go mod tidy ; cd -