package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ruspatrick/go-toff/application/services"
	"github.com/ruspatrick/go-toff/domain/models"
)

func BookController(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		addBook(w, req)
	}
}

func addBook(w http.ResponseWriter, req *http.Request) {
	book := new(models.Book)

	if err := json.NewDecoder(req.Body).Decode(book); err != nil {
		return
	}

	services.AddBook(book)
}
