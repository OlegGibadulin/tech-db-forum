package main

import (
	"database/sql"
	"log"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"github.com/OlegGibadulin/tech-db-forum/config"
	"github.com/OlegGibadulin/tech-db-forum/internal/mwares"

	userHandler "github.com/OlegGibadulin/tech-db-forum/internal/user/delivery"
	userRepo "github.com/OlegGibadulin/tech-db-forum/internal/user/repository"
	userUsecase "github.com/OlegGibadulin/tech-db-forum/internal/user/usecases"

	forumHandler "github.com/OlegGibadulin/tech-db-forum/internal/forum/delivery"
	forumRepo "github.com/OlegGibadulin/tech-db-forum/internal/forum/repository"
	forumUsecase "github.com/OlegGibadulin/tech-db-forum/internal/forum/usecases"

	threadHandler "github.com/OlegGibadulin/tech-db-forum/internal/thread/delivery"
	threadRepo "github.com/OlegGibadulin/tech-db-forum/internal/thread/repository"
	threadUsecase "github.com/OlegGibadulin/tech-db-forum/internal/thread/usecases"
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
	userRepo := userRepo.NewUserPgRepository(dbConnection)
	forumRepo := forumRepo.NewForumPgRepository(dbConnection)
	threadRepo := threadRepo.NewThreadPgRepository(dbConnection)

	// Usecases
	userUcase := userUsecase.NewUserUsecase(userRepo)
	forumUcase := forumUsecase.NewForumUsecase(forumRepo)
	threadUcase := threadUsecase.NewThreadUsecase(threadRepo)

	// Middleware
	e := echo.New()
	mw := mwares.NewMiddlewareManager()
	e.Use(mw.PanicRecovering, mw.AccessLog)

	// Delivery
	userHandler := userHandler.NewUserHandler(userUcase)
	forumHandler := forumHandler.NewForumHandler(forumUcase, userUcase)
	threadHandler := threadHandler.NewThreadHandler(threadUcase, userUcase, forumUcase)

	userHandler.Configure(e, mw)
	forumHandler.Configure(e, mw)
	threadHandler.Configure(e, mw)

	log.Fatal(e.Start(config.GetServerConnString()))
}
