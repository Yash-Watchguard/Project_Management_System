package user

import "github.com/Yash-Watchguard/Tasknest/internal/model/roles"



type User struct {
	Id          string `json:"id"`
	Role        roles.Role `json:"role"`
	Name        string `json:"name"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phonenumber"`
	Email       string `json:"email"`
}