package main

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"log"

	"github.com/OlegGibadulin/tech-db-forum/config"
	"github.com/OlegGibadulin/tech-db-forum/internal/mwares"
)

func main() {
	config, err := config.LoadConfig("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	// Database
	dbConnection, err := sql.Open("postgres", config.GetDbConnString())
	if err != nil {
		log.Fatal(err)
	}
	defer dbConnection.Close()

	if err := dbConnection.Ping(); err != nil {
		log.Fatal(err)
	}

	// Repository

	// Usecases

	// Middleware
	e := echo.New()
	mw := mwares.NewMiddlewareManager()
	e.Use(mw.PanicRecovering, mw.AccessLog)

	// Delivery

	log.Fatal(e.Start(config.GetServerConnString()))
}
