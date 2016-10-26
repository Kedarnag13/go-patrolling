package account

import (
	"encoding/json"
	"github.com/Kedarnag13/go-patrolling/api/v1/controllers"
	"github.com/Kedarnag13/go-patrolling/api/v1/models"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	// "log"
	"net/http"
)

type registrationController struct{}

var Registration registrationController

func (r registrationController) Create(rw http.ResponseWriter, req *http.Request) {

	// To Connect with the Database
	db, err := gorm.Open("postgres", "host=localhost user=postgres password=password dbname=go_patrolling_development sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var user models.User
	// var session models.Session

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		panic(err)
	}

	find_by_mobile_number := db.Where("mobile_number = ?", user.MobileNumber).Find(&user)

	if find_by_mobile_number.RecordNotFound() == true {

		key := []byte("traveling is fun")
		password := []byte(user.Password)
		confirm_password := []byte(user.PasswordConfirmation)

		var user = models.User{FirstName: user.FirstName, LastName: user.LastName, Email: user.Email, MobileNumber: user.MobileNumber, Password: controllers.Encrypt(key, password), PasswordConfirmation: controllers.Encrypt(key, confirm_password), DeviseToken: user.DeviseToken}

		db.Create(&user)

		get_user, err := db.Model(&user).Where("mobile_number = ?", user.MobileNumber).Select("id").Rows()

		defer get_user.Close()

		for get_user.Next() {
			var id int
			err = get_user.Scan(&id)
			if err != nil {
				panic(err)
			}

			var session = models.Session{UserID: id, DeviseToken: user.DeviseToken}

			db.Create(&session)

			var device = models.Device{Token: user.DeviseToken}
			db.Create(&device)
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
