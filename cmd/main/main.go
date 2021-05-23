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
	forumUsecase := forumUsecase.NewForumUsecase(forumRepo)
	forumHandler := forumHandler.NewForumHandler(forumUsecase)

	postRepo := postRepo.NewPostRepository(postgresDB.GetDatabase())
	postUsecase := postUsecase.NewForumUsecase(postRepo)
	postHandler := postHandler.NewPostHandler(postUsecase)

	serviceRepo := serviceRepo.NewServiceRepository(postgresDB.GetDatabase())
	serviceUsecase := serviceUsecase.NewForumUsecase(serviceRepo)
	serviceHandler := serviceHandler.NewServiceHandler(serviceUsecase)

	threadRepo := threadRepo.NewThreadRepository(postgresDB.GetDatabase())
	threadUsecase := threadUsecase.NewForumUsecase(threadRepo)
	threadHandler := threadHandler.NewThreadHandler(threadUsecase)

	userRepo := userRepo.NewUserRepository(postgresDB.GetDatabase())
	userUsecase := userUsecase.NewForumUsecase(userRepo)
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
