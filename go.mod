module github.com/off-grid-block/core-service

go 1.13

replace github.com/off-grid-block/vote => /Users/brianli/deon/src/vote

replace github.com/off-grid-block/core-interface => /Users/brianli/deon/core-interface

require (
	github.com/gorilla/mux v1.7.4
	github.com/hyperledger/fabric-protos-go v0.0.0-20200124220212-e9cfc186ba7b
	github.com/hyperledger/fabric-sdk-go v1.0.0-beta2
	github.com/ipfs/go-ipfs-api v0.0.3
	github.com/off-grid-block/core-interface v0.0.0-20200614195207-2ed65a989bd0
	github.com/off-grid-block/vote v0.0.0-20200614200315-e2f36038c477
	github.com/pkg/errors v0.9.1
)
