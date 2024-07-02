package usergrpc

import (
	userpb "github.com/majidmohsenifar/heli-tech/data-contracts/proto/user"
	"github.com/majidmohsenifar/heli-tech/user-service/service/user"
)

type server struct {
	userService *user.Service
	userpb.UnimplementedUserServer
}

func NewServer(
	userService *user.Service,
) userpb.UserServer {
	return &server{
		userService: userService,
	}
}
