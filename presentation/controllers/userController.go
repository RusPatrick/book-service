package controllers

import (
	"github.com/ruspatrick/book-service/domain/models" 
	"encoding/json"
	"net/http"
)

func Signup(w http.ResponseWriter, req *http.Request) {
	userInfo := new(models.User)
	if err := json.NewDecoder(req.Body).Decode(userInfo); err != nil {
		writeError(w, err)
		return
	}

	services.
}
