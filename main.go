package main

import (
	"github.com/off-grid-block/core-service/web"
	"github.com/off-grid-block/vote"
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"time"
	// "github.com/off-grid-block/core-service/blockchain"
	"github.com/off-grid-block/core-interface/pkg/sdk"
	// "github.com/pkg/errors"
	"fmt"
)


func main() {

	// Initialize SDK
	fabricSDK, err := sdk.SetupSDK()
	if err != nil {
		fmt.Errorf("Failed to set up and initialize SDK: %v", err)
	}

	// set up admin app
	adminApp := web.SetupApp(fabricSDK)

	// set up vote app
	voteApp, err := vote.SetupApp(fabricSDK)
	if err != nil {
		fmt.Errorf("Failed to set up voting app: %v", err)
		return
	}

	err = fabricSDK.ChainCodeInstallationInstantiation()
	if err != nil {
		fmt.Errorf("Failed to install & instantiate chaincode: %v", err)
		return
	}

	// create router object
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()

	/********************************/
	/* identity management endpoint */
	/********************************/
	api.HandleFunc("/application", adminApp.UserHandler).Methods("POST")

	/*********************************/
	/*	subrouter for "poll" prefix  */
	/*********************************/
	poll := api.PathPrefix("/poll").Subrouter()

	// handler for initPoll
	poll.HandleFunc("", voteApp.InitPollHandler).Methods("POST")

	// handler for queryAllPolls
	poll.HandleFunc("", voteApp.QueryAllPollsHandler).Methods("GET")

	// handler for updatePollStatus
	poll.HandleFunc("/{pollid}/status", voteApp.UpdatePollStatusHandler).Methods("PUT")

	// handler for getPoll
	poll.HandleFunc("/{pollid}", voteApp.GetPollHandler).Methods("GET")

	// // handler for getPollPrivateDetails
	// poll.HandleFunc("/{pollid}/private", voteApp.getPollPrivateDetailsHandler).Methods("GET")

	/*********************************/
	/*	subrouter for "vote" prefix  */
	/*********************************/
	vote := api.PathPrefix("/vote").Subrouter()

	// handler for initVote
	vote.HandleFunc("", voteApp.InitVoteHandler).Methods("POST")

	// // handler for getVotePrivateDetails
	// vote.HandleFunc("/{pollid}/{voterid}/private", voteApp.getVotePrivateDetailsHandler).Methods("GET")

	// // handler for getVotePrivateDetailsHash
	// vote.HandleFunc("/{pollid}/{voterid}/hash", voteApp.getVotePrivateDetailsHashHandler).Methods("GET")

	// handler for getVote
	vote.HandleFunc("/{pollid}/{voterid}", voteApp.GetVoteHandler).Methods("GET")

	// handler for queryVotePrivateDetailsByPoll
	vote.HandleFunc("", voteApp.QueryVotePrivateDetailsByPollHandler).
		Methods("GET").
		Queries("type", "private", "pollid", "{pollid}")

	// handler for queryVotesByPoll
	vote.HandleFunc("", voteApp.QueryVotesByPollHandler).
		Methods("GET").
		Queries("type", "public", "pollid", "{pollid}")

	// handler for getVotesByVoter
	vote.HandleFunc("", voteApp.QueryVotesByVoterHandler).
		Methods("GET").
		Queries("voterid", "{voterid}")

	// Start http server
	srv := &http.Server{
		Handler: 	r,
		Addr:		"127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("Listening on %v...\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())

}