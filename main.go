package main

import (
	"net/http"

	"github.com/ruspatrick/go-toff/application/services"
	"github.com/ruspatrick/go-toff/presentation/core/config"
	"github.com/ruspatrick/go-toff/presentation/core/routers"
)

func init() {
	config.Read()
	services.Init()
}

func main() {
	r := routers.NewRouter()
	http.ListenAndServe(config.Get().App.Port, r)
}
