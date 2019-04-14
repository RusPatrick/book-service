package routers

import (
	"net/http"

	"github.com/ruspatrick/go-toff/presentation/controllers"
)

func NewRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/books", controllers.BookController)
	return router
}
