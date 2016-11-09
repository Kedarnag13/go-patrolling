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

	// To Connect with the Database
	// db, err := gorm.Open("postgres", "host=localhost user=postgres password=password dbname=go_patrolling_development sslmode=disable")
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()

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
		log.Fatal(err)
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

	get_mobile_number_url := f.StartAt(user.MobileNumber).EndAt(user.MobileNumber).OrderBy("mobile_number").String()

	find_by_mobile_number, err := http.Get(get_mobile_number_url)
	if err != nil {
		panic(err)
	}

	if find_by_mobile_number == nil {
		log.Printf("I am inside")
		key := []byte("traveling is fun")
		password := []byte(user.Password)
		confirm_password := []byte(user.PasswordConfirmation)

		user = models.User{Id: user_id, FirstName: user.FirstName, LastName: user.LastName, Email: user.Email, MobileNumber: user.MobileNumber, Password: controllers.Encrypt(key, password), PasswordConfirmation: controllers.Encrypt(key, confirm_password), DeviseToken: user.DeviseToken}

		child_user, err := f.Child("Users").Push(user)
		if err != nil {
			panic(err)
		}

		get_user_url := f.EqualTo(user_id).String()

		get_user, err := http.Get(get_user_url)
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

		if get_user == nil {
			session = models.Session{Id: session_id, UserID: user_id, DeviseToken: user.DeviseToken}
			child_session, err := child_user.Child("Session").Push(session)
			if err != nil || child_session == nil {
				panic(err)
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
			device = models.Device{Id: device_id, Token: user_id}
			child_device, err := child_user.Child("Device").Push(device)
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
			goto end
		} else {
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
		}
	} else {
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
	}
end:
}
