package interfaces

import "github.com/Yash-Watchguard/Tasknest/model"

type UserRepository interface {
	SaveUser(user *model.User)error
	IsUserPresent(name,email,password string)(*model.User,error)
	ViewProfile(user *model.User)
	GetAllUsers()[]model.User
	DeleteUserById(userId string)error
	GetAllManager()error
}