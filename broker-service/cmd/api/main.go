package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "80"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {

	// try to connect to rbbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	// start listeing to messages
	log.Println("Listening for consuming RabbitMq messages...")

	app := &Config{
		Rabbit: rabbitConn,
	}
	log.Printf("Starting broker server on port %s", webPort)

	// http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()

	if err != nil {
		log.Panic("Failed to start broker server", err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backoff time.Duration
	var connection *amqp.Connection

	// dont continue untill rabbit is ready

	for {
		cnn, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Printf("rabit mq not yet ready")
			counts++
		} else {
			log.Println("connected to rabbit mq")

			connection = cnn
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}
		backoff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off")
		time.Sleep(backoff)
		continue

	}
	return connection, nil
}
