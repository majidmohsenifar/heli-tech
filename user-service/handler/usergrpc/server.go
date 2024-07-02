package usergrpc

import (
	"git.energy/corepass/notification-notifier-service/service/user"
	userpb "git.energy/ting/data-contracts/proto/user"
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
