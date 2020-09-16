module github.com/off-grid-block/core-service

go 1.13

replace github.com/off-grid-block/core-interface => /pkg/core-interface
replace github.com/off-grid-block/vote => /src/vote

require (
	github.com/gorilla/mux v1.7.3
	github.com/hyperledger/fabric-protos-go v0.0.0-20200124220212-e9cfc186ba7b
	github.com/hyperledger/fabric-sdk-go v1.0.0-beta2
	github.com/ipfs/go-ipfs-api v0.0.3
	github.com/off-grid-block/core-interface v0.0.0-20200915132455-14563a2cd9b8
	github.com/off-grid-block/vote v0.0.0-20200915132652-dfea0a1f8428
	github.com/pkg/errors v0.9.1
)
