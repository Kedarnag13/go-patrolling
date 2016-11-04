package main

import (
	"github.com/Kedarnag13/go-patrolling/api/v1/controllers/account"
	"github.com/Kedarnag13/go-patrolling/api/v1/controllers/tracker"
	"github.com/Kedarnag13/go-patrolling/api/v1/models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
)

func main() {

	// To Connect with the Database
	db, err := gorm.Open("postgres", "host=localhost user=postgres password=password dbname=go_patrolling_development sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Setup the Tables
	db.CreateTable(&models.User{}, &models.Session{}, &models.Device{}, &models.Tracker{}, &models.Route{})

	// Migrations
	db.AutoMigrate(&models.User{}, &models.Session{}, &models.Device{}, &models.Tracker{}, &models.Route{})

	// Routes
	r := mux.NewRouter()
	r.HandleFunc("/sign_up", account.Registration.Create).Methods("POST")
	r.HandleFunc("/sign_in", account.Session.Create).Methods("POST")
	r.HandleFunc("/record", tracker.Track.Route).Methods("POST")
	http.Handle("/", r)
	log.Printf("main : Started : Listening on: http://localhost:3000")
	http.ListenAndServe("0.0.0.0:3000", nil)
}
