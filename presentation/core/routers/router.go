package routers

import (
	"net/http"

	"github.com/ruspatrick/book-service/presentation/controllers"
)

const (
	apiV1Url = "/api/v1"
)

func NewRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc(apiV1Url+"/books", controllers.BooksController)
	router.HandleFunc(apiV1Url+"/books/", controllers.BookController)
	router.HandleFunc(apiV1Url+"/signup", controllers.Signup)
	router.HandleFunc(apiV1Url+"/login", controllers.Login)
	return router
}
