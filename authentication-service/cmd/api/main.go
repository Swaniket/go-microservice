package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var dbConnRetryCount int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func openDBConn(dsn string) (*sql.DB, error) {
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

	// Infinite for loop till we connect to the dabatabse with backoff mechanism (May want to implement exponential backoff)
	for {
		connection, err := openDBConn(dsn)
		if err != nil {
			log.Println("err", err)
			log.Println("Postgrest not yet ready...")
			dbConnRetryCount++
		} else {
			log.Println("Connected to postgres!")
			return connection
		}

		if dbConnRetryCount > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for 2 seconds...")
		time.Sleep(2 * time.Second)
		continue
	}

}

func main() {
	log.Println("Starting authetication micro-service")

	// @TODO: Connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}

	// Setup config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
