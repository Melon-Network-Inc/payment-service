// @title Payment Service
// @version 1.0
// @description The MelonWallet microservice responsible for dealing with payment and crypto transaction information.

package main

import (
	"context"
	"os"

	"github.com/Melon-Network-Inc/account-service/pkg/friend"
	"github.com/Melon-Network-Inc/account-service/pkg/user"

	dbcontext "github.com/Melon-Network-Inc/common/pkg/dbcontext"
	"github.com/Melon-Network-Inc/common/pkg/log"

	"github.com/Melon-Network-Inc/payment-service/docs"
	"github.com/Melon-Network-Inc/payment-service/pkg/transaction"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Version indicates the current version of the application.
var Version = "1.0.0"
var swagHandler gin.HandlerFunc

func init() {
	swagHandler = ginSwagger.WrapHandler(swaggerfiles.Handler)
}

func main() {
	// create root logger tagged with server version
	logger := log.New().With(context.Background(), "version", Version)
	logger.Info("Payment Service Started")

	viper.SetConfigFile("./pkg/envs/.env")

	var port, dbUrl, redisUrl string
	if err := viper.ReadInConfig(); err == nil {
		port = viper.Get("DB_PORT").(string)
		dbUrl = viper.Get("DB_URL").(string)
		redisUrl = viper.Get("CACHE_URL").(string)
	} else {
		port = ":7000"
		dbUrl = "postgres://postgres:postgres@localhost:5432/melon_service"
		redisUrl = "localhost:6379"
	}
	
	router := gin.Default()
	db, err := dbcontext.ConnectToDatabase(dbUrl)
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}
	
	cache := dbcontext.ConnectToCache(redisUrl)

	buildHandlers(router.Group(""), dbcontext.NewDatabase(db), dbcontext.NewCache(cache), logger)

	router.Run(port)
}

func buildHandlers(router *gin.RouterGroup, db *dbcontext.DB, cache *dbcontext.Cache, logger log.Logger) {
	transactionrRepo := transaction.NewRepository(db, logger)
	userRepo 	     := user.NewRepository(db, cache, logger)
	friendRepo  	 := friend.NewRepository(db, logger)

	transactionService := transaction.NewService(transactionrRepo, userRepo, friendRepo, logger)

	v1 := router.Group("api/v1")
	transaction.RegisterHandlers(v1, transactionService, logger)

	if swagHandler != nil {
		buildSwagger()
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
}

func buildSwagger() {
	docs.SwaggerInfo.Title = "Payment Service API"
	docs.SwaggerInfo.Description = "This is payment server for Melon Wallet."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:7000"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}

