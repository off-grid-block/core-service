package web

import (
	"github.com/off-grid-block/deon-library/sdk"
	// ipfs "github.com/ipfs/go-ipfs-api"
)

type Application struct {
	FabricSDK *sdk.SDKConfig
}

// Set up DEON Admin app
func SetupApp(config *sdk.SDKConfig) *Application {
	return &Application{config}
}