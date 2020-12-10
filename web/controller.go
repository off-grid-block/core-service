package web

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/off-grid-block/controller"
	"github.com/off-grid-block/core-service/blockchain"
	caMsp "github.com/off-grid-block/fabric-sdk-go/pkg/client/msp"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ControllerManager struct {
	admin *controller.AdminController
	clients ClientControllerStore
	ledgerUrl string
}

type ClientController struct {
	*controller.ClientController
	Initialized bool
}

type ClientControllerStore map[string]*ClientController

// initialize controllers
func NewControllerManager(ledgerUrl string) (*ControllerManager, error) {
	var mgr ControllerManager

	// initialize map
	mgr.clients = make(ClientControllerStore)

	// add ledger url
	mgr.ledgerUrl = ledgerUrl

	return &mgr, nil
}

// util that generates seeds for did registration with ledger
func Seed() string {
	seed := "my_seed_000000000000000000000000"
	randInt := rand.Intn(800000) + 100000
	seed = seed + strconv.Itoa(randInt)
	return seed[len(seed)-32:]
}

func (app *Application) newAdminController(req NewControllerRequest) (string, error) {
	var err error
	app.ControllerMgr.admin, err = controller.NewAdminController("http://admin.example.com:8021")
	if err != nil {
		return "", err
	}

	// register admin agent DID with ledger
	did, err := controller.RegisterDidWithLedger(app.ControllerMgr.admin, Seed(), app.ControllerMgr.ledgerUrl)
	if err != nil {
		return "", err
	}

	return did, nil
}

func (app *Application) newClientController(req NewControllerRequest) (string, error) {

	// create new client controller
	cc, _ := controller.NewClientController(req.Alias, req.AgentUrl)

	// register public DID on ledger
	_, err := controller.RegisterDidWithLedger(cc, Seed(), app.ControllerMgr.ledgerUrl)
	if err != nil {
		return "", err
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
		return "", err
	}

	id := uuid.New().String()

	app.ControllerMgr.clients[id] = &ClientController{cc, false}

	return id, nil
}

type NewControllerRequest struct {
	AgentType string `json:"agent_type"`
	Alias string `json:"alias"`
	AgentUrl string `json:"agent_url"`
	Name string `json:"name"`
	Secret string `json:"secret"`
	Type string `json:"type"`
}

func (app *Application) NewControllerHandler(w http.ResponseWriter, r *http.Request) {

	var req NewControllerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "New controller request body badly formed", 400)
		return
	}

	var resp string

	if req.AgentType == "client" {
		id, err := app.newClientController(req)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to create new client controller", 500)
			return
		}
		log.Printf("Agent initialized\n")
		resp = fmt.Sprintf("Client controller ID: %v\n", id)

	} else if req.AgentType == "admin" {
		did, err := app.newAdminController(req)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to create new admin controller", 500)
			return
		}
		log.Printf("Agent initialized with DID: %v\n", did)
		resp = fmt.Sprintf("Created admin controller")
	}

	w.Write([]byte(resp))
}



func (app *Application) GetControllerByAliasHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	alias := vars["alias"]

	type Response struct {
		Initialized bool `json:"initialized"`
	}

	for _, client := range app.ControllerMgr.clients {
		if client.Alias() == alias {

			body := Response{Initialized: client.Initialized}
			resp, err := json.Marshal(body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.Write([]byte(resp))
			return
		}
	}

	// if no client found
	resp, _ := json.Marshal(false)
	w.Write([]byte(resp))
}

// func (app *Application) RegisterLedgerHandler(w http.ResponseWriter, r *http.Request) {

// 	// retrieve controller manager
// 	mgr := app.ControllerMgr

// 	vars := mux.Vars(r)
// 	id := vars["agent_id"]

// 	if id == "1" {
// 		_, err := controller.RegisterDidWithLedger(mgr.admin, Seed(), mgr.ledgerUrl)
// 		if err != nil {
// 			fmt.Println(err)
// 			http.Error(w, "Failed to register public admin did", 500)
// 			return
// 		}
// 	} else {
// 		// get client controller with id
// 		client := mgr.clients[id]
// 		if client == nil {
// 			http.Error(w, "Client controller with give id doesn't exist", 400)
// 			return
// 		}
// 		_, err := controller.RegisterDidWithLedger(client, Seed(), mgr.ledgerUrl)
// 		if err != nil {
// 			fmt.Println(err)
// 			http.Error(w, "Failed to register public client did", 500)
// 			return
// 		}
// 	}
// 	w.Write([]byte("Registered public DID"))
// }


func (app *Application) EstablishConnectionHandler(w http.ResponseWriter, r *http.Request) {

	// retrieve controller manager
	mgr := app.ControllerMgr

	vars := mux.Vars(r)
	id := vars["agent_id"]
	client, found := mgr.clients[id]
	if !found {
		http.Error(w, "Client agent not found", 400)
		return
	}

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

	vars := mux.Vars(r)
	id := vars["agent_id"]
	client, found := mgr.clients[id]
	if !found {
		http.Error(w, "Client agent not found", 400)
		return
	}
	schemaID, err := mgr.admin.RegisterSchema("deon_schema")
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

	// get app_name and app_id from request body
	var req IssueCredentialRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Issue credential request badly formed", 400)
		return
	}

	log.Printf("app name: %v\n", req.AppName)
	log.Printf("app id:   %v\n", req.AppID)

	credExID, err := mgr.admin.IssueCredential(req.AppName, req.AppID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to issue credential", 500)
		return
	}
	log.Printf("Credential Issued with ID: %v\n", credExID)

	client.Initialized = true
	w.Write([]byte("Issued credential"))
}