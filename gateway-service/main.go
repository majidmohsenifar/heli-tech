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

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	socialHandler := api.NewSocialHandler(socialService, contentService, validator)
	contentHandler := api.NewContentHandler(contentService, validator)
	notifHandler := api.NewNotifHandler(notifService, validator)
	roomHandler := api.NewRoomHandler(roomService, validator)
	api.InitialSwagger()
	router := router.New(
		notifHandler,
		userHandler,
		socialHandler,
		contentHandler,
		roomHandler,
		authService,
		userService,
		logger,
	)
	if err != nil {
		logger.Fatal("failed to initialize router", zap.Error(err))
	}
	httpServer := &http.Server{
		Addr:    config.ServerHttpAddress(),
		Handler: router.Handler,
	}
	logger.Info("starting HTTP server on %s", zap.String("HTTP server address: ", config.ServerHttpAddress()))
	err = httpServer.ListenAndServe()
	if err != nil {
		logger.Fatal(err.Error())
	}
}
