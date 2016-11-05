package tracker

import (
	"encoding/json"
	"github.com/Kedarnag13/go-patrolling/api/v1/models"
	"github.com/jinzhu/gorm"
	// "github.com/lib/pq"
	"github.com/zabawaba99/firego"
	"io/ioutil"
	"log"
	"net/http"
)

type RecordController struct{}

var Track RecordController

func (r RecordController) Route(rw http.ResponseWriter, req *http.Request) {

	// To Connect with the Database
	db, err := gorm.Open("postgres", "host=localhost user=postgres password=password dbname=go_patrolling_development sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var track models.Tracker
	var session models.Session

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &track)
	if err != nil {
		panic(err)
	}

	get_user := db.Model(&session).Where("user_id = ?", track.UserID).Select("id").Row()

	var id int
	err = get_user.Scan(&id)
	if err != nil {
		panic(err)
	}

	fb := firego.New("https://go-patrolling.firebaseio.com/", nil)
	fb.Auth("P0xReX74eqJ6dgZhaujvdamVtzp0o7ik20nLuIGO")
	fb.Shallow(true)
	fb.IncludePriority(true)

	// for _, v := range track.Route {
	track = models.Tracker{StartLocation: track.StartLocation, Route: track.Route, EndLocation: track.EndLocation, UserID: id}

	newRef, err := fb.Push(&track)
	if err != nil || newRef == nil {
		log.Fatal(err)
	}

	// db.Create(&track)

	// // log.Println("track", &track)
	b, err := json.Marshal(models.Message{
		Success: true,
		Message: "Track recorded Successfully!",
		Error:   "",
	})
	if err != nil {
		panic(err)
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(b)
}
