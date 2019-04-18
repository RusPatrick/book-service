package routers

import (
	"net/http"

	"github.com/ruspatrick/book-service/application/services"
	"github.com/ruspatrick/book-service/presentation/controllers"
	"github.com/ruspatrick/book-service/presentation/middlewares"
)

const (
	apiV1Url = "/api/v1"
)

func NewRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.Handle(apiV1Url+"/books", middlewares.AuthMiddleware(http.HandlerFunc(controllers.BooksController), services.GetDB()))
	router.Handle(apiV1Url+"/books/", middlewares.AuthMiddleware(http.HandlerFunc(controllers.BookController), services.GetDB()))
	router.HandleFunc(apiV1Url+"/signup", controllers.Signup)
	router.HandleFunc(apiV1Url+"/login", controllers.Login)
	router.Handle(apiV1Url+"/me", middlewares.AuthMiddleware(http.HandlerFunc(controllers.Me), services.GetDB()))
	return router
}
