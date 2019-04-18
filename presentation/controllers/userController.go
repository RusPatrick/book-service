package controllers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ruspatrick/book-service/application/services"
	"github.com/ruspatrick/book-service/domain/models"
	"github.com/ruspatrick/book-service/infrastructure/errors"
)

const (
	maxLengthSalt = 13
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var (
	ErrDifferentPasswords = fmt.Errorf("Неверный email или пароль")
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

	if userInfo.Password != userInfo.ConfirmPassword {
		writeError(w, errors.CreateBusinessError(ErrDifferentPasswords, ErrDifferentPasswords.Error()))
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

func Me(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodDelete:
		deleteUser(w, req)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(nil)
	}
}

func deleteUser(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("session_id")
	if err != nil {
		writeError(w, err)
		return
	}
	err = services.DeleteUser(cookie.Value)
	if err != nil {
		writeError(w, err)
		return
	}
	writeSuccess(w, http.StatusOK, nil, nil)
}
