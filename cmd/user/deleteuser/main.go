package main

import (
	"context"
	"fmt"

	"github.com/Yash-Watchguard/Tasknest/internal/config"
	"github.com/Yash-Watchguard/Tasknest/internal/middleware"
	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
	"github.com/Yash-Watchguard/Tasknest/internal/response"

	"github.com/Yash-Watchguard/Tasknest/internal/repository"

	"github.com/Yash-Watchguard/Tasknest/internal/service1"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var userService service1.UserServiceInterface

func init(){
 dynamoClient:=config.GetDyanoDbCliebt()

 userRepo:=repository.NewUserRepo(dynamoClient,"TaskNest")
 userService=service1.NewUserService(userRepo)
}
func main(){
   lambda.Start(middleware.LambdaAuthMiddleWare(handler))
}


func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

    email := req.PathParameters["user_email"]
    role:=ctx.Value(ContextKey.UserRole).(roles.Role)
	if(role!=roles.Admin){
		return response.LambdaErrorResponse(nil,"unauthorized access",403,403),nil
	}
	fmt.Println("Extracted email:", email)

	if email == "" {	
		return response.LambdaErrorResponse(nil,"missing user_email",400,400),nil
	}
	error:= userService.DeleteUser(email)
	if error != nil {
			return response.LambdaErrorResponse(nil,"error deleting user"+error.Error(),500,500),nil
	}
	   return response.LambdaSuccessResponse(nil,"user deleted successfully",200,200),nil
}

