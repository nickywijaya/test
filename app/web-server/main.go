package main

import (
	"net/http"

	"github.com/bukalapak/go-xample/handler"
	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.HandlerFunc("GET", "/healthz", handler.Healthz)
	router.HandlerFunc("GET", "/metrics", handler.Metric)
	http.ListenAndServe(":1234", router)
}
