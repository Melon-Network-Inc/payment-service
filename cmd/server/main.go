// @title Payment Service
// @version 1.0
// @description The MelonWallet microservice responsible for dealing with payment and crypto transaction information.

package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	dbcontext "github.com/Melon-Network-Inc/payment-service/pkg/dbcontext"
	"github.com/Melon-Network-Inc/payment-service/pkg/log"
	transaction "github.com/Melon-Network-Inc/payment-service/pkg/transaction"
	"github.com/gorilla/mux"
)

// Version indicates the current version of the application.
var Version = "1.0.0"

func main() {

	// create root logger tagged with server version
	logger := log.New().With(context.Background(), "version", Version)
	logger.Info("Payment Service Started")

	var wait time.Duration
	flag.DurationVar(
		&wait,
		"graceful-timeout",
		time.Second*15,
		"the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m",
	)
	flag.Parse()

	db, err := dbcontext.Connect("postgres://postgres:123456@localhost:5432/melon_service")
	if err != nil {
		logger.Error(err)
	}
	r := mux.NewRouter()

	transactionRepo := transaction.NewRepository(dbcontext.New(db), logger)
	transaction.RegisterHandlers(r, transaction.NewService(transactionRepo, logger), logger)

	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Info(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	logger.Info("shutting down")
	os.Exit(0)
}
