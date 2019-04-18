package main

import (
	"net/http"

	"github.com/ruspatrick/book-service/application/services"
	"github.com/ruspatrick/book-service/presentation/core/config"
	"github.com/ruspatrick/book-service/presentation/core/routers"
)

func init() {
	config.Read()
	services.Init()
}

func main() {
	r := routers.NewRouter()
	go services.PeriodicalClearSessions()
	http.ListenAndServe(config.Get().App.Port, r)
}
