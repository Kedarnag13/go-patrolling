package tracker

import (
"encoding/json"
"github.com/Kedarnag13/go-patrolling/api/v1/models"
"github.com/jinzhu/gorm"
"io/ioutil"
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

  track = models.Tracker{StartLocation: track.StartLocation, Routes: track.Routes, EndLocation: track.EndLocation, UserID: id}

  db.Create(&track)

  b, err := json.Marshal(models.Message{
    Success: true,
    Message: "Track recorded Successfully!",
    Error:   err.Error(),
    })
  if err != nil {
    panic(err)
  }
  rw.Header().Set("Content-Type", "application/json")
  rw.Write(b)
}
