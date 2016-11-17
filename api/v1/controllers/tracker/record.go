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

	var track models.Tracker

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &track)
	if err != nil {
		panic(err)
	}

	var get_session map[string]interface{}
	if err = f.Child("Sessions").EqualTo(mobile_number).OrderBy("mobile_number").Value(&get_session); err != nil {
		panic(err)
	}

	log.Println(get_session)

	if len(get_session) == 0 {
		b, err := json.Marshal(models.Message{
			Success: false,
			Message: "",
			Error:   "You need to be logged in!",
		})
		if err != nil {
			panic(err)
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(b)
		goto end
	} else {
		track = models.Tracker{Id: track_id, StartLocation: track.StartLocation, Routes: track.Routes, EndLocation: track.EndLocation, MobileNumber: track.MobileNumber}
		child_track, err := f.Child("Tracker").Push(track)
		if err != nil || child_track == nil {
			panic(err)
		}
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
