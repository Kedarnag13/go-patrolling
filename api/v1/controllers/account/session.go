package account

import (
	"encoding/json"
	// "github.com/Kedarnag13/go-patrolling/api/v1/controllers"
	"github.com/Kedarnag13/go-patrolling/api/v1/models"
	"gopkg.in/zabawaba99/firego.v1"
	"io/ioutil"
	"log"
	"net/http"
)

type SessionController struct{}

var Session SessionController

func (s SessionController) Create(rw http.ResponseWriter, req *http.Request) {

	// flag := 0

	f := firego.New("https://go-patrolling.firebaseio.com/", nil)
	f.Auth("P0xReX74eqJ6dgZhaujvdamVtzp0o7ik20nLuIGO")

	var user models.User
	var session models.Session

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &session)
	if err != nil {
		panic(err)
	}

	var get_sessions map[string]interface{}
	if err := f.Child("Sessions").EqualTo(user.MobileNumber).OrderBy("mobile_number").Value(&get_sessions); err != nil {
		panic(err)
	}
	for _, value := range get_sessions {
		get_session := value.(map[string]interface{})
		if get_session["devise_token"] == session.DeviseToken {
			b, err := json.Marshal(models.Message{
				Success: true,
				Message: "Session already exists!",
				Error:   "",
			})
			if err != nil {
				panic(err)
			}
			rw.Header().Set("Content-Type", "application/json")
			rw.Write(b)
			goto end
		} else {
			var get_user map[string]interface{}
			if err := f.Child("Users").EqualTo(user.MobileNumber).OrderBy("mobile_number").Value(&get_user); err != nil {
				panic(err)
			}
			log.Println(get_user)
		}
	}
end:
}

func (s SessionController) Destroy(rw http.ResponseWriter, req *http.Request) {

}
