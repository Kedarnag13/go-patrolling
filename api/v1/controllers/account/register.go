package account

import (
	"encoding/json"
	"github.com/Kedarnag13/go-patrolling/api/v1/models"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"log"
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

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		panic(err)
	}

	find_by_mobile_number := db.Where("mobile_number = ?", user.MobileNumber).First(&user)

	if find_by_mobile_number.RecordNotFound() == true {
		db.Create(&user)
	} else {
		log.Println(find_by_mobile_number.GetErrors())
	}

}
