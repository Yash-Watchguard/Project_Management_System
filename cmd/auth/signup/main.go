package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Yash-Watchguard/Tasknest/internal/config"
	"github.com/Yash-Watchguard/Tasknest/internal/model"
	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"

	
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
	"github.com/Yash-Watchguard/Tasknest/internal/repository"
	"github.com/Yash-Watchguard/Tasknest/internal/response"

	
	"github.com/Yash-Watchguard/Tasknest/internal/service1"
	"github.com/Yash-Watchguard/Tasknest/internal/util"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"golang.org/x/crypto/bcrypt"
)

var authService service1.AuthServiceInterface
var userservice service1.UserServiceInterface

func init() {
	dynmoDbClient:=config.GetDyanoDbCliebt()

	userRepo:=repository.NewUserRepo(dynmoDbClient,"TaskNest")

	authService = service1.NewAuthService(userRepo)
	userservice = service1.NewUserService(userRepo)
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var err error
	var newUser struct {
		Name        string `json:"name" validate:"required"`
		Email       string `json:"email" validate:"required,email"`
		Password    string `json:"password" validate:"required,min=8"`
		PhoneNumber string `json:"phonenumber" validate:"required"`
	}

	if err = json.Unmarshal([]byte(event.Body), &newUser); err != nil {
		return response.LambdaErrorResponse(nil,"Invalid request body",1001,http.StatusBadRequest), nil
	}

		flag := userservice.CheckUserExist(newUser.Email)
		if flag {
			return response.LambdaErrorResponse(nil,"user already exists",1001,http.StatusConflict), nil
		}

		if err = util.ValidatePassword(newUser.Password); err != nil {
			return response.LambdaErrorResponse(nil,"Invalid password",1001,http.StatusBadRequest), nil
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
		if err != nil {
			return response.LambdaErrorResponse(nil,"error in hashing",1001,http.StatusInternalServerError), nil
		}

		newUser.Password = string(hashedPassword)

		NewUser := &user.User{
			Id:          util.GenerateUniqueUUID(),
			Name:        newUser.Name,
			Email:       newUser.Email,
			Password:    newUser.Password,
			Role:        roles.Employee,
			PhoneNumber: newUser.PhoneNumber,
			Status:user.Active,
		}
		err = authService.Signup(NewUser)
		if err != nil {
           return response.LambdaErrorResponse(nil,fmt.Sprint("Signup failed",err.Error()),1001,http.StatusInternalServerError), nil
		}

		userDto:=model.UserDto{
		Id: NewUser.Id,
		Name: NewUser.Name,
		Email: NewUser.Email,
		PhoneNumber: NewUser.PhoneNumber,
		Role: roles.RoleParser(NewUser.Role),
		}

		
		return response.LambdaErrorResponse(userDto,"User created Successfully",1001,http.StatusCreated), nil
}

