package interfaces

import (

	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
)

//go:generate mockgen -source=user_repo_interface.go -destination=../mocks/mock_userrepository.go -package=mocks
type UserRepository interface {
	SaveUser(user *user.User)error
	IsUserPresent(name,email,password string)(*user.User,error)
	ViewProfile(userId string)([]user.User,error)
	GetAllUsers()([]user.User,error)
	DeleteUserById(userId string)error
	UpdateProfile(userId string,mp map[string]interface{})error
	GetAllManager()([]user.User,error)
	PromoteEmployee(employeeId string) error
	ViewAllEmployee() ([]user.User, error) 
	GetUserByEmail(email string)(*user.User,error)
}