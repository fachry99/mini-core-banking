package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/fachry/mini-core-banking/internal/config"
	"github.com/fachry/mini-core-banking/internal/handler"
	"github.com/fachry/mini-core-banking/internal/middleware"
	"github.com/fachry/mini-core-banking/internal/repository"
	"github.com/fachry/mini-core-banking/internal/service"
	"github.com/fachry/mini-core-banking/internal/worker"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := config.NewPostgres(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Worker create channel dan start goroutine
	auditWorker := worker.NewAuditWorker(100)
	auditWorker.Start(ctx)

	// Repository
	userRepo := repository.NewUserRepository(db)
	accountRepo := repository.NewAccountRepository(db)
	idempotencyRepo := repository.NewIdempotencyRepository(db)

	// Service
	transferService := service.NewTransferService(
		db,
		auditWorker,
	)

	// Handler
	userHandler := handler.NewUserHandler(userRepo)
	accountHandler := handler.NewAccountHandler(accountRepo)
	depositHandler := handler.NewDepositHandler(accountRepo)
	transferHandler := handler.NewTransferHandler(
		transferService,
		idempotencyRepo,
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/users", userHandler.CreateUser)
	mux.HandleFunc("/accounts", accountHandler.CreateAccount)
	mux.HandleFunc("/deposit", depositHandler.Deposit)
	mux.HandleFunc("/transfer", transferHandler.Transfer)

	handlerWithMiddleware := middleware.CORSMiddleware(
		middleware.RequestID(mux),
	)

	// Graceful shutdown
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		cancel()
	}()

	log.Println(cfg.AppName, "running on port", cfg.AppPort)
	log.Fatal(http.ListenAndServe(":"+cfg.AppPort, handlerWithMiddleware))
}
