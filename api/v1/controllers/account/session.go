package account

import (
	"encoding/json"
	"github.com/Kedarnag13/go-patrolling/api/v1/controllers"
	"github.com/Kedarnag13/go-patrolling/api/v1/models"
	"github.com/zabawaba99/fireauth"
	"gopkg.in/zabawaba99/firego.v1"
	"io/ioutil"
	"log"
	"net/http"
)

type SessionController struct{}

var Session SessionController

func (s SessionController) Create(rw http.ResponseWriter, req *http.Request) {

	flag := 0

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
		log.Fatal(err)
	}

	mobile_number := `"` + session.MobileNumber + `"`

	var get_entire_user map[string]interface{}

	if err = f.Child("Users").EqualTo(mobile_number).OrderBy("mobile_number").Value(&get_entire_user); err != nil {
		panic(err)
	}

	for u_key := range get_entire_user {
		get_user := get_entire_user[u_key]
		user := get_user.(map[string]interface{})
		if user["mobile_number"] == session.MobileNumber {
			var get_entire_session map[string]interface{}
			if err = f.Child("Sessions").EqualTo(mobile_number).OrderBy("mobile_number").Value(&get_entire_session); err != nil {
				panic(err)
			}
			if get_entire_session == nil {
				encrypt_decrypt_key := []byte("traveling is fun")
				decrypt_password := controllers.Decrypt(encrypt_decrypt_key, user["password"].(string))
				for s_key := range get_entire_session {
					get_session := get_entire_session[s_key]
					user_session := get_session.(map[string]interface{})
					if user_session["password"] == decrypt_password && session.MobileNumber == user_session["mobile_number"] {
						session = models.Session{Id: session_id, MobileNumber: mobile_number, UserID: user_session["user_id"].(string), DeviseToken: user_session["devise_token"].(string)}
						child_session, err := f.Child("Sessions").Push(session)
						if err != nil || child_session == nil {
							panic(err)
						}
						device = models.Device{Id: device_id, Token: user_session["user_id"].(string)}
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
					}
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
					goto end
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

// end:

// 	var get_sessions map[string]interface{}
// 	if err := f.Child("Sessions").EqualTo(user.MobileNumber).OrderBy("mobile_number").Value(&get_sessions); err != nil {
// 		panic(err)
// 	}
// 	for _, value := range get_sessions {
// 		get_session := value.(map[string]interface{})
// 		if get_session["devise_token"] == session.DeviseToken {
// 			b, err := json.Marshal(models.Message{
// 				Success: true,
// 				Message: "Session already exists!",
// 				Error:   "",
// 			})
// 			if err != nil {
// 				panic(err)
// 			}
// 			rw.Header().Set("Content-Type", "application/json")
// 			rw.Write(b)
// 			goto end
// 		} else {
// 			var get_user map[string]interface{}
// 			if err := f.Child("Users").EqualTo(user.MobileNumber).OrderBy("mobile_number").Value(&get_user); err != nil {
// 				panic(err)
// 			}
// 			log.Println(get_user)
// 		}
// 	}
// end:
// }

func (s SessionController) Destroy(rw http.ResponseWriter, req *http.Request) {

}
