package main

import (
	"net/http"
	"os"

	userpb "github.com/majidmohsenifar/heli-tech/data-contracts/proto/user"
	"github.com/majidmohsenifar/heli-tech/gateway-service/config"
	"github.com/majidmohsenifar/heli-tech/gateway-service/handler/api"
	"github.com/majidmohsenifar/heli-tech/gateway-service/handler/api/router"
	"github.com/majidmohsenifar/heli-tech/gateway-service/logger"
	"github.com/majidmohsenifar/heli-tech/gateway-service/service/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	Address = "localhost:8080"
)

func main() {
	logger := logger.NewLogger()
	viper := config.NewViper()

	userConn, err := grpc.NewClient(
		viper.GetString("usersrv.address"),
		//config.UserServiceUrl(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Error("Failed to dial user service", err)
		os.Exit(1)
	}
	userClient := userpb.NewUserClient(userConn)

	userService := user.NewService(userClient, logger)

	userHandler := api.NewUserHandler(userService, contentService, validator)
	api.InitialSwagger()
	router := router.New(
		userHandler,
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
