package web

import (
	"fmt"
	"github.com/gorilla/mux"
	ipfs "github.com/ipfs/go-ipfs-api"
	"github.com/off-grid-block/core-service/blockchain"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
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
	proxyPath := getProxyPath(proxyUrl, r)
	serveRedirect(proxyUrl, proxyPath, w, r)
}

func getProxyPath(proxyUrl string, r *http.Request) string {

	if proxyUrl == os.Getenv("CLIENT_ACAPY_AGENT_URL") {
		return strings.TrimPrefix(r.URL.Path, "/api/v1/client")

	} else if proxyUrl == os.Getenv("CI_MSP_ACAPY_AGENT_URL") {
		return strings.TrimPrefix(r.URL.Path, "/api/v1/ci-msp")
	}

	return r.URL.Path
}

// Retrieve the correct host:port for the specified application
func getProxyUrl(r *http.Request) string {

	fmt.Println("Path: " + r.URL.Path)

	if strings.HasPrefix(r.URL.Path, "/api/v1/vote") || strings.HasPrefix(r.URL.Path, "/api/v1/poll") {
		fmt.Println("Redirecting to vote service...")
		return os.Getenv("VOTE_URL")

	} else if strings.HasPrefix(r.URL.Path, "/api/v1/client") {
		fmt.Println("Redirecting to client ACA-Py agent...")
		return os.Getenv("CLIENT_ACAPY_AGENT_URL")

	} else if strings.HasPrefix(r.URL.Path, "/api/v1/ci-msp") {
		fmt.Println("Redirecting to CI/MSP ACA-Py agent...")
		return os.Getenv("CI_MSP_ACAPY_AGENT_URL")

	} else {
		log.Fatalf("Failed to match path: %v\n", r.URL.Path)
		return ""
	}
}

// Redirect the request to the URL specified
func serveRedirect(host string, path string, w http.ResponseWriter, r *http.Request) {

	u, err := url.Parse(host)
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(u)
	
	r.URL.Path = path
	r.URL.Host = u.Host
	r.URL.Scheme = u.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = u.Host

	proxy.ServeHTTP(w, r)
}

// Serve core DEON service API
func Serve(app *Application) {

	// initialize controller manager
	mgr := NewControllerManager()

	// create router object
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()

	// identity management endpoints
	api.HandleFunc("/admin/app/register", app.UserHandler).Methods("POST")
	api.HandleFunc("/admin/agents/connect", mgr.EstablishConnectionHandler).Methods("POST")
	api.HandleFunc("/admin/agents/register-public-did", mgr.RegisterPublicDidHandler).Methods("POST")
	api.HandleFunc("/admin/agents/issue-credential", mgr.IssueCredentialHandler).Methods("POST")

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
