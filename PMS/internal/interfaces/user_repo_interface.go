package interfaces

import (

	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
)

type UserRepository interface {
	SaveUser(user *user.User)error
	// IsUserPresent(name,email,password string)(*user.User,error)
	ViewProfile(userId string)error
	// GetAllUsers()[]user.User
	DeleteUserById(userId string)error
	// GetAllManager()error
}