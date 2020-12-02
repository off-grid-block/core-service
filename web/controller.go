package web

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/off-grid-block/controller"
	"github.com/off-grid-block/core-service/blockchain"
	caMsp "github.com/off-grid-block/fabric-sdk-go/pkg/client/msp"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ControllerManager struct {
	admin *controller.AdminController
	clients ClientControllerStore
}

type ClientControllerStore map[string]*controller.ClientController

// initialize controllers
func NewControllerManager() (*ControllerManager, error) {
	var mgr ControllerManager

	// initialize admin controller
	mgr.admin, _ = controller.NewAdminController()

	//// register admin agent DID with ledger
	//_, err := controller.RegisterDidWithLedger(mgr.admin, Seed())
	//if err != nil {
	//	return nil, err
	//}

	// initialize map
	mgr.clients = make(ClientControllerStore)

	return &mgr, nil
}

// util that generates seeds for did registration with ledger
func Seed() string {
	seed := "my_seed_000000000000000000000000"
	randInt := rand.Intn(800000) + 100000
	seed = seed + strconv.Itoa(randInt)
	return seed[len(seed)-32:]
}


func (app *Application) NewControllerHandler(w http.ResponseWriter, r *http.Request) {

	req := struct {
		Alias string `json:"alias"`
		AgentUrl string `json:"agent_url"`
		Name string `json:"name"`
		Secret string `json:"secret"`
		Type string `json:"type"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to decode new controller request", 500)
		return
	}

	// create new client controller
	cc, _ := controller.NewClientController(req.Alias, req.AgentUrl)

	// register public DID on ledger
	_, err = controller.RegisterDidWithLedger(cc, Seed())
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to register public client did", 500)
		return
	}

	// register app DID with DEON network
	affl := strings.ToLower("org1") + ".department1"
	data := caMsp.RegistrationRequest{
		Name: req.Name,
		Secret: req.Secret,
		Type: req.Type,
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

	_, err = blockchain.Register(app.FabricSDK, data)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to register app with DEON network", 500)
		return
	}

	id := uuid.New().String()
	app.ControllerMgr.clients[id] = cc

	resp := fmt.Sprintf("Client controller ID: %v\n", id)
	w.Write([]byte(resp))
}


func (app *Application) RegisterLedgerHandler(w http.ResponseWriter, r *http.Request) {

	// retrieve controller manager
	mgr := app.ControllerMgr

	vars := mux.Vars(r)
	id := vars["agent_id"]

	if id == "1" {
		_, err := controller.RegisterDidWithLedger(mgr.admin, Seed())
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to register public admin did", 500)
			return
		}
	} else {
		// get client controller with id
		client := mgr.clients[id]
		if client == nil {
			http.Error(w, "Client controller with give id doesn't exist", 400)
			return
		}
		_, err := controller.RegisterDidWithLedger(client, Seed())
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to register public client did", 500)
			return
		}
	}
	w.Write([]byte("Registered public DID"))
}


func (app *Application) EstablishConnectionHandler(w http.ResponseWriter, r *http.Request) {

	// retrieve controller manager
	mgr := app.ControllerMgr

	vars := mux.Vars(r)
	id := vars["agent_id"]
	client := mgr.clients[id]

	invitation, err := controller.CreateInvitation(mgr.admin)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to create connection invitation", 500)
		return
	}

	_, err = controller.ReceiveInvitation(client, invitation)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to establish connection", 500)
		return
	}

	time.Sleep(1 * time.Second)

	// Get connection details of connection between admin and client
	_, err = mgr.admin.GetConnection()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to get connection details for admin", 500)
		return
	}

	_, err = client.GetConnection()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to get connection details for client", 500)
		return
	}

	w.Write([]byte("Connection established"))
}

type IssueCredentialRequest struct {
	AppName string `json:"app_name"`
	AppID string `json:"app_id"`
}

// Issue credential based on DEON app credential definition
func (app *Application) IssueCredentialHandler(w http.ResponseWriter, r *http.Request) {

	mgr := app.ControllerMgr

	schemaID, err := mgr.admin.RegisterSchema("schema")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to register schema", 500)
		return
	}

	_, err = mgr.admin.RegisterCredentialDefinition(schemaID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to register cred def", 500)
		return
	}

	var req IssueCredentialRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Issue credential request badly formed", 400)
		return
	}

	fmt.Printf("app name: %v\n", req.AppName)
	fmt.Printf("app id:   %v\n", req.AppID)

	err = mgr.admin.IssueCredential(req.AppName, req.AppID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to issue credential", 500)
		return
	}

	w.Write([]byte("Issued credential"))
}