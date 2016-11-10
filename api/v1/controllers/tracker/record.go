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
end:
}
