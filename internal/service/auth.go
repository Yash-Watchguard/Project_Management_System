// package service

// import (
// 	"github.com/Yash-Watchguard/Tasknest/internal/interfaces"
// 	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
// )

// type AuthService struct {
// 	Repo interfaces.UserRepository
// }

// func NewAuthService(repo interfaces.UserRepository) *AuthService {
// 	return &AuthService{Repo: repo}
// }

// func (s *AuthService) Signup(user *user.User) error {
// 	return s.Repo.SaveUser(user)
// }
// func (s *AuthService) Login(name, email, password string) (*user.User, error) {
// 	return s.Repo.IsUserPresent(name, email, password)
// }
package service