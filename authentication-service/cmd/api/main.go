package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/authentication-service/data"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var count int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Printf("starting authentication service on port %s...", webPort)
	app := Config{}
	// http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	//connect to db
	conn := connectToDB()
	if conn == nil {
		log.Panic("cannot connect to postgres")
	}

	app.DB = conn
	app.Models = data.New(conn)

	err := srv.ListenAndServe()

	if err != nil {
		log.Panic("Failed to start authentication server", err)
	}

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err

	}
	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("postgress not yet ready ....")
			count++
		} else {
			log.Println("connected to postgress...")
			return connection
		}

		if count > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for 2 seconds...")

		time.Sleep(2 * time.Second)
		continue
	}
}
