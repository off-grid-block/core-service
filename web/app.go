package web

import (
	"github.com/off-grid-block/core-service/blockchain"
	ipfs "github.com/ipfs/go-ipfs-api"
	"github.com/gorilla/mux"
	"fmt"
	"log"
	"net/http"
	// "net/http/httputil"
	// "net/url"
	"time"
	// "os"
	// "strings"

)

type Application struct {
	FabricSDK *blockchain.SetupSDK
	IpfsShell *ipfs.Shell
}

// Homepage
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Homepage\n"))
}

func Serve(app *Application) {
	// create router object
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()

	// test api homepage
	api.HandleFunc("/", HomeHandler)

	/********************************/
	/* identity management endpoint */
	/********************************/
	api.HandleFunc("/register", app.UserHandler).Methods("POST")

	// Start http server
	srv := &http.Server{
		Handler: 	r,
		Addr:		"0.0.0.0:8001",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("Listening on %v...\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
	
}
