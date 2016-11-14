package main

import (
	"github.com/Kedarnag13/go-patrolling/api/v1/controllers/account"
	"github.com/Kedarnag13/go-patrolling/api/v1/controllers/tracker"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {

	// Routes
	r := mux.NewRouter()
	r.HandleFunc("/sign_up", account.Registration.Create).Methods("POST")
	r.HandleFunc("/sign_in", account.Session.Create).Methods("POST")
	r.HandleFunc("/sign_out/", account.Session.Destroy).Methods("GET")
	r.HandleFunc("/record", tracker.Track.Route).Methods("POST")
	http.Handle("/", r)
	log.Printf("main : Started : Listening on: http://localhost:3000")
	http.ListenAndServe("0.0.0.0:3000", nil)
}
