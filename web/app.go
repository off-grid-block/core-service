package web

import (
	"github.com/off-grid-block/core-service/blockchain"
	ipfs "github.com/ipfs/go-ipfs-api"
	"github.com/gorilla/mux"
	"fmt"
	"log"
	"net/http"
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

	/*********************************/
	/*	subrouter for "poll" prefix  */
	/*********************************/
	// poll := api.PathPrefix("/poll").Subrouter()

	// // handler for initPoll
	// poll.HandleFunc("", voteApp.InitPollHandler).Methods("POST")

	// // handler for queryAllPolls
	// poll.HandleFunc("", voteApp.QueryAllPollsHandler).Methods("GET")

	// // handler for updatePollStatus
	// poll.HandleFunc("/{pollid}/status", voteApp.UpdatePollStatusHandler).Methods("PUT")

	// // handler for getPoll
	// poll.HandleFunc("/{pollid}", voteApp.GetPollHandler).Methods("GET")

	// // handler for getPollPrivateDetails
	// poll.HandleFunc("/{pollid}/private", voteApp.getPollPrivateDetailsHandler).Methods("GET")

	/*********************************/
	/*	subrouter for "vote" prefix  */
	/*********************************/
	// vote := api.PathPrefix("/vote").Subrouter()

	// // handler for initVote
	// vote.HandleFunc("", voteApp.InitVoteHandler).Methods("POST")

	// // handler for getVotePrivateDetails
	// vote.HandleFunc("/{pollid}/{voterid}/private", voteApp.getVotePrivateDetailsHandler).Methods("GET")

	// // handler for getVotePrivateDetailsHash
	// vote.HandleFunc("/{pollid}/{voterid}/hash", voteApp.getVotePrivateDetailsHashHandler).Methods("GET")

	// // handler for getVote
	// vote.HandleFunc("/{pollid}/{voterid}", voteApp.GetVoteHandler).Methods("GET")

	// // handler for queryVotePrivateDetailsByPoll
	// vote.HandleFunc("", voteApp.QueryVotePrivateDetailsByPollHandler).
	// 	Methods("GET").
	// 	Queries("type", "private", "pollid", "{pollid}")

	// // handler for queryVotesByPoll
	// vote.HandleFunc("", voteApp.QueryVotesByPollHandler).
	// 	Methods("GET").
	// 	Queries("type", "public", "pollid", "{pollid}")

	// // handler for getVotesByVoter
	// vote.HandleFunc("", voteApp.QueryVotesByVoterHandler).
	// 	Methods("GET").
	// 	Queries("voterid", "{voterid}")

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