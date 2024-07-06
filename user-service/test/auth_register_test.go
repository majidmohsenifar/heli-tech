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

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestAuth_Register_InvalidInputs(t *testing.T) {
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
	req := userpb.RegisterRequest{
		Email:           " ",
		Password:        "",
		ConfirmPassword: "",
	}
	res, err := client.Register(ctx, &req)
	assert.Nil(res)
	e, ok := status.FromError(err)
	assert.True(ok)
	assert.Equal(e.Code(), codes.Code(400))
	assert.Equal(e.Message(), "email is empty")

	//invalid email
	req = userpb.RegisterRequest{
		Email:           "invalid",
		Password:        "",
		ConfirmPassword: "",
	}
	res, err = client.Register(ctx, &req)
	assert.Nil(res)
	e, ok = status.FromError(err)
	assert.True(ok)
	assert.Equal(e.Code(), codes.Code(400))
	assert.Equal(e.Message(), "invalid email")

	//empty password
	req = userpb.RegisterRequest{
		Email:           "test@test.com",
		Password:        "",
		ConfirmPassword: "",
	}
	res, err = client.Register(ctx, &req)
	assert.Nil(res)
	e, ok = status.FromError(err)
	assert.True(ok)
	assert.Equal(e.Code(), codes.Code(400))
	assert.Equal(e.Message(), "password is empty")

	//empty confirmPassword
	req = userpb.RegisterRequest{
		Email:           "test@test.com",
		Password:        "123456789",
		ConfirmPassword: "",
	}
	res, err = client.Register(ctx, &req)
	assert.Nil(res)
	e, ok = status.FromError(err)
	assert.True(ok)
	assert.Equal(e.Code(), codes.Code(400))
	assert.Equal(e.Message(), "confirmPassword is empty")

	//different password and confirmPassword
	req = userpb.RegisterRequest{
		Email:           "test@test.com",
		Password:        "123456789",
		ConfirmPassword: "987654321",
	}
	res, err = client.Register(ctx, &req)
	assert.Nil(res)
	e, ok = status.FromError(err)
	assert.True(ok)
	assert.Equal(e.Code(), codes.Code(400))
	assert.Equal(e.Message(), "password and confirmPassword are not the same")
}

func TestAuth_Register_EmailAlreadyExist(t *testing.T) {
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

	_, err = repo.CreateUser(ctx, db, repository.CreateUserParams{
		Email:    "test@test.com",
		Password: "password",
	})
	assert.Nil(err)
	client := userpb.NewUserClient(userConn)

	req := userpb.RegisterRequest{
		Email:           "test@test.com",
		Password:        "123456789",
		ConfirmPassword: "123456789",
	}
	res, err := client.Register(ctx, &req)
	assert.Nil(res)
	e, ok := status.FromError(err)
	assert.True(ok)
	assert.Equal(e.Code(), codes.Code(422))
	assert.Equal(e.Message(), "email already exist")
}

func TestAuth_Register_Successful(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	db := getDB()
	err := truncateDB()
	assert.Nil(err)
	repo := repository.New()
	passwordEncoder := core.NewPasswordEncoder()
	roleRouteManager := auth.NewRoleRouteManager(db, repo)
	authService := auth.NewService(db, repo, passwordEncoder, nil, getLogger(), roleRouteManager)
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

	role, err := repo.CreateRole(ctx, db, auth.RoleEndUser)
	assert.Nil(err)

	client := userpb.NewUserClient(userConn)

	req := userpb.RegisterRequest{
		Email:           "new@test.com",
		Password:        "123456789",
		ConfirmPassword: "123456789",
	}
	res, err := client.Register(ctx, &req)
	assert.Nil(err)
	assert.True(res.Ok)
	//check new created user
	createdUser, err := repo.GetUserByEmail(ctx, db, "new@test.com")
	assert.Nil(err)
	err = passwordEncoder.CompareHashAndPassword(createdUser.Password, "123456789")
	assert.Nil(err)

	//check the user role
	userRoles, err := repo.GetUserRolesByUserID(ctx, db, createdUser.ID)
	assert.Nil(err)
	assert.Equal(len(userRoles), 1)
	assert.Equal(userRoles[0].RoleID, role.ID)
}
