package test

import (
	"context"
	"net"
	"testing"

	userpb "github.com/majidmohsenifar/heli-tech/data-contracts/proto/user"
	"github.com/majidmohsenifar/heli-tech/user-service/handler/usergrpc"

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
