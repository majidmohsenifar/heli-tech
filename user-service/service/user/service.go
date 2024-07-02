package user

import "errors"

var (
	ErrEmailAlreadyExist = errors.New("email already exist")
	ErrUserNotFound      = errors.New("user not found")
)

type Service struct{}

type RegisterParams struct {
	Email    string
	Password string
}

type LoginParams struct {
	Email    string
	Password string
}

type GetUserDataByTokenParams struct {
	Token string
	Path  string
}

type GetUserDataByTokenResponse struct {
	Email string
	ID    int64
}

func (s *Service) Register(params RegisterParams) error {
	panic("imple later")
}

func (s *Service) Login(params LoginParams) (string, error) {
	panic("imple later")
}

func (s *Service) GetUserDataByToken(params GetUserDataByTokenParams) (GetUserDataByTokenResponse, error) {
	panic("imple later")
}

func NewService() *Service {
	return &Service{}
}
