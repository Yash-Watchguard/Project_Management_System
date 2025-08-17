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
	UpdateProfile(userId string,field,updatedData string)error
	GetAllManager()([]user.User,error)
	PromoteEmployee(employeeId string) error
	ViewAllEmployee() ([]user.User, error) 
}