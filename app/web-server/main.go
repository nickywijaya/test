package main

import (
	"net/http"

	"github.com/bukalapak/packen/middleware"
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

	router.POST("/users", middleware.MonitorHTTP("create-user", gxHandler.CreateUser))
	router.GET("/users/:id", middleware.MonitorHTTP("get-user", gxHandler.GetUser))

	router.POST("/login", middleware.MonitorHTTP("login", gxHandler.Login))

	http.ListenAndServe(":1234", router)
}
