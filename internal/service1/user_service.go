package service1

import (
	"context"
	"errors"
     "golang.org/x/crypto/bcrypt"
	"github.com/Yash-Watchguard/Tasknest/internal/interfaces"
	"github.com/Yash-Watchguard/Tasknest/internal/util"

	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
)
//go:generate mockgen -source=user_service.go -destination=../mocks/mock_userservice.go -package=mocks
type UserServiceInterface interface{
    ViewProfile( userId string) ([]user.User, error)
    ViewAllUsers() ([]user.User, error)
    DeleteUser(userId string) error
    GetAllManager(ctx context.Context) ([]user.User,error)
    UpdateUser(id string, updates map[string]interface{}) error
    PromoteEmployee( employeeId string) error
    ViewAllEmplpyee(ctx context.Context)([]user.User,error)
    IsUserPresent(name string, email string,password string)(*user.User,error)
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
        mp:=make(map[string]interface{})
        mp["status"]="InActive"
		return u.userRepo.UpdateProfile(userId,mp)
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
        finalUpdates["name"] = name
    }
    if email, ok := updates["email"].(string); ok {
        // validate email format
        if err := util.ValidateEmail(email); err != nil {
            return err
        }

        // check uniqueness
        existingUser, err := us.userRepo.GetUserByEmail(email)
        if err == nil && existingUser != nil {
            return errors.New("email already exists")
        }
        finalUpdates["email"] = email
    }
    if password, ok := updates["password"].(string); ok {
        // validate password strength
        if err := util.ValidatePassword(password); err != nil {
            return err
        }

        hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        finalUpdates["password"] = string(hashedPassword)
    }
    if phone, ok := updates["phoneNumber"].(string); ok {
        finalUpdates["phone_number"] = phone // match DB column name
    }

    if len(finalUpdates) == 0 {
        return errors.New("no valid fields to update")
    }

    
    return us.userRepo.UpdateProfile(id, finalUpdates)
}



func (u * UserService) PromoteEmployee( employeeId string) error {
	
	return u.userRepo.PromoteEmployee(employeeId)
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
func(u *UserService)IsUserPresent(name string, email string,password string)(*user.User,error){
    return u.userRepo.IsUserPresent(name,email,password)
}