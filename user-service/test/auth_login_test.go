package test

import (
	"context"
	"net"
	"testing"

	userpb "github.com/majidmohsenifar/heli-tech/data-contracts/proto/user"
	"github.com/majidmohsenifar/heli-tech/user-service/core"
	"github.com/majidmohsenifar/heli-tech/user-service/handler/usergrpc"
	"github.com/majidmohsenifar/heli-tech/user-service/repository"
	"github.com/majidmohsenifar/heli-tech/user-service/service/auth"
	"github.com/majidmohsenifar/heli-tech/user-service/service/jwt"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestAuth_Login_InvalidInputs(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	userGrpcServer := usergrpc.NewServer(nil)

	l, err := net.Listen("tcp", "127.0.0.1:0")
	assert.Nil(err)
	defer l.Close()
	googleGrpcServer := grpc.NewServer()

	userpb.RegisterUserServer(googleGrpcServer, userGrpcServer)
	userConn, err := grpc.NewClient(
		l.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.Nil(err)
	go googleGrpcServer.Serve(l)
	client := userpb.NewUserClient(userConn)

	//empty email
	req := userpb.LoginRequest{
		Email:    " ",
		Password: "",
	}
	res, err := client.Login(ctx, &req)
	assert.Nil(res)
	e, ok := status.FromError(err)
	assert.True(ok)
	assert.Equal(e.Code(), codes.Code(400))
	assert.Equal(e.Message(), "email is empty")

	//invalid email
	req = userpb.LoginRequest{
		Email:    "invalid",
		Password: "",
	}
	res, err = client.Login(ctx, &req)
	assert.Nil(res)
	e, ok = status.FromError(err)
	assert.True(ok)
	assert.Equal(e.Code(), codes.Code(400))
	assert.Equal(e.Message(), "invalid email")

	//empty password
	req = userpb.LoginRequest{
		Email:    "test@test.com",
		Password: "",
	}
	res, err = client.Login(ctx, &req)
	assert.Nil(res)
	e, ok = status.FromError(err)
	assert.True(ok)
	assert.Equal(e.Code(), codes.Code(400))
	assert.Equal(e.Message(), "password is empty")

}

func TestAuth_Login_EmailDoesNotExist(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	db := getDB()
	err := truncateDB()
	assert.Nil(err)
	repo := repository.New()
	authService := auth.NewService(db, repo, nil, nil, getLogger(), nil)
	userGrpcServer := usergrpc.NewServer(authService)

	l, err := net.Listen("tcp", "127.0.0.1:0")
	assert.Nil(err)
	defer l.Close()
	googleGrpcServer := grpc.NewServer()

	userpb.RegisterUserServer(googleGrpcServer, userGrpcServer)
	userConn, err := grpc.NewClient(
		l.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.Nil(err)
	go googleGrpcServer.Serve(l)
	client := userpb.NewUserClient(userConn)

	req := userpb.LoginRequest{
		Email:    "test@test.com",
		Password: "123456789",
	}
	res, err := client.Login(ctx, &req)
	assert.Nil(res)
	e, ok := status.FromError(err)
	assert.True(ok)
	assert.Equal(e.Code(), codes.Code(401))
	assert.Equal(e.Message(), "invalid username or password")
}

func TestAuth_Login_IncorrectPassword(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	db := getDB()
	err := truncateDB()
	passwordEncoder := core.NewPasswordEncoder()
	assert.Nil(err)
	repo := repository.New()
	authService := auth.NewService(db, repo, passwordEncoder, nil, getLogger(), nil)
	userGrpcServer := usergrpc.NewServer(authService)

	l, err := net.Listen("tcp", "127.0.0.1:0")
	assert.Nil(err)
	defer l.Close()
	googleGrpcServer := grpc.NewServer()

	userpb.RegisterUserServer(googleGrpcServer, userGrpcServer)
	userConn, err := grpc.NewClient(
		l.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.Nil(err)
	go googleGrpcServer.Serve(l)

	hashedPassword, err := passwordEncoder.GenerateFromPassword("123456789")
	assert.Nil(err)
	_, err = repo.CreateUser(ctx, db, repository.CreateUserParams{
		Email:    "test@test.com",
		Password: string(hashedPassword),
	})
	assert.Nil(err)

	client := userpb.NewUserClient(userConn)

	//incorrect password
	req := userpb.LoginRequest{
		Email:    "test@test.com",
		Password: "incorrect",
	}
	res, err := client.Login(ctx, &req)
	assert.Nil(res)
	e, ok := status.FromError(err)
	assert.True(ok)
	assert.Equal(e.Code(), codes.Code(401))
	assert.Equal(e.Message(), "invalid username or password")
}

func TestAuth_Login_Successful(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	db := getDB()
	err := truncateDB()
	passwordEncoder := core.NewPasswordEncoder()
	assert.Nil(err)
	repo := repository.New()
	viper := getViperConfig()
	jwtService, err := jwt.NewService(viper)
	assert.Nil(err)
	authService := auth.NewService(db, repo, passwordEncoder, jwtService, getLogger(), nil)
	userGrpcServer := usergrpc.NewServer(authService)

	l, err := net.Listen("tcp", "127.0.0.1:0")
	assert.Nil(err)
	defer l.Close()
	googleGrpcServer := grpc.NewServer()

	userpb.RegisterUserServer(googleGrpcServer, userGrpcServer)
	userConn, err := grpc.NewClient(
		l.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.Nil(err)
	go googleGrpcServer.Serve(l)

	hashedPassword, err := passwordEncoder.GenerateFromPassword("123456789")
	assert.Nil(err)
	_, err = repo.CreateUser(ctx, db, repository.CreateUserParams{
		Email:    "test@test.com",
		Password: string(hashedPassword),
	})
	assert.Nil(err)

	client := userpb.NewUserClient(userConn)
	req := userpb.LoginRequest{
		Email:    "test@test.com",
		Password: "123456789",
	}
	res, err := client.Login(ctx, &req)
	assert.Nil(err)
	assert.NotEqual(res.Token, "")
}
