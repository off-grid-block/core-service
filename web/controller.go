package web

import (
	"encoding/json"
	"fmt"
	"github.com/off-grid-block/controller"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type ControllerManager struct {
	admin *controller.AdminController
	client *controller.ClientController
}

// Check if Controller Manager is initialized
func (mgr *ControllerManager) Initialized() bool {
	return (mgr.admin != nil) && (mgr.client != nil)
}

// initialize controllers
func NewControllerManager() *ControllerManager {
	var mgr ControllerManager

	mgr.admin, _ = controller.NewAdminController()
	mgr.client, _ = controller.NewClientController()

	return &mgr
}


// util that generates seeds for did registration with ledger
func Seed() string {
	seed := "my_seed_000000000000000000000000"
	randInt := rand.Intn(800000) + 100000
	seed = seed + strconv.Itoa(randInt)
	return seed[len(seed)-32:]
}

func (mgr *ControllerManager) RegisterPublicDidHandler(w http.ResponseWriter, r *http.Request) {

	var err error

	_, err = controller.RegisterDidWithLedger(mgr.admin, Seed())
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to register public admin did", 500)
	}

	_, err = controller.RegisterDidWithLedger(mgr.client, Seed())
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to register public client did", 500)
	}

	w.Write([]byte("Registered public DIDs"))
}


func (mgr *ControllerManager) EstablishConnectionHandler(w http.ResponseWriter, r *http.Request) {

	if !mgr.Initialized() {
		http.Error(w, "Initialize controllers before establishing connection", 500)
		return
	}

	invitation, err := controller.CreateInvitation(mgr.admin)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to create connection invitation", 500)
		return
	}

	_, err = controller.ReceiveInvitation(mgr.client, invitation)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to establish connection", 500)
		return
	}

	time.Sleep(2 * time.Second)

	// Get connection details of connection between admin and client
	_, err = mgr.admin.GetConnection()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to get connection details for admin", 500)
		return
	}

	_, err = mgr.client.GetConnection()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to get connection details for client", 500)
		return
	}

	w.Write([]byte("Connection established"))
}

type IssueCredentialRequest struct {
	appName string `json:"app_name"`
	appID string `json:"app_id"`
}

// Issue credential based on DEON app credential definition
func (mgr *ControllerManager) IssueCredentialHandler(w http.ResponseWriter, r *http.Request) {

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

	err = mgr.admin.IssueCredential(req.appName, req.appID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to issue credential", 500)
		return
	}

	w.Write([]byte("Issued credential"))
}