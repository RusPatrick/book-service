package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/ruspatrick/go-toff/domain/models"
)

const (
	errTitle          = "Ошибка"
	serverTypeError   = "server"
	businessTypeError = "business"
)

var (
	errIncorrectId = errors.New("Incorrect book id")
)

func writeServerError(w http.ResponseWriter, title, detail string) {
	w.WriteHeader(http.StatusInternalServerError)
	writeError(w, title, detail, serverTypeError)
}

func writeBusinessError(w http.ResponseWriter, title, detail string) {
	w.WriteHeader(http.StatusBadRequest)
	writeError(w, title, detail, businessTypeError)
}

func writeError(w http.ResponseWriter, title, detail, typeErr string) {
	response, err := json.Marshal(models.Error{
		Title:  title,
		Type:   typeErr,
		Detail: detail,
	})
	if err != nil {
		log.Println(err)
	}

	w.Write(response)
}
