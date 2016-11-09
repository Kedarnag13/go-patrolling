package tracker

import (
	"encoding/json"
	"github.com/Kedarnag13/go-patrolling/api/v1/models"
	"github.com/zabawaba99/fireauth"
	"gopkg.in/zabawaba99/firego.v1"
	"io/ioutil"
	"log"
	"net/http"
)

type RecordController struct{}

var Track RecordController

func (r RecordController) Route(rw http.ResponseWriter, req *http.Request) {

	f := firego.New("https://go-patrolling.firebaseio.com/", nil)
	f.Auth("P0xReX74eqJ6dgZhaujvdamVtzp0o7ik20nLuIGO")

	generate_track_id := fireauth.New("go-patrolling")

	data := fireauth.Data{"uid": "1"}
	options := &fireauth.Option{
		NotBefore:  2,
		Expiration: 3,
		Admin:      false,
		Debug:      true,
	}
	track_id, err := generate_track_id.CreateToken(data, options)
	if err != nil {
		panic(err)
	}

	var session models.Session
	var track models.Tracker

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &track)
	if err != nil {
		panic(err)
	}

	log.Printf("%v", track)

	var get_user map[string]interface{}
	if err := f.Child("Users").EqualTo(session.Id).OrderBy("mobile_number").Value(&get_user); err != nil {
		panic(err)
	}
	for key, value := range get_user {
		mapped_value := value.(map[string]interface{})
		track = models.Tracker{Id: track_id, StartLocation: track.StartLocation, Routes: track.Routes, EndLocation: track.EndLocation, UserID: mapped_value["id"].(string)}
		child_user := f.Child("Users")
		if child_user == nil {
			panic(err)
		}
		child_track := child_user.Child(key)
		if child_track == nil {
			panic(err)
		}
		child_track.Child("Tracker").Push(track)
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
		goto end
	}

	// track = models.Tracker{StartLocation: track.StartLocation, Routes: track.Routes, EndLocation: track.EndLocation, UserID: id}

	// f.Push(track)

	// To Connect with the Database
	// db, err := gorm.Open("postgres", "host=localhost user=postgres password=password dbname=go_patrolling_development sslmode=disable")
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()

	// var track models.Tracker
	// var session models.Session

	// body, err := ioutil.ReadAll(req.Body)
	// if err != nil {
	// 	panic(err)
	// }

	// err = json.Unmarshal(body, &track)
	// if err != nil {
	// 	panic(err)
	// }

	// get_user := db.Model(&session).Where("user_id = ?", track.UserID).Select("id").Row()

	// var id int
	// err = get_user.Scan(&id)
	// if err != nil {
	// 	panic(err)
	// }

	// track = models.Tracker{StartLocation: track.StartLocation, Routes: track.Routes, EndLocation: track.EndLocation, UserID: id}

	// // var routes models.Tracker
	// // err = json.Unmarshal(body, &routes)
	// // if err != nil {
	// // 	panic(err)
	// // }

	// // var route string
	// // track.Routes = []string{}
	// // for _, route = range routes.Routes {
	// // 	track.Routes = append(track.Routes, route)
	// // }

	// db.Create(&track)

	// b, err := json.Marshal(models.Message{
	// 	Success: true,
	// 	Message: "Track recorded Successfully!",
	// 	Error:   "",
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// rw.Header().Set("Content-Type", "application/json")
	// rw.Write(b)
end:
}
