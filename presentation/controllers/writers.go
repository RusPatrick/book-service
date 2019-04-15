package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ruspatrick/book-service/infrastructure/errors"
)

const (
	errTitle          = "Ошибка"
	serverTypeError   = "server"
	businessTypeError = "business"
)

func writeError(w http.ResponseWriter, err error) {
	switch v := err.(type) {
	case errors.Error:
		log.Printf("%s: %s", v.Source, v.ErrDescription.Error())

		response, errMarshal := json.Marshal(err)
		if errMarshal != nil {
			log.Println(errMarshal)
		}

		w.WriteHeader(v.StatusCode)
		w.Write(response)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
	}
}

func writeSuccess(w http.ResponseWriter, statusCode int, headers http.Header, body []byte) {
	for k, v := range headers {
		w.Header().Add(k, v[0])
	}
	w.WriteHeader(statusCode)
	w.Write(body)
}
