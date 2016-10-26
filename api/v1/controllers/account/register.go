package account

import (
	// "log"
	"net/http"
)

type RegisterController struct{}

var Registration RegisterController

func (r *RegisterController) Create(rw http.ResponseWriter, req *http.Request) {

}
