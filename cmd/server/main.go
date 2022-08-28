// @title Payment Service
// @version 1.0
// @description The MelonWallet microservice responsible for dealing with payment and crypto transaction information.

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Melon-Network-Inc/account-service/pkg/friend"
	"github.com/Melon-Network-Inc/account-service/pkg/user"

	"github.com/Melon-Network-Inc/common/pkg/config"
	dbcontext "github.com/Melon-Network-Inc/common/pkg/dbcontext"
	"github.com/Melon-Network-Inc/common/pkg/log"

	"github.com/Melon-Network-Inc/payment-service/docs"
	"github.com/Melon-Network-Inc/payment-service/pkg/transaction"

	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// ServiceName indicates the name of current service.
// Version indicates the current version of the application.
const (
	ServiceName       	= "payment-service"
	Version           	= "1.0.0"
	ServiceConfigPath 	= "../config/prod.yml"
	DefaultServicePort 	= 7001
)

var swagHandler gin.HandlerFunc

type Server struct {
	App      *gin.Engine
	Cache    *dbcontext.Cache
	Database *dbcontext.DB
	Logger   log.Logger
}

func init() {
	swagHandler = ginSwagger.WrapHandler(swaggerfiles.Handler)
}

// @title Melon Wallet Service API
// @host localhost:8080
// @description This is backend server for Melon Wallet..
// @version 1.0
// @BasePath /api/v1

// @contact.name API Support
// @contact.url https://melonnetwork.io
// @contact.email support@melonnetwork.io

// @query.collection.format  multi
func main() {
	serverConfig := config.BuildServerConfig(ServiceName, Version, DefaultServicePort, ServiceConfigPath)

	// create root logger tagged with server version
	logger := log.New(serverConfig.ServiceName).With(context.Background(), "version", serverConfig.Version)

	db, err := dbcontext.ConnectToDatabase(serverConfig.DatabaseUrl)
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}

	s := Server{
		App:      gin.Default(),
		Cache:    dbcontext.NewCache(dbcontext.ConnectToCache(serverConfig.CacheUrl)),
		Database: dbcontext.NewDatabase(db),
		Logger:   logger,
	}

	s.buildHandlers()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", serverConfig.ServerPort),
		Handler: s.App,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
		case <-ctx.Done():
			logger.Info("timeout of 5 seconds.")
	}
	logger.Info("Server exiting")
}

func (s *Server) buildHandlers() {
	transactionRepo := transaction.NewRepository(s.Database, s.Logger)
	userRepo := user.NewRepository(s.Database, s.Cache, s.Logger)
	friendRepo := friend.NewRepository(s.Database, s.Logger)

	transactionService := transaction.NewService(transactionRepo, userRepo, friendRepo, s.Logger)

	v1 := s.App.Group("api/v1")
	transaction.RegisterHandlers(v1, transactionService, s.Logger)

	if swagHandler != nil {
		s.buildSwagger()
		s.App.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
}

func (s *Server) buildSwagger() {
	docs.SwaggerInfo.Title = "Melon Wallet Service API"
	docs.SwaggerInfo.Description = "This is backend server for Melon Wallet."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}
