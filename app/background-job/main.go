package main

import (
	"net/http"
	"os"
	"time"

	"github.com/bukalapak/packen/instrument"
	"github.com/julienschmidt/httprouter"
	"github.com/subosito/gotenv"

	gx "github.com/bukalapak/go-xample"
	"github.com/bukalapak/go-xample/connection"
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

	ecOpt := connection.Option{
		Timeout: 3 * time.Second,
	}

	mysql, _ := database.NewMySQL(dbOpt)
	rmq, _ := messenger.NewRabbitMQ(rmqOpt)
	ec := connection.NewEmailChecker(ecOpt)

	goXample := gx.NewGoXample(mysql, rmq, ec)
	gxHandler := handler.NewHandler(goXample)

	router := httprouter.New()
	router.GET("/healthz", gxHandler.Healthz)
	router.HandlerFunc("GET", "/metrics", instrument.Handler)

	go rmq.Listen(goXample)

	http.ListenAndServe(":1235", router)
}
