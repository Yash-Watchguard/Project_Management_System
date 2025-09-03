package user

import (
	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
	 
)



type User struct {
    Id          string     `json:"id" db:"id"`            
    Role        roles.Role `json:"role" db:"role"`         
    Name        string     `json:"name" db:"name"`
    Password    string     `json:"password" db:"password"`
    PhoneNumber string     `json:"phonenumber" db:"phone_number"`
    Email       string     `json:"email" db:"email"`
    Status      UserStatus `json:"status" db:"status"`
}