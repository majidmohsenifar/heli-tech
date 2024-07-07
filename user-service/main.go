package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"

	userpb "github.com/majidmohsenifar/heli-tech/data-contracts/proto/user"
	"github.com/majidmohsenifar/heli-tech/user-service/config"
	"github.com/majidmohsenifar/heli-tech/user-service/core"
	"github.com/majidmohsenifar/heli-tech/user-service/handler/usergrpc"
	"github.com/majidmohsenifar/heli-tech/user-service/logger"
	"github.com/majidmohsenifar/heli-tech/user-service/repository"
	"github.com/majidmohsenifar/heli-tech/user-service/service/auth"
	"github.com/majidmohsenifar/heli-tech/user-service/service/jwt"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	grpcPort = "50051"
)

func main() {
	ctx := context.Background()
	viper := config.NewViper("./config/")
	logger := logger.NewLogger()
	dbClient, err := core.NewDBClient(ctx, viper.GetString("db.dsn"))
	if err != nil {
		logger.Error("failed to initiate a db client", err)
		os.Exit(1)
	}
	defer dbClient.Close()
	repo := repository.New()
	passwordEncoder := core.NewPasswordEncoder()
	jwtService, err := jwt.NewService(viper)
	if err != nil {
		logger.Error("failed to create jwt service", err)
		os.Exit(1)
	}
	roleRouteManager := auth.NewRoleRouteManager(dbClient, repo)
	authService := auth.NewService(
		dbClient,
		repo,
		passwordEncoder,
		jwtService,
		logger,
		roleRouteManager,
	)

	grpcPanicRecoveryHandler := func(p any) error {
		err := errors.New("recovered from panic")
		tempErr, ok := p.(error)
		if ok {
			err = tempErr
		} else {
			panicStr, ok := p.(string)
			if ok {
				err = errors.New(panicStr)
			}
		}
		logger.Error("recovered from panic", err)
		return status.Errorf(codes.Internal, "%s", "something went wrong")
	}
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),
		grpc.ChainStreamInterceptor(
			recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),
	)
	userGrpcServer := usergrpc.NewServer(
		authService,
	)
	userpb.RegisterUserServer(grpcServer, userGrpcServer)
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", grpcPort))
	if err != nil {
		logger.Error("can not listen to grpcPort", err)
		os.Exit(1)
	}
	err = grpcServer.Serve(l)
	if err != nil {
		logger.Error("can not serv grpc", err)
		os.Exit(1)
	}
}
