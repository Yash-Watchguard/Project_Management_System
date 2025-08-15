package interfaces

import (

	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
)

type UserRepository interface {
	SaveUser(user *user.User)error
	IsUserPresent(name,email,password string)(*user.User,error)
	ViewProfile(userId string)([]user.User,error)
	GetAllUsers()([]user.User,error)
	DeleteUserById(userId string)error
	UpdateProfile(userId string,name string, email string,password string,number string)error
	GetAllManager()error
	PromoteEmployee(employeeId string) error
	ViewAllEmployee() ([]user.User, error) 
}