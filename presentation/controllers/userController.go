package controllers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ruspatrick/book-service/application/services"
	"github.com/ruspatrick/book-service/domain/models"
)

const (
	maxLengthSalt = 13
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Signup(w http.ResponseWriter, req *http.Request) {
	userInfo := new(models.User)
	if err := json.NewDecoder(req.Body).Decode(userInfo); err != nil {
		writeError(w, err)
		return
	}

	salt := randStringBytes(rand.Int() % 13)

	if err := services.Signup(*userInfo, salt); err != nil {
		writeError(w, err)
		return
	}

	cookie, err := services.Login(*userInfo)
	if err != nil {
		writeError(w, err)
		return
	}

	http.SetCookie(w, cookie)
	writeSuccess(w, http.StatusOK, nil, []byte("Authorized"))
}

func Login(w http.ResponseWriter, req *http.Request) {
	userInfo := new(models.User)
	if err := json.NewDecoder(req.Body).Decode(userInfo); err != nil {
		writeError(w, err)
		return
	}

	cookie, err := services.Login(*userInfo)
	if err != nil {
		writeError(w, err)
		return
	}

	http.SetCookie(w, cookie)
	writeSuccess(w, http.StatusOK, nil, []byte("Authorized"))
}

func UserController(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
	case http.MethodDelete:
		deleteUser(e, req)
	}
}

func deleteUser(w http.ResponseWriter, req *http.Request) {
	userID, err := getUserID(req.URL.Path)
	if err != nil {
		writeError(w, err)
		return
	}
	services.DeleteUser()
}

func getUserID(str string) (int, error) {
	return strconv.Atoi(strings.TrimPrefix(str, "/api/v1/users/"))
}
