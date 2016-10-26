package account

import (
	"fmt"
	"net/http"
)

type RegisterController struct{}

var Registration RegisterController

func (r *RegisterController) Create(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("I am here")
}
