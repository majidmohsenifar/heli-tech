package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	transactionpb "github.com/majidmohsenifar/heli-tech/data-contracts/proto/transaction"
	userpb "github.com/majidmohsenifar/heli-tech/data-contracts/proto/user"
	"github.com/majidmohsenifar/heli-tech/gateway-service/config"
	"github.com/majidmohsenifar/heli-tech/gateway-service/handler/api"
	"github.com/majidmohsenifar/heli-tech/gateway-service/handler/api/router"
	"github.com/majidmohsenifar/heli-tech/gateway-service/logger"
	"github.com/majidmohsenifar/heli-tech/gateway-service/service/transaction"
	"github.com/majidmohsenifar/heli-tech/gateway-service/service/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	Address = "0.0.0.0:8080"
)

func main() {
	logger := logger.NewLogger()
	viper := config.NewViper()
	gin.SetMode(gin.ReleaseMode)
	userConn, err := grpc.NewClient(
		viper.GetString("usersrv.address"),
		//config.UserServiceUrl(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Error("Failed to dial user service", err)
		os.Exit(1)
	}
	transactionConn, err := grpc.NewClient(
		viper.GetString("transactionsrv.address"),
		//config.UserServiceUrl(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Error("Failed to dial transaction service", err)
		os.Exit(1)
	}
	userClient := userpb.NewUserClient(userConn)
	transactionClient := transactionpb.NewTransactionClient(transactionConn)

	userService := user.NewService(userClient, logger)
	transactionService := transaction.NewService(transactionClient, logger)
	validator := validator.New()
	userHandler := api.NewUserHandler(userService, validator)
	transactionHandler := api.NewTransactionHandler(transactionService, validator)
	api.InitiateSwagger()
	router := router.New(
		userHandler,
		transactionHandler,
		userService,
		logger,
	)
	if err != nil {
		logger.Error("failed to initialize router", err)
		os.Exit(1)
	}
	httpServer := &http.Server{
		Addr:    Address,
		Handler: router.Handler,
	}
	err = httpServer.ListenAndServe()
	if err != nil {
		logger.Error("cannot listen and serv", err)
		os.Exit(1)
	}
}
