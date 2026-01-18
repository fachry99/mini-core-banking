package main

import (
	"log"
	"net/http"

	"github.com/fachry/mini-core-banking/internal/config"
	"github.com/fachry/mini-core-banking/internal/handler"
	"github.com/fachry/mini-core-banking/internal/repository"
	"github.com/fachry/mini-core-banking/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := config.NewPostgres(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Repository
	userRepo := repository.NewUserRepository(db)
	accountRepo := repository.NewAccountRepository(db)

	// Service
	transferService := service.NewTransferService(db)

	// Handler
	userHandler := handler.NewUserHandler(userRepo)
	accountHandler := handler.NewAccountHandler(accountRepo)
	depositHandler := handler.NewDepositHandler(accountRepo)
	transferHandler := handler.NewTransferHandler(transferService)

	// Routes
	http.HandleFunc("/users", userHandler.CreateUser)
	http.HandleFunc("/accounts", accountHandler.CreateAccount)
	http.HandleFunc("/deposit", depositHandler.Deposit)
	http.HandleFunc("/transfer", transferHandler.Transfer)

	log.Println(cfg.AppName, "running on port", cfg.AppPort)
	log.Fatal(http.ListenAndServe(":"+cfg.AppPort, nil))
}
