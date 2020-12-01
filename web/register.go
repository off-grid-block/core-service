package web

import (
	"encoding/json"
	"net/http"
	"strings"
	"fmt"

	"github.com/off-grid-block/core-service/blockchain"

	caMsp "github.com/off-grid-block/fabric-sdk-go/pkg/client/msp"
)

func (app *Application) UserHandler(w http.ResponseWriter, r *http.Request) {

	affl := strings.ToLower("org1") + ".department1"

	data := caMsp.RegistrationRequest{
		Name: "email",
		Secret: "password",
		Type: "peer",
		MaxEnrollments: -1,
		Affiliation: affl,
		Attributes: []caMsp.Attribute{
			{
				Name: "role",
				Value: "user",
				ECert: true,
			},
		},
		CAName: "ca.org1.example.com",
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		panic(err)
	}

	_, err = blockchain.Register(app.FabricSDK, data)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to register app with DEON network", 500)
	}
}
