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

type DynamoUser struct{
    PK          string     `json:"PK" dynamodbav:"PK"`            
    SK          string     `json:"SK" dynamodbav:"SK"`
    Id          string     `json:"Id" dynamodbav:"Id"`
    Role        roles.Role `json:"Role" dynamodbav:"Role"`
    Name        string     `json:"Name" dynamodbav:"Name"`
    Password    string     `json:"Password" dynamodbav:"Password"`
    PhoneNumber string     `json:"PhoneNumber" dynamodbav:"PhoneNumber"`
    Email       string     `json:"Email" dynamodbav:"Email"`
    Status      UserStatus `json:"Status" dynamodbav:"Status"`
}