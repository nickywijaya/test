package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	gx "github.com/bukalapak/go-xample"
	"github.com/bukalapak/go-xample/handler"
)

func main() {
	gX := gx.NewGoXample()
	gxHandler := handler.NewHandler(gX)

	router := httprouter.New()
	router.GET("/healthz", gxHandler.Healthz)
	router.GET("/metrics", gxHandler.Metric)
	router.POST("/users", gxHandler.CreateUser)
	http.ListenAndServe(":1234", router)
}
