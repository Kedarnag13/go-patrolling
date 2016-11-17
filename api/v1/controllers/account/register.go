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

type registrationController struct{}

var Registration registrationController

func (r registrationController) Create(rw http.ResponseWriter, req *http.Request) {

	flag := 1

	f := firego.New("https://go-patrolling.firebaseio.com/", nil)
	f.Auth("P0xReX74eqJ6dgZhaujvdamVtzp0o7ik20nLuIGO")

	generate_user_id := fireauth.New("go-patrolling")

	data := fireauth.Data{"uid": "1"}
	options := &fireauth.Option{
		NotBefore:  2,
		Expiration: 3,
		Admin:      false,
		Debug:      true,
	}
	user_id, err := generate_user_id.CreateToken(data, options)
	if err != nil {
		panic(err)
	}

	var user models.User
	var session models.Session
	var device models.Device

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		panic(err)
	}

	mobile_number := `"` + user.MobileNumber + `"`

	var get_user map[string]interface{}
	if err := f.Child("Users").OrderBy("mobile_number").EqualTo(mobile_number).Value(&get_user); err != nil {
		panic(err)
	}

	if len(get_user) == 0 {
		key := []byte("traveling is fun")
		password := []byte(user.Password)
		confirm_password := []byte(user.PasswordConfirmation)

		user = models.User{Id: user_id, FirstName: user.FirstName, LastName: user.LastName, Email: user.Email, MobileNumber: user.MobileNumber, Password: controllers.Encrypt(key, password), PasswordConfirmation: controllers.Encrypt(key, confirm_password), DeviseToken: user.DeviseToken}

		child_user, err := f.Child("Users").Push(user)
		if err != nil || child_user == nil {
			panic(err)
		}

		var get_session map[string]interface{}
		if err := f.Child("Sessions").EqualTo(mobile_number).OrderBy("mobile_number").Value(&get_session); err != nil {
			panic(err)
		}

		generate_sesssion_id := fireauth.New("go-patrolling")

		data = fireauth.Data{"uid": "1"}
		options = &fireauth.Option{
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

		data := fireauth.Data{"uid": "1"}
		options := &fireauth.Option{
			NotBefore:  2,
			Expiration: 3,
			Admin:      false,
			Debug:      true,
		}
		device_id, err := generate_device_id.CreateToken(data, options)
		if err != nil {
			log.Fatal(err)
		}

		if len(get_session) == 0 {
			session = models.Session{Id: session_id, MobileNumber: user.MobileNumber, UserID: user_id, DeviseToken: user.DeviseToken}
			child_session, err := f.Child("Sessions").Push(session)
			if err != nil || child_session == nil {
				panic(err)
			}
			device = models.Device{Id: device_id, Token: user_id, SessionID: session_id}
			child_device, err := f.Child("Devices").Push(device)
			if err != nil || child_device == nil {
				panic(err)
			}
			b, err := json.Marshal(models.Message{
				Success: true,
				Message: "Session created Successfully!",
				Error:   "",
				User:    user,
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
				Error:   "Session already exists!",
			})
			if err != nil {
				panic(err)
			}
			rw.Header().Set("Content-Type", "application/json")
			rw.Write(b)
			flag = 0
			goto end
		}
	}
	if flag == 1 {
		b, err := json.Marshal(models.Message{
			Success: false,
			Message: "",
			Error:   "Mobile Number already exists!",
		})
		if err != nil {
			panic(err)
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(b)
		goto end
	}
end:
	// 	var get_users map[string]interface{}
	// 	if err := f.Child("Users").OrderBy("id").EqualTo(user.Id).Value(&get_users); err != nil {
	// 		panic(err)
	// 	}

	// 	for key := range get_users {
	// 		var get_user map[string]interface{}
	// 		if err := f.Child("Users").Child(key).OrderBy("mobile_number").Value(&get_user); err != nil {
	// 			panic(err)
	// 		}
	// 		if get_user["mobile_number"] == user.MobileNumber {
	// 			b, err := json.Marshal(models.Message{
	// 				Success: false,
	// 				Message: "",
	// 				Error:   "Mobile Number already exists!",
	// 			})
	// 			if err != nil {
	// 				panic(err)
	// 			}
	// 			rw.Header().Set("Content-Type", "application/json")
	// 			rw.Write(b)
	// 			goto end
	// 		}
	// 	}
	// 	if flag == 0 {
	// 		key := []byte("traveling is fun")
	// 		password := []byte(user.Password)
	// 		confirm_password := []byte(user.PasswordConfirmation)

	// 		user = models.User{Id: user_id, FirstName: user.FirstName, LastName: user.LastName, Email: user.Email, MobileNumber: user.MobileNumber, Password: controllers.Encrypt(key, password), PasswordConfirmation: controllers.Encrypt(key, confirm_password), DeviseToken: user.DeviseToken}

	// 		child_user, err := f.Child("Users").Push(user)
	// 		if err != nil || child_user == nil {
	// 			panic(err)
	// 		}
	// 		var get_session map[string]interface{}
	// 		if err := f.Child("Sessions").EqualTo(session.UserID).OrderBy("mobile_number").Value(&get_session); err != nil {
	// 			panic(err)
	// 		}
	// 		generate_sesssion_id := fireauth.New("go-patrolling")

	// 		data = fireauth.Data{"uid": "1"}
	// 		options = &fireauth.Option{
	// 			NotBefore:  2,
	// 			Expiration: 3,
	// 			Admin:      false,
	// 			Debug:      true,
	// 		}
	// 		session_id, err := generate_sesssion_id.CreateToken(data, options)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}

	// 		var get_sessions map[string]interface{}
	// 		if err := f.Child("Sessions").OrderBy("id").EqualTo(session.Id).Value(&get_sessions); err != nil {
	// 			panic(err)
	// 		}

	// 		for key := range get_sessions {
	// 			var get_session map[string]interface{}
	// 			if err := f.Child("Sessions").Child(key).EqualTo(user.MobileNumber).OrderBy("mobile_number").Value(&get_session); err != nil {
	// 				panic(err)
	// 			}
	// 			if get_session["mobile_number"] != nil {
	// 				b, err := json.Marshal(models.Message{
	// 					Success: true,
	// 					Message: "Session already exists!",
	// 					Error:   "",
	// 				})
	// 				if err != nil {
	// 					panic(err)
	// 				}
	// 				rw.Header().Set("Content-Type", "application/json")
	// 				rw.Write(b)
	// 				goto end
	// 			} else {
	// 				session = models.Session{Id: session_id, MobileNumber: user.MobileNumber, UserID: user_id, DeviseToken: user.DeviseToken}
	// 				child_session, err := f.Child("Sessions").Push(session)
	// 				if err != nil || child_session == nil {
	// 					panic(err)
	// 				}
	// 				generate_device_id := fireauth.New("go-patrolling")

	// 				data := fireauth.Data{"uid": "1"}
	// 				options := &fireauth.Option{
	// 					NotBefore:  2,
	// 					Expiration: 3,
	// 					Admin:      false,
	// 					Debug:      true,
	// 				}
	// 				device_id, err := generate_device_id.CreateToken(data, options)
	// 				if err != nil {
	// 					log.Fatal(err)
	// 				}
	// 				device = models.Device{Id: device_id, Token: user_id}
	// 				child_device, err := f.Child("Devices").Push(device)
	// 				if err != nil || child_device == nil {
	// 					panic(err)
	// 				}
	// 				b, err := json.Marshal(models.Message{
	// 					Success: true,
	// 					Message: "Session created Successfully!",
	// 					Error:   "",
	// 				})
	// 				if err != nil {
	// 					panic(err)
	// 				}
	// 				rw.Header().Set("Content-Type", "application/json")
	// 				rw.Write(b)
	// 				flag = 0
	// 				goto end
	// 			}
	// 		}
	// 	}
	// end:
}

// 	if get_user["mobile_number"] == user.MobileNumber {
// 		b, err := json.Marshal(models.Message{
// 			Success: false,
// 			Message: "",
// 			Error:   "Mobile Number already exists!",
// 		})
// 		if err != nil {
// 			panic(err)
// 		}
// 		rw.Header().Set("Content-Type", "application/json")
// 		rw.Write(b)
// 		goto end
// 	} else {

// end:
// }

// 	for _, value := range get_users {
// 		get_user := value.(map[string]interface{})
// 		if get_user["mobile_number"] == user.MobileNumber {
// 			b, err := json.Marshal(models.Message{
// 				Success: false,
// 				Message: "",
// 				Error:   "Mobile Number already exists!",
// 			})
// 			if err != nil {
// 				panic(err)
// 			}
// 			rw.Header().Set("Content-Type", "application/json")
// 			rw.Write(b)
// 			goto end
// 		}
// 	}
// 	if len(get_users) == 0 && get_users == nil {
// 		key := []byte("traveling is fun")
// 		password := []byte(user.Password)
// 		confirm_password := []byte(user.PasswordConfirmation)

// 		user = models.User{Id: user_id, FirstName: user.FirstName, LastName: user.LastName, Email: user.Email, MobileNumber: user.MobileNumber, Password: controllers.Encrypt(key, password), PasswordConfirmation: controllers.Encrypt(key, confirm_password), DeviseToken: user.DeviseToken}

// 		child_user, err := f.Child("Users").Push(user)
// 		if err != nil || child_user == nil {
// 			panic(err)
// 		}

// 		var get_session map[string]interface{}
// 		if err := f.Child("Sessions").EqualTo(session.Id).OrderBy("mobile_number").Value(&get_session); err != nil {
// 			panic(err)
// 		}

// 		generate_sesssion_id := fireauth.New("go-patrolling")

// 		data = fireauth.Data{"uid": "1"}
// 		options = &fireauth.Option{
// 			NotBefore:  2,
// 			Expiration: 3,
// 			Admin:      false,
// 			Debug:      true,
// 		}
// 		session_id, err := generate_sesssion_id.CreateToken(data, options)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		if get_session == nil {
// 			session = models.Session{Id: session_id, MobileNumber: user.MobileNumber, UserID: user_id, DeviseToken: user.DeviseToken}
// 			child_session, err := f.Child("Sessions").Push(session)
// 			if err != nil || child_session == nil {
// 				panic(err)
// 			}
// 			generate_device_id := fireauth.New("go-patrolling")

// 			data := fireauth.Data{"uid": "1"}
// 			options := &fireauth.Option{
// 				NotBefore:  2,
// 				Expiration: 3,
// 				Admin:      false,
// 				Debug:      true,
// 			}
// 			device_id, err := generate_device_id.CreateToken(data, options)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			device = models.Device{Id: device_id, Token: user_id}
// 			child_device, err := f.Child("Devices").Push(device)
// 			if err != nil || child_device == nil {
// 				panic(err)
// 			}
// 			b, err := json.Marshal(models.Message{
// 				Success: true,
// 				Message: "Session created Successfully!",
// 				Error:   "",
// 			})
// 			if err != nil {
// 				panic(err)
// 			}
// 			rw.Header().Set("Content-Type", "application/json")
// 			rw.Write(b)
// 			goto end
// 		} else {
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
// 		}
// 	}
// 	end:
// }
