package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ruspatrick/book-service/application/services"
	"github.com/ruspatrick/book-service/presentation/core/config"
	"github.com/ruspatrick/book-service/presentation/core/routers"
)

func init() {
	time.Sleep(5 * time.Second)
	config.Read()
	services.Init()
}

func main() {
	r := routers.NewRouter()
	go services.PeriodicalClearSessions()
	log.Fatal(http.ListenAndServe(config.Get().App.Port, r))
}
