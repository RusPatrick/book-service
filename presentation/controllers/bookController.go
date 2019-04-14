package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ruspatrick/go-toff/application/services"
	"github.com/ruspatrick/go-toff/domain/models"
)

func BookController(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		addBook(w, req)
	case http.MethodGet:

	case http.MethodPatch:
		modifyBook(w, req)
	case http.MethodDelete:
		deleteBook(w, req)
	}

}

func addBook(w http.ResponseWriter, req *http.Request) {
	book := new(models.Book)

	if err := json.NewDecoder(req.Body).Decode(book); err != nil {
		return
	}

	services.AddBook(book)
}

func modifyBook(w http.ResponseWriter, req *http.Request) {
	id, err := getBookID(req.URL.Path)
	if err != nil {
		writeBusinessError(w, errTitle, errIncorrectId.Error())
		return
	}

	book := new(models.Book)
	if err := json.NewDecoder(req.Body).Decode(book); err != nil {
		return
	}
	book.ID = id

	services.ModifyBook(book)

}

func deleteBook(w http.ResponseWriter, req *http.Request) {

}

func getBookID(idStr string) (int, error) {
	return strconv.Atoi(strings.TrimPrefix(idStr, "/books/"))
}
