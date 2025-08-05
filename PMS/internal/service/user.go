package service

import "github.com/Yash-Watchguard/Tasknest/internal/model/user"

type UserService struct {
	userrepo *user.User
}

func NewUserService(userrepo *user.User) *UserService {
	return &UserService{userrepo: userrepo}
}
func(a * AdminService)ViewProfile(ctx context.Context){
	userId:=ctx.Value(ContextKey.UserId).(string)
	err:=a.userRepo.ViewProfile(userId)
	if err!=nil{
	color.Red("%v",err)
	}
}

func()SaveUser(){

}
func()UpdateUser(){

}
func(a *AdminService)DeleteUser(userId string)error{
	return a.userRepo.DeleteUserById(userId)
}