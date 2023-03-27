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

	"github.com/Melon-Network-Inc/common/pkg/blockchain"
	"github.com/Melon-Network-Inc/payment-service/pkg/taskq"

	gcs "cloud.google.com/go/storage"
	accountRepo "github.com/Melon-Network-Inc/account-service/pkg/repository"
	"github.com/Melon-Network-Inc/common/pkg/config"
	"github.com/Melon-Network-Inc/common/pkg/dbcontext"
	"github.com/Melon-Network-Inc/common/pkg/log"
	message "github.com/Melon-Network-Inc/common/pkg/notification"
	"github.com/Melon-Network-Inc/common/pkg/storage"
	"github.com/Melon-Network-Inc/common/pkg/utils"
	"github.com/Melon-Network-Inc/payment-service/docs"
	"github.com/Melon-Network-Inc/payment-service/pkg/activity"
	"github.com/Melon-Network-Inc/payment-service/pkg/news"
	"github.com/Melon-Network-Inc/payment-service/pkg/repository"
	"github.com/Melon-Network-Inc/payment-service/pkg/transaction"
	paymentUtils "github.com/Melon-Network-Inc/payment-service/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var swagHandler gin.HandlerFunc

type Server struct {
	App           *gin.Engine
	Cache         *dbcontext.Cache
	Database      *dbcontext.DB
	Cronjob       *gocron.Scheduler
	QueueManager  taskq.QueueManager
	BlockClient   blockchain.BlockDaemonClient
	StorageClient *storage.StorageClient
	FcmClient     *message.FCMClient
	Logger        log.Logger
}

func init() {
	swagHandler = ginSwagger.WrapHandler(swaggerFiles.Handler)
}

// @title Melon Wallet Service API
// @host localhost:8080
// @description This is backend server for Melon Wallet.
// @version 1.0
// @BasePath /api/v1

// @contact.name API Support
// @contact.url https://melonnetwork.io
// @contact.email support@melonnetwork.io

// @query.collection.format  multi
func main() {
	serverConfig := config.BuildServerConfig("config/payment.yml")
	blockchainConfig := config.BuildBlockchainConfig("config/blockchain.yml")

	// create root logger tagged with server version
	logger := log.New(serverConfig.ServiceName).Default(context.Background(), serverConfig, "version", serverConfig.Version)

	db, err := dbcontext.ConnectToDatabase(serverConfig.DatabaseUrl)
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}

	router := gin.Default()
	router.Use(log.GinLogger(logger), log.GinRecovery(logger, true))

	serverLocation, err := paymentUtils.GetPSTLocation()
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}

	gcsClient, err := gcs.NewClient(context.Background())
	if err != nil {
		logger.Errorf("error initializing Google Cloud Storage: %v\n", err)
		os.Exit(-1)
	}
	storageClient := storage.NewStorageClient(gcsClient)

	fcmClient, err := message.NewNotificationClient(context.Background(), logger)
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}

	blockClient := blockchain.NewBlockDaemonClient(blockchainConfig.UbiAccessToken, blockchainConfig.UbiEndpoint, logger)

	s := Server{
		App:           router,
		Cache:         dbcontext.NewCache(dbcontext.ConnectToCache(serverConfig.CacheUrl), logger),
		Database:      dbcontext.NewDatabase(db),
		Cronjob:       gocron.NewScheduler(serverLocation),
		QueueManager:  taskq.NewTaskQueueManager(serverConfig),
		BlockClient:   blockClient,
		StorageClient: &storageClient,
		FcmClient:     fcmClient,
		Logger:        logger,
	}

	// Bind all handlers to wallet server
	s.buildHandlers()

	if !utils.IsProdEnvironment() {
		logger.Debug(router.Run(fmt.Sprintf(":%d", serverConfig.ServerPort)))
	} else {
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
}

func (s *Server) buildHandlers() {
	transactionRepo := repository.NewTransactionRepository(s.Database, s.Logger)
	newsRepo := repository.NewNewsRepository(s.Database, s.Logger)

	userRepo := accountRepo.NewUserRepository(s.Database, s.Cache, s.StorageClient, s.Logger)
	friendRepo := accountRepo.NewFriendRepository(s.Database, s.Logger)
	notificationRepo := accountRepo.NewNotificationRepository(s.Database, s.Logger)
	deviceRepo := accountRepo.NewDeviceRepository(s.Database, s.Logger)

	newsClient := news.NewClient(s.Logger)

	transactionService := transaction.NewService(
		transactionRepo,
		userRepo,
		friendRepo,
		deviceRepo,
		notificationRepo,
		s.QueueManager,
		s.BlockClient,
		s.FcmClient,
		s.Logger)
	activityService := activity.NewService(userRepo, transactionRepo, friendRepo, s.Logger)
	newsService := news.NewService(newsRepo, newsClient, s.Logger)
	taskqService := taskq.NewService(s.QueueManager, s.Logger)

	newsConsumer := news.NewConsumer(newsService, s.Logger)
	transactionConsumer := transaction.NewConsumer(transactionService, s.Logger)

	v1 := s.App.Group("api/v1")
	transaction.RegisterHandlers(v1, transactionService, s.Logger)
	activity.RegisterHandler(v1, activityService, s.Logger)
	news.RegisterHandler(v1, newsService, s.Logger)
	taskq.RegisterHandler(v1, taskqService, s.Logger)

	if !utils.IsProdEnvironment() && swagHandler != nil {
		s.buildSwagger()
		s.App.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	s.setupCronJob(newsConsumer, transactionConsumer)
}

func (s *Server) setupCronJob(
	newsConsumer news.Consumer, 
	transactionConsumer transaction.Consumer) {
	var err error
	_, err = s.Cronjob.Every(1).Day().At("8:00").Do(newsConsumer.Collect)
	if err != nil {
		s.Logger.Error("cannot schedule new cron job due to ", err)
	}

	s.Cronjob.StartAsync()
}

func (s *Server) buildSwagger() {
	docs.SwaggerInfo.Title = "Melon Wallet Service API"
	docs.SwaggerInfo.Description = "This is backend server for Melon Wallet."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}
