package main

import (
	"github.com/Kedarnag13/go-patrolling/api/v1/controllers/account"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/sign_up", account.Registration.Create).Methods("POST")
	http.Handle("/", r)
	log.Printf("main : Started : Listening on: http://localhost:3000")
	http.ListenAndServe("0.0.0.0:3000", nil)
}
