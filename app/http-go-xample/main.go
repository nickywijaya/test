package main

import (
	"net/http"
	"os"

	"github.com/bukalapak/packen/middleware"
	"github.com/julienschmidt/httprouter"
	"github.com/subosito/gotenv"

	gx "github.com/bukalapak/go-xample"
	"github.com/bukalapak/go-xample/database"
	"github.com/bukalapak/go-xample/handler"
	"github.com/bukalapak/go-xample/messenger"
)

func main() {
	gotenv.Load("../../.env")

	dbOpt := database.Option{
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     os.Getenv("MYSQL_PORT"),
		Database: os.Getenv("MYSQL_DATABASE"),
		Charset:  os.Getenv("MYSQL_CHARSET"),
	}

	rmqOpt := messenger.RabbitMQOption{
		Username:     os.Getenv("RABBITMQ_USER"),
		Password:     os.Getenv("RABBITMQ_PASSWORD"),
		Host:         os.Getenv("RABBITMQ_HOST"),
		VHost:        os.Getenv("RABBITMQ_VHOST"),
		ExchangeName: os.Getenv("RABBITMQ_EXCHANGE_NAME"),
		ExchangeType: os.Getenv("RABBITMQ_EXCHANGE_TYPE"),
		RoutingKey:   os.Getenv("RABBITMQ_ROUTING_KEY"),
		Durable:      true,
		Exclusive:    false,
	}

	mysql, _ := database.NewMySQL(dbOpt)
	rmq, _ := messenger.NewRabbitMQ(rmqOpt)
	gX := gx.NewGoXample(mysql, rmq)
	gxHandler := handler.NewHandler(gX)

	router := httprouter.New()
	router.GET("/healthz", gxHandler.Healthz)
	router.GET("/metrics", gxHandler.Metric)

	router.POST("/users", middleware.MonitorHTTP("create-user", gxHandler.CreateUser))
	router.GET("/users/:id", middleware.MonitorHTTP("get-user", gxHandler.GetUser))

	router.POST("/login", middleware.MonitorHTTP("login", gxHandler.Login))
	router.GET("/logout", middleware.MonitorHTTP("logout", gxHandler.Logout))

	http.ListenAndServe(":1234", router)
}
