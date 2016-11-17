package account

import (
	"encoding/json"
	"github.com/Kedarnag13/go-patrolling/api/v1/controllers"
	"github.com/Kedarnag13/go-patrolling/api/v1/models"
	"github.com/gorilla/mux"
	"github.com/zabawaba99/fireauth"
	"gopkg.in/zabawaba99/firego.v1"
	"io/ioutil"
	"log"
	"net/http"
)

type SessionController struct{}

var Session SessionController

func (s SessionController) Create(rw http.ResponseWriter, req *http.Request) {

	flag := 1

	f := firego.New("https://go-patrolling.firebaseio.com/", nil)
	f.Auth("P0xReX74eqJ6dgZhaujvdamVtzp0o7ik20nLuIGO")

	var session models.Session
	var device models.Device

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &session)
	if err != nil {
		panic(err)
	}

	generate_sesssion_id := fireauth.New("go-patrolling")

	data := fireauth.Data{"uid": "1"}
	options := &fireauth.Option{
		NotBefore:  2,
		Expiration: 3,
		Admin:      false,
		Debug:      true,
	}
	session_id, err := generate_sesssion_id.CreateToken(data, options)
	if err != nil {
		log.Fatal(err)
	}

	generate_device_id := fireauth.New("go-patrolling")

	data = fireauth.Data{"uid": "1"}
	options = &fireauth.Option{
		NotBefore:  2,
		Expiration: 3,
		Admin:      false,
		Debug:      true,
	}
	device_id, err := generate_device_id.CreateToken(data, options)
	if err != nil {
		panic(err)
	}

	mobile_number := `"` + session.MobileNumber + `"`

	var get_entire_user map[string]interface{}

	if err = f.Child("Users").EqualTo(mobile_number).OrderBy("mobile_number").Value(&get_entire_user); err != nil {
		panic(err)
	}

	for u_key := range get_entire_user {
		get_user := get_entire_user[u_key]
		user := get_user.(map[string]interface{})
		var get_entire_session map[string]interface{}
		if err = f.Child("Sessions").EqualTo(mobile_number).OrderBy("mobile_number").Value(&get_entire_session); err != nil {
			panic(err)
		}
		if len(get_entire_session) == 0 {
			encrypt_decrypt_key := []byte("traveling is fun")
			decrypt_password := controllers.Decrypt(encrypt_decrypt_key, user["password"].(string))
			if session.MobileNumber == user["mobile_number"] && decrypt_password == session.Password {
				session = models.Session{Id: session_id, MobileNumber: session.MobileNumber, UserID: user["id"].(string), DeviseToken: user["devise_token"].(string)}
				child_session, err := f.Child("Sessions").Push(session)
				if err != nil || child_session == nil {
					panic(err)
				}
				device = models.Device{Id: device_id, Token: user["id"].(string), SessionID: session_id}
				child_device, err := f.Child("Devices").Push(device)
				if err != nil || child_device == nil {
					panic(err)
				}
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
				flag = 0
				goto end
			} else {
				b, err := json.Marshal(models.Message{
					Success: false,
					Message: "",
					Error:   "Invalid MobileNumber or Password",
				})
				if err != nil {
					panic(err)
				}
				rw.Header().Set("Content-Type", "application/json")
				rw.Write(b)
				flag = 0
			}
		} else {
			b, err := json.Marshal(models.Message{
				Success: false,
				Message: "",
				Error:   "Session already exists!",
			})
			if err != nil {
				panic(err)
			}
			rw.Header().Set("Content-Type", "application/json")
			rw.Write(b)
			goto end
		}
	}
	if flag == 1 {
		b, err := json.Marshal(models.Message{
			Success: false,
			Message: "",
			Error:   "MobileNumber does not exist!",
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

func (s SessionController) Destroy(rw http.ResponseWriter, req *http.Request) {

	f := firego.New("https://go-patrolling.firebaseio.com/", nil)
	f.Auth("P0xReX74eqJ6dgZhaujvdamVtzp0o7ik20nLuIGO")

	vars := mux.Vars(req)
	mobile_number := `"` + vars["mobile_number"] + `"`

	var get_entire_session map[string]interface{}
	if err := f.Child("Sessions").EqualTo(mobile_number).OrderBy("mobile_number").Value(&get_entire_session); err != nil {
		panic(err)
	}
	for s_key := range get_entire_session {
		key := s_key
		f.Child("Sessions").Child(key).Remove()
		b, err := json.Marshal(models.Message{
			Success: true,
			Message: "",
			Error:   "Logged out Successfully!",
		})
		if err != nil {
			panic(err)
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(b)
	}
	b, err := json.Marshal(models.Message{
		Success: true,
		Message: "",
		Error:   "You are not logged in!",
	})
	if err != nil {
		panic(err)
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(b)
}
