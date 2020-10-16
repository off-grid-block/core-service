package web

import (
	"github.com/off-grid-block/core-service/blockchain"
	ipfs "github.com/ipfs/go-ipfs-api"
	"github.com/gorilla/mux"
	"fmt"
	"log"
	"net/url"
	"net/http"
	"net/http/httputil"
	"time"
	"os"
	"strings"
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

// Proxy handler to redirect requests
func RequestAndRedirectHandler(w http.ResponseWriter, r *http.Request) {
	proxyUrl := getProxyUrl(r)
	serveRedirect(proxyUrl, w, r)
}

// Retrieve the correct host:port for the specified application
func getProxyUrl(r *http.Request) string {

	fmt.Println("Path: " + r.URL.Path)

	if strings.HasPrefix(r.URL.Path, "/api/v1/vote") || strings.HasPrefix(r.URL.Path, "/api/v1/poll") {
		fmt.Println("Redirecting to vote service...")
		return os.Getenv("VOTE_URL")

	} else {
		log.Fatalf("Failed to match path: %v\n", r.URL.Path)
		return ""
	}
}

// Redirect the request to the URL specified
func serveRedirect(host string, w http.ResponseWriter, r *http.Request) {

	u, err := url.Parse(host)
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(u)
	
	r.URL.Host = u.Host
	r.URL.Scheme = u.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = u.Host

	proxy.ServeHTTP(w, r)
}


func Serve(app *Application) {
	// create router object
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()

	/********************************/
	/* identity management endpoint */
	/********************************/
	api.HandleFunc("/register", app.UserHandler).Methods("POST")

	// api.HandleFunc("/", HomeHandler)

	// Redirect requests using reverse proxy
	api.PathPrefix("/").HandlerFunc(RequestAndRedirectHandler)

	// Start http server
	srv := &http.Server{
		Handler: 	r,
		Addr:		os.Getenv("CORE_URL"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("Listening on %v...\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
	
}
