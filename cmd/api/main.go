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
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	// DB connection
	db, err := config.NewPostgres(cfg.DB)
	if err != nil {
		log.Fatal("failed to connect db:", err)
	}
	defer db.Close()

	// Repository
	userRepo := repository.NewUserRepository(db)
	accountRepo := repository.NewAccountRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	// Service
	transferService := service.NewTransferService(db, accountRepo, transactionRepo)

	// Handler
	userHandler := handler.NewUserHandler(userRepo)
	accountHandler := handler.NewAccountHandler(accountRepo)
	transferHandler := handler.NewTransferHandler(transferService)

	// Routes
	http.HandleFunc("/users", userHandler.CreateUser)
	http.HandleFunc("/accounts", accountHandler.CreateAccount)
	http.HandleFunc("/transfer", transferHandler.Transfer)

	log.Println(cfg.AppName, "running on port", cfg.AppPort)
	log.Fatal(http.ListenAndServe(":"+cfg.AppPort, nil))
}
