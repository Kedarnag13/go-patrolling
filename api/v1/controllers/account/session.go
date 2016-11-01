package account

import (
	"encoding/json"
	"github.com/Kedarnag13/go-patrolling/api/v1/models"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type SessionController struct{}

var Session SessionController

func (s SessionController) Create(rw http.ResponseWriter, req *http.Request) {

	flag := 1

	// To Connect with the Database
	db, err := gorm.Open("postgres", "host=localhost user=postgres password=password dbname=go_patrolling_development sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

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

	find_user, err := db.Model(&user).Where("devise_token = ?", session.DeviseToken).Select("id").Rows()
	if err != nil {
		panic(err)
	}

	defer find_user.Close()

	for find_user.Next() {
		flag = 0
		var id int
		err = find_user.Scan(&id)
		find_session := db.Find(&session, "user_id = ?", id)
		if find_session.RecordNotFound() == true {
			b, err := json.Marshal(models.Message{
				Success: false,
				Message: "",
				Error:   "Session does not exist!",
			})
			if err != nil {
				panic(err)
			}
			rw.Header().Set("Content-Type", "application/json")
			rw.Write(b)
		} else if find_session.RecordNotFound() == false {
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
		} else {
			session = models.Session{UserID: id, DeviseToken: session.DeviseToken}
			db.Create(session)
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
		}
		goto end
	}
	if flag == 1 {
		b, err := json.Marshal(models.Message{
			Success: false,
			Message: "",
			Error:   "User does not exist!",
		})
		if err != nil {
			panic(err)
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(b)
		goto end
	}

	// find_user := db.Find(&user, "mobile_number = ? AND devise_token = ?", session.User.MobileNumber, session.DeviseToken)

	// 	if find_user.RecordNotFound() == true {
	// 		b, err := json.Marshal(models.Message{
	// 			Success: false,
	// 			Message: "",
	// 			Error:   "User does not exist!",
	// 		})
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		rw.Header().Set("Content-Type", "application/json")
	// 		rw.Write(b)
	// 		goto end
	// 	} else {
	// 		var id int
	// 		err = find_user.Scan(&id)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		find_session := db.Find(&session, "user_id = ?", id)
	// 		if find_session.RecordNotFound() == true {
	// 			b, err := json.Marshal(models.Message{
	// 				Success: false,
	// 				Message: "",
	// 				Error:   "Session does not exist!",
	// 			})
	// 			if err != nil {
	// 				panic(err)
	// 			}
	// 			rw.Header().Set("Content-Type", "application/json")
	// 			rw.Write(b)
	// 			goto end
	// 		} else {
	// 			session = models.Session{UserID: id, DeviseToken: session.DeviseToken}
	// 			db.Create(session)
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
	// 		}
	// 	}
end:
}
