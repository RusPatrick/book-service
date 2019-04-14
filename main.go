package main

import (
	"net/http"

	"github.com/ruspatrick/go-toff/presentation/core/routers"
)

func main() {
	r := routers.NewRouter()
	http.ListenAndServe(":8000", r)
}
