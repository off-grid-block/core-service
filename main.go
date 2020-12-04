package main

import (
	"github.com/off-grid-block/core-service/blockchain"
	"github.com/off-grid-block/core-service/web"
	"fmt"
	// "github.com/pkg/errors"
	ipfs "github.com/ipfs/go-ipfs-api"
	"os"
)

func main() {

	fSetup := blockchain.SetupSDK {
		OrdererID: 			"orderer.example.com",
		ChannelID: 			"mychannel",
		ChannelConfig:		os.Getenv("CHANNEL_CONFIG"),
		ChaincodeGoPath:	os.Getenv("CHAINCODE_GOPATH"),
		ChaincodePath:		make(map[string]string),
		OrgAdmin:			"Admin",
		OrgName:			"org1",
		ConfigFile:			"/src/config.yaml",
		UserName:			"User1",
	}

	err := fSetup.Initialization()
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
		return
	}

	// Close SDK
	defer fSetup.CloseSDK()

	err = fSetup.AdminSetup()
	if err != nil {
		fmt.Printf("Failed to set up network admin: %v\n", err)
		return
	}

	err = fSetup.ChannelSetup()
	if err != nil {
		fmt.Printf("Failed to set up channel: %v\n", err)
		return
	}

	// create shell to connect to IPFS
	sh := ipfs.NewShell(os.Getenv("IPFS_ENDPOINT"))

	// initialize manager for deon network's aca-py agents
	mgr, err := web.NewControllerManager(os.Args[1])
	if err != nil {
		fmt.Printf("Failed to start controller manager: %v\n", err)
		return
	}

	app := &web.Application{
		FabricSDK: &fSetup,
		IpfsShell: sh,
		ControllerMgr: mgr,
	}

	web.Serve(app)
}