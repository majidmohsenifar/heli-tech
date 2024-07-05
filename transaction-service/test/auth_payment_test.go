package test

import (
	"context"
	"net"
	"testing"

	paymentpb "github.com/majidmohsenifar/heli-tech/data-contracts/proto/payment"
	"github.com/majidmohsenifar/heli-tech/transaction-service/handler/transactiongrpc"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestPayment_Withdraw_InvalidInputs(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	paymentGrpcServer := transactiongrpc.NewServer(nil)

	l, err := net.Listen("tcp", "127.0.0.1:0")
	assert.Nil(err)
	defer l.Close()
	googleGrpcServer := grpc.NewServer()

	paymentpb.RegisterPaymentServer(googleGrpcServer, paymentGrpcServer)
	userConn, err := grpc.NewClient(
		l.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.Nil(err)
	go googleGrpcServer.Serve(l)
	client := paymentpb.NewPaymentClient(userConn)

	req := paymentpb.RegisterRequest{
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
	req = paymentpb.RegisterRequest{
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
	req = paymentpb.RegisterRequest{
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
	req = paymentpb.RegisterRequest{
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
	req = paymentpb.RegisterRequest{
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
