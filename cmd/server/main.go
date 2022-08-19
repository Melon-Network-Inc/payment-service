// @title Payment Service
// @version 1.0
// @description The MelonWallet microservice responsible for dealing with payment and crypto transaction information.

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Melon-Network-Inc/account-service/pkg/friend"
	"github.com/Melon-Network-Inc/account-service/pkg/user"

	dbcontext "github.com/Melon-Network-Inc/common/pkg/dbcontext"
	"github.com/Melon-Network-Inc/common/pkg/log"

	"github.com/Melon-Network-Inc/payment-service/config"
	"github.com/Melon-Network-Inc/payment-service/docs"
	"github.com/Melon-Network-Inc/payment-service/pkg/transaction"

	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Version indicates the current version of the application.
var ServiceName = "payment-service"
var Version = "1.0.0"
var swagHandler gin.HandlerFunc

func init() {
	swagHandler = ginSwagger.WrapHandler(swaggerfiles.Handler)
}

type Server struct {
	App  			*gin.Engine
	Cache  			*dbcontext.Cache
	Database 		*dbcontext.DB
	Logger          log.Logger
}

func main() {
	config := config.BuildServerConfig(ServiceName, Version)

	// create root logger tagged with server version
	logger := log.New(config.ServiceName).With(context.Background(), "version", config.Version)
	
	db, err := dbcontext.ConnectToDatabase(config.DatabaseUrl)
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}
	
	s := Server{
		App : gin.Default(), 
		Cache: dbcontext.NewCache(dbcontext.ConnectToCache(config.CacheUrl)), 
		Database: dbcontext.NewDatabase(db), 
		Logger: logger,
	}

	s.buildHandlers()
	s.App.Run(fmt.Sprintf(":%d", config.ServerPort))
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

func (s Server) buildSwagger() {
	docs.SwaggerInfo.Title = "Payment Service API"
	docs.SwaggerInfo.Description = "This is payment server for Melon Wallet."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:7000"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}

