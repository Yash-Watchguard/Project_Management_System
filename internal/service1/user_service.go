package service1

import (
	"context"
	"errors"
	"fmt"

	"github.com/Yash-Watchguard/Tasknest/internal/interfaces"
	"github.com/Yash-Watchguard/Tasknest/internal/util"
	"golang.org/x/crypto/bcrypt"

	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
)

//go:generate mockgen -source=user_service.go -destination=../mocks/mock_userservice.go -package=mocks
type UserServiceInterface interface{
    ViewProfile( userId string) ([]user.User, error)
    ViewAllUsers() ([]user.User, error)
    DeleteUser(email string) error
    GetAllManager(ctx context.Context) ([]user.User,error)
    UpdateUser(id string, updates map[string]interface{}) error
    PromoteEmployee( employeeId string) error
    ViewAllEmplpyee(ctx context.Context)([]user.User,error)
    IsUserPresent(name,email string,password string)(*user.User,error)
    CheckUserExist(email string)(bool)
}
type UserService struct{
	userRepo    interfaces.UserRepository
}

func NewUserService(userRepo interfaces.UserRepository)UserServiceInterface{
return &UserService{userRepo: userRepo}
}

func (u *UserService) ViewProfile( userId string) ([]user.User, error) {
		return u.userRepo.ViewProfile(userId)
}

func (u *UserService) ViewAllUsers() ([]user.User, error) {
	return u.userRepo.GetAllUsers()
}
func (u * UserService) DeleteUser(userId string) error {
       return u.userRepo.DeleteUserById(userId)
}

func (u *UserService) GetAllManager(ctx context.Context) ([]user.User,error) {
	userId := ctx.Value(ContextKey.UserRole).(roles.Role)
	if userId != 0 {
		return nil,errors.New("unautherized access")
	}
	return u.userRepo.GetAllManager()
}

func (us *UserService) UpdateUser(id string, updates map[string]interface{}) error {
    // check if user exists
    _, err := us.userRepo.ViewProfile(id)
    if err != nil {
        return err
    }

    finalUpdates := make(map[string]interface{})

    // Apply partial updates
    if name, ok := updates["name"].(string); ok {
        if(len(name)!=0){
            finalUpdates["name"] = name
        }
    }
    if email, ok := updates["email"].(string); ok {
        if(len(email)!=0){
             if err := util.ValidateEmail(email); err != nil {
            return errors.New("email is not valid")
        }

        // check uniqueness
        existingUser, err := us.userRepo.GetUserByEmail(email)
        if err == nil && existingUser != nil {
            return errors.New("email already exists")
        }
        finalUpdates["email"] = email
        }
        // validate email format
       
    }
    if password, ok := updates["password"].(string); ok {

        // validate password strength
        if(len(password)!=0){
            if err := util.ValidatePassword(password); err != nil {
            return errors.New("please enter valid password")
        }

        hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        finalUpdates["password"] = string(hashedPassword)
        }
        
    }
    if phone, ok := updates["phoneNumber"].(string); ok {
        if(len(phone)!=0){
             if err:=util.ValidateMobileNumber(phone);err!=nil{
            return errors.New("phone number is not valid")
        }
        finalUpdates["phone_number"] = phone 
        }
       
    }

    if len(finalUpdates) == 0 {
        return errors.New("no valid fields to update")
    }

    fmt.Println(finalUpdates)
    return us.userRepo.UpdateProfile(id, finalUpdates)
}



func (u * UserService) PromoteEmployee( email string) error {
	
	return u.userRepo.PromoteEmployee(email)
}
func(u *UserService)ViewAllEmplpyee(ctx context.Context)([]user.User,error){
	userRole:=ctx.Value(ContextKey.UserRole).(roles.Role)

    if userRole==2{
		return []user.User{},errors.New("unauthorized access")
	}
	return u.userRepo.ViewAllEmployee()
}
func(u *UserService)CheckUserExist(email string)(bool){
	user, err:=u.userRepo.GetUserByEmail(email)

	if err == nil && user!=nil{
		 return true
	}
	return false
}
func(u *UserService)IsUserPresent(name,email string,password string)(*user.User,error){
    return u.userRepo.IsUserPresent(name,email,password)
}