package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/alSergey/TechMain_2021_db_forum/configs"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/middleware"
	"github.com/alSergey/TechMain_2021_db_forum/internal/app/tools/database"

	forumHandler "github.com/alSergey/TechMain_2021_db_forum/internal/app/forum/delivery/http"
	forumRepo "github.com/alSergey/TechMain_2021_db_forum/internal/app/forum/repository/postgres"
	forumUsecase "github.com/alSergey/TechMain_2021_db_forum/internal/app/forum/usecase"

	postHandler "github.com/alSergey/TechMain_2021_db_forum/internal/app/post/delivery/http"
	postRepo "github.com/alSergey/TechMain_2021_db_forum/internal/app/post/repository/postgres"
	postUsecase "github.com/alSergey/TechMain_2021_db_forum/internal/app/post/usecase"

	serviceHandler "github.com/alSergey/TechMain_2021_db_forum/internal/app/service/delivery/http"
	serviceRepo "github.com/alSergey/TechMain_2021_db_forum/internal/app/service/repository/postgres"
	serviceUsecase "github.com/alSergey/TechMain_2021_db_forum/internal/app/service/usecase"

	threadHandler "github.com/alSergey/TechMain_2021_db_forum/internal/app/thread/delivery/http"
	threadRepo "github.com/alSergey/TechMain_2021_db_forum/internal/app/thread/repository/postgres"
	threadUsecase "github.com/alSergey/TechMain_2021_db_forum/internal/app/thread/usecase"

	userHandler "github.com/alSergey/TechMain_2021_db_forum/internal/app/user/delivery/http"
	userRepo "github.com/alSergey/TechMain_2021_db_forum/internal/app/user/repository/postgres"
	userUsecase "github.com/alSergey/TechMain_2021_db_forum/internal/app/user/usecase"
)

func main() {
	postgresDB, err := database.NewPostgres(configs.Configs.GetPostgresConfig())
	if err != nil {
		log.Fatal(err)
	}

	forumRepo := forumRepo.NewForumRepository(postgresDB.GetDatabase())
	postRepo := postRepo.NewPostRepository(postgresDB.GetDatabase())
	serviceRepo := serviceRepo.NewServiceRepository(postgresDB.GetDatabase())
	threadRepo := threadRepo.NewThreadRepository(postgresDB.GetDatabase())
	userRepo := userRepo.NewUserRepository(postgresDB.GetDatabase())

	forumUsecase := forumUsecase.NewForumUsecase(forumRepo)
	postUsecase := postUsecase.NewForumUsecase(postRepo, threadRepo)
	serviceUsecase := serviceUsecase.NewForumUsecase(serviceRepo)
	threadUsecase := threadUsecase.NewForumUsecase(threadRepo, forumRepo)
	userUsecase := userUsecase.NewForumUsecase(userRepo)

	forumHandler := forumHandler.NewForumHandler(forumUsecase, threadUsecase)
	postHandler := postHandler.NewPostHandler(postUsecase)
	serviceHandler := serviceHandler.NewServiceHandler(serviceUsecase)
	threadHandler := threadHandler.NewThreadHandler(threadUsecase, postUsecase)
	userHandler := userHandler.NewUserHandler(userUsecase)

	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.JSONMiddleware)

	forumHandler.Configure(api)
	postHandler.Configure(api)
	serviceHandler.Configure(api)
	threadHandler.Configure(api)
	userHandler.Configure(api)

	server := http.Server{
		Addr:         fmt.Sprint(":", configs.Configs.GetMainPort()),
		Handler:      router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
