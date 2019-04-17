package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ruspatrick/book-service/application/services"
	"github.com/ruspatrick/book-service/domain/models"
)

const (
	errInternalServerErrorMsg = "Внутренняя ошибка сервиса"
)

func BooksController(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		addBook(w, req)
	case http.MethodGet:
		getBooks(w, req)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(nil)
	}
}

func BookController(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		getBookByID(w, req)
	case http.MethodPatch:
		updateBook(w, req)
	case http.MethodDelete:
		deleteBook(w, req)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(nil)
	}
}

func addBook(w http.ResponseWriter, req *http.Request) {
	book := new(models.Book)

	if err := json.NewDecoder(req.Body).Decode(book); err != nil {
		return
	}

	if err := services.AddBook(book); err != nil {
		writeError(w, err)
		return
	}

	responseBody, err := marshalJSON(book)
	if err != nil {
		writeError(w, err)
		return
	}

	headers := http.Header{}
	headers.Add("location", fmt.Sprintf("%s/%d", req.URL.Path, book.ID))
	writeSuccess(w, http.StatusCreated, headers, responseBody)
}

func updateBook(w http.ResponseWriter, req *http.Request) {
	id, err := getBookID(req.URL.Path)
	if err != nil {
		writeError(w, err)
		return
	}

	book := new(models.Book)
	if err := json.NewDecoder(req.Body).Decode(book); err != nil {
		return
	}
	book.ID = id

	if err := services.UpdateBook(book); err != nil {
		writeError(w, err)
		return
	}

	responseBody, err := marshalJSON(book)
	if err != nil {
		writeError(w, err)
		return
	}

	writeSuccess(w, http.StatusOK, nil, responseBody)
}

func deleteBook(w http.ResponseWriter, req *http.Request) {
	bookID, err := getBookID(req.URL.Path)
	if err != nil {
		writeError(w, err)
		return
	}

	if err := services.DeleteBook(bookID); err != nil {
		writeError(w, err)
		return
	}

	writeSuccess(w, http.StatusOK, nil, nil)
}

func getBooks(w http.ResponseWriter, req *http.Request) {
	booksQuery := parseQuery(req.URL.Query())

	books, err := services.GetBooks(booksQuery)
	if err != nil {
		writeError(w, err)
		return
	}

	responseBody, err := marshalJSON(books)
	if err != nil {
		writeError(w, err)
		return
	}
	writeSuccess(w, http.StatusOK, nil, responseBody)
}

func getBookByID(w http.ResponseWriter, req *http.Request) {
	id, err := getBookID(req.URL.Path)
	if err != nil {
		writeError(w, err)
		return
	}
	book, err := services.GetBookByID(id)
	if err != nil {
		writeError(w, err)
		return
	}
	responseBody, err := marshalJSON(book)
	if err != nil {
		writeError(w, err)
		return
	}
	writeSuccess(w, http.StatusOK, nil, responseBody)
}

func getBookID(idStr string) (int, error) {
	return strconv.Atoi(strings.TrimPrefix(idStr, "/api/v1/books/"))
}
