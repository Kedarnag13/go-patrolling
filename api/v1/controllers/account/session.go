package account

import (
	"encoding/json"
	"github.com/Kedarnag13/go-patrolling/api/v1/models"
	"gopkg.in/zabawaba99/firego.v1"
	"io/ioutil"
	"net/http"
)

type SessionController struct{}

var Session SessionController

func (s SessionController) Create(rw http.ResponseWriter, req *http.Request) {

	f := firego.New("https://go-patrolling.firebaseio.com/", nil)
	f.Auth("P0xReX74eqJ6dgZhaujvdamVtzp0o7ik20nLuIGO")

	var session models.Session

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &session)
	if err != nil {
		panic(err)
	}

	var get_user map[string]interface{}
	if err := f.Child("Users").EqualTo(session.Id).OrderBy("mobile_number").Value(&get_user); err != nil {
		panic(err)
	}
	for key, value := range get_user {
		mapped_value := value.(map[string]interface{})
		if mapped_value["id"] == nil {
			session = models.Session{UserID: mapped_value["id"].(string), DeviseToken: session.DeviseToken}
			child_user := f.Child("Users")
			if child_user == nil {
				panic(err)
			}
			child_track := child_user.Child(key)
			if child_track == nil {
				panic(err)
			}
			child_track.Child("Tracker").Push(session)
			b, err := json.Marshal(models.Message{
				Success: true,
				Message: "Session created Successfully!",
				Error:   "",
			})
			if err != nil {
				panic(err)
			}
			rw.Header().Set("Content-Type", "application/json")
			rw.Write(b)
			goto end
		} else {
			b, err := json.Marshal(models.Message{
				Success: false,
				Message: "Session already exists!",
				Error:   "",
			})
			if err != nil {
				panic(err)
			}
			rw.Header().Set("Content-Type", "application/json")
			rw.Write(b)
			goto end
		}
	}
end:
}
