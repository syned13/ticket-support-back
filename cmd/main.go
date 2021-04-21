package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	authHandler "github.com/syned13/ticket-support-back/internal/handlers/auth"
	ticketsHandler "github.com/syned13/ticket-support-back/internal/handlers/tickets"
	ticketsRepository "github.com/syned13/ticket-support-back/internal/repositories/tickets/postgres"
	usersRepo "github.com/syned13/ticket-support-back/internal/repositories/users/postgres"
	authService "github.com/syned13/ticket-support-back/internal/service/auth"
	ticketsService "github.com/syned13/ticket-support-back/internal/service/tickets"
	"github.com/syned13/ticket-support-back/pkg/config"
)

func main() {
	config, err := config.GetConfigFromEnv()
	if err != nil {
		log.Fatal("getting_config_failed: " + err.Error())
	}

	fmt.Println(config.DatabaseConfig.Connection)

	ctx := context.Background()

	var pool *pgxpool.Pool
	isConnected := false

	for i := 0; i < 10; i++ {
		pool, err = pgxpool.Connect(ctx, config.DatabaseConfig.Connection)
		if err != nil {
			time.Sleep(time.Second * 2)
			continue
		} else {
			isConnected = true
			break
		}
	}

	if !isConnected {
		log.Fatal("failed_connecting_to_database: " + err.Error())
	}

	usersRepo, err := usersRepo.New(pool)
	if err != nil {
		log.Fatal("users_repo_initialization_failed")
	}

	authService := authService.New(usersRepo)

	router := mux.NewRouter()

	authHandler.SetupRoutes(ctx, authService, router)

	ticketsRepo, err := ticketsRepository.New(pool)
	if err != nil {
		log.Fatal("tickets_repo_initialization_failed")
	}

	ticketsService := ticketsService.New(ticketsRepo, usersRepo)

	ticketsHandler.SetupRoutes(ctx, ticketsService, router)
	fmt.Printf("Listeting on port :%s\n", config.Port)

	err = http.ListenAndServe(":"+config.Port, router)
	if err != nil {
		log.Fatal("initializing_server_failed: " + err.Error())
	}
}
