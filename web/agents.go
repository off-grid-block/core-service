package web

import (
	"fmt"
	"github.com/off-grid-block/controller"
	"net/http"
	"time"
)

// Issue credential based on DEON app credential definition
func (mgr *ControllerManager) InitializeControllersHandler(w http.ResponseWriter, r *http.Request) {

	var err error

	mgr.admin, err = controller.NewAdminController()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to initialize admin controller", 500)
		return
	}

	mgr.client, err = controller.NewClientController()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to initialize client controller", 500)
		return
	}

	w.Write([]byte("Initialized"))
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

	err = mgr.admin.IssueCredential("voter", "101")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to issue credential", 500)
		return
	}

	w.Write([]byte("Issued credential"))
}