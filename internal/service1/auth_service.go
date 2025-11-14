package service1

import (
	"github.com/Yash-Watchguard/Tasknest/internal/interfaces"
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
)
//go:generate mockgen -source=auth_service.go -destination=../mocks/mock_authservice.go -package=mocks
type AuthServiceInterface interface{
	Signup(user *user.User) error
	Login(name, email, password string) (*user.User, error)
}

type AuthService struct {
	Repo interfaces.UserRepository
}

func NewAuthService(repo interfaces.UserRepository) AuthServiceInterface{
	return &AuthService{Repo: repo}
}

func (s *AuthService) Signup(user *user.User) error {
	return s.Repo.SaveUser(user)
}
func (s *AuthService) Login(name, email, password string) (*user.User, error) {
	return s.Repo.IsUserPresent(name,email, password)
}
