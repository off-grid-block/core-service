package web

import (
	"github.com/off-grid-block/core-interface/pkg/sdk"
	// ipfs "github.com/ipfs/go-ipfs-api"
)

type Application struct {
	FabricSDK *sdk.SDKConfig
}

// Set up DEON Admin app
func SetupApp(config *sdk.SDKConfig) *Application {
	return &Application{config}
}