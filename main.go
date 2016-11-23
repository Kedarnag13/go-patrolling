package main

import (
	"github.com/Kedarnag13/go-patrolling/api/v1/controllers/account"
	"github.com/Kedarnag13/go-patrolling/api/v1/controllers/tracker"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func main() {

	// Routes
	r := mux.NewRouter()
	r.HandleFunc("/sign_up", account.Registration.Create).Methods("POST")
	r.HandleFunc("/sign_in", account.Session.Create).Methods("POST")
	r.HandleFunc("/sign_out/{mobile_number:([0-9]+)?}", account.Session.Destroy).Methods("GET")
	r.HandleFunc("/record", tracker.Track.Route).Methods("POST")
	r.HandleFunc("/get_routes_for/{mobile_number:([0-9]+)?}", tracker.Track.Get_Routes_For).Methods("GET")
	handler := cors.Default().Handler(r)
	http.Handle("/", handler)
	log.Printf("main : Started : Listening on: http://localhost:3010")
	http.ListenAndServe("0.0.0.0:3010", nil)
}
