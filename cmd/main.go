package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	authHandler "github.com/syned13/ticket-support-back/internal/handlers/auth"
	usersRepo "github.com/syned13/ticket-support-back/internal/repositories/users/postgres"
	authService "github.com/syned13/ticket-support-back/internal/service/auth"
	"github.com/syned13/ticket-support-back/pkg/config"
)

func main() {
	config, err := config.GetConfigFromEnv()
	if err != nil {
		log.Fatal("getting_config_failed: " + err.Error())
	}

	ctx := context.Background()

	pool, err := pgxpool.Connect(ctx, config.DatabaseConfig.Connection)
	if err != nil {
		log.Fatal("failed_connecting_to_database")
	}

	usersRepo, err := usersRepo.New(pool)
	if err != nil {
		log.Fatal("users_repo_initialization_failed")
	}

	authService := authService.New(usersRepo)

	router := mux.NewRouter()

	authHandler.SetupRoutes(ctx, authService, router)

	fmt.Println("hello world!")
}
