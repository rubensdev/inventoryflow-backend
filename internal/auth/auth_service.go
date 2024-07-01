package auth

import (
	"errors"

	userdom "github.com/rubensdev/inventoryflow-backend/internal/user"
)

type AuthService struct {
	userSrv userdom.UserService
}

func NewAuthService(userSrv userdom.UserService) *AuthService {
	return &AuthService{
		userSrv: userSrv,
	}
}

func (s *AuthService) ValidateCredentials(loginReq LoginRequest) (*userdom.User, error) {
	user, err := s.userSrv.GetByUsername(loginReq.Username)
	if err != nil {
		return nil, errors.New("wrong credentials")
	}

	matches, err := user.Password.Matches(loginReq.Password)
	if err != nil {
		return nil, err
	}
	if !matches {
		return nil, errors.New("wrong credentials")
	}

	return user, nil
}
