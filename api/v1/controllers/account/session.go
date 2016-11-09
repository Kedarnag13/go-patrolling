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

// 	find_user, err := db.Model(&user).Where("mobile_number = ?", session.MobileNumber).Select("id").Rows()
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer find_user.Close()

// 	for find_user.Next() {
// 		flag = 0
// 		var id int
// 		err = find_user.Scan(&id)
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
// 		} else if find_session.RecordNotFound() == false {
// 			b, err := json.Marshal(models.Message{
// 				Success: false,
// 				Message: "Session already exists!",
// 				Error:   "",
// 			})
// 			if err != nil {
// 				panic(err)
// 			}
// 			rw.Header().Set("Content-Type", "application/json")
// 			rw.Write(b)
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
// 		}
// 		goto end
// 	}
// 	if flag == 1 {
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
// 	}
// end:
// }
