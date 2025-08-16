package service1

import (
	"context"
	"errors"

	"github.com/Yash-Watchguard/Tasknest/internal/interfaces"
	
	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
	"github.com/fatih/color"
)
type UserService struct{
	userRepo    interfaces.UserRepository
}

func NewUserService(userRepo interfaces.UserRepository)*UserService{
return &UserService{userRepo: userRepo}
}

func (u *UserService) ViewProfile(ctx context.Context, userId string) ([]user.User, error) {
	userID := ctx.Value(ContextKey.UserId).(string)
	userRole := ctx.Value(ContextKey.UserRole).(roles.Role)

	if userID == userId || userRole == 0 {
		return u.userRepo.ViewProfile(userId)
	}
	return nil, errors.New("unauthorized access")
}

func (u *UserService) ViewAllUsers(ctx context.Context) ([]user.User, error) {

	userID := ctx.Value(ContextKey.UserRole).(roles.Role)
	if userID != 0 {
		return []user.User{}, errors.New("unautherized access")
	}
	return u.userRepo.GetAllUsers()

}
func (u * UserService) DeleteUser(ctx context.Context, userId string) error {
	userID := ctx.Value(ContextKey.UserId).(string)
	userRole := ctx.Value(ContextKey.UserRole).(roles.Role)

	if userId == userID || userRole == 0 {
		err := u.userRepo.DeleteUserById(userId)
		if err != nil {
			color.Red("%v", err)
		}
	} else {
		return errors.New("unauthorized access")
	}
	return nil
}

func (u *UserService) GetAllManager(ctx context.Context) error {
	userId := ctx.Value(ContextKey.UserRole).(roles.Role)
	if userId != 0 {
		return errors.New("unautherized access")
	}
	return u.userRepo.GetAllManager()
}

func (u * UserService) UpdateProfile(userId string, ctx context.Context, name string, email string, password string, number string) error {
	userID := ctx.Value(ContextKey.UserId).(string)

	if userID != userId {
		return errors.New("unauthorized access")
	}
	return u.userRepo.UpdateProfile(userId, name, email, password, number)
}

func (u * UserService) PromoteEmployee(ctx context.Context, employeeId string) error {
	userRole := ctx.Value(ContextKey.UserRole).(roles.Role)

	if userRole ==2  {
		return errors.New("unauthorized person")
	}
	return u.userRepo.PromoteEmployee(employeeId)
}
func(u *UserService)ViewAllEmplpyee(ctx context.Context)([]user.User,error){
	userRole:=ctx.Value(ContextKey.UserRole).(roles.Role)

    if userRole==2{
		return []user.User{},errors.New("unauthorized access")
	}
	return u.userRepo.ViewAllEmployee()
}