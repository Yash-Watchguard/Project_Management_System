package service

import (
	"github.com/Yash-Watchguard/Tasknest/interfaces"
	"github.com/Yash-Watchguard/Tasknest/model"
)

type AuthService struct {
	Repo interfaces.UserRepository
}

func NewAuthService(repo interfaces.UserRepository) *AuthService {
	return &AuthService{Repo: repo}
}

func (s *AuthService) Signup(user *model.User) error {
	return s.Repo.SaveUser(user)
}
func(s * AuthService)Login(name,email,password string)(*model.User,error){
	return s.Repo.IsUserPresent(name,email,password)
}
