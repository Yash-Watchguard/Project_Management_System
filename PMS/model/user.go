package model

type User struct{
	Id string `json:"id"`
	Role string `jason:"role"`
	Name string `json:"name"`
	Password string `json:"password"`
	PhoneNumber string `json:"phonenumber"`
	Email string `json:"email"`
}