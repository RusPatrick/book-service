package routers

import (
	"net/http"

	"github.com/ruspatrick/go-toff/presentation/controllers"
)

const (
	apiV1Url = "/api/v1"
)

func NewRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc(apiV1Url+"/books", controllers.BookController)
	return router
}
