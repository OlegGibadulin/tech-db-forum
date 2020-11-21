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

	threadHandler "github.com/OlegGibadulin/tech-db-forum/internal/thread/delivery"
	threadRepo "github.com/OlegGibadulin/tech-db-forum/internal/thread/repository"
	threadUsecase "github.com/OlegGibadulin/tech-db-forum/internal/thread/usecases"

	forumHandler "github.com/OlegGibadulin/tech-db-forum/internal/forum/delivery"
	forumRepo "github.com/OlegGibadulin/tech-db-forum/internal/forum/repository"
	forumUsecase "github.com/OlegGibadulin/tech-db-forum/internal/forum/usecases"

	postHandler "github.com/OlegGibadulin/tech-db-forum/internal/post/delivery"
	postRepo "github.com/OlegGibadulin/tech-db-forum/internal/post/repository"
	postUsecase "github.com/OlegGibadulin/tech-db-forum/internal/post/usecases"
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
	threadRepo := threadRepo.NewThreadPgRepository(dbConnection)
	forumRepo := forumRepo.NewForumPgRepository(dbConnection)
	postRepo := postRepo.NewPostPgRepository(dbConnection)

	// Usecases
	userUcase := userUsecase.NewUserUsecase(userRepo)
	threadUcase := threadUsecase.NewThreadUsecase(threadRepo)
	forumUcase := forumUsecase.NewForumUsecase(forumRepo)
	postUcase := postUsecase.NewPostUsecase(postRepo)

	// Middleware
	e := echo.New()
	mw := mwares.NewMiddlewareManager()
	e.Use(mw.PanicRecovering, mw.AccessLog)

	// Delivery
	userHandler := userHandler.NewUserHandler(userUcase)
	threadHandler := threadHandler.NewThreadHandler(threadUcase, userUcase, postUcase)
	forumHandler := forumHandler.NewForumHandler(forumUcase, userUcase, threadUcase)
	postHandler := postHandler.NewPostHandler(postUcase, userUcase, threadUcase, forumUcase)

	userHandler.Configure(e, mw)
	threadHandler.Configure(e, mw)
	forumHandler.Configure(e, mw)
	postHandler.Configure(e, mw)

	log.Fatal(e.Start(config.GetServerConnString()))
}
