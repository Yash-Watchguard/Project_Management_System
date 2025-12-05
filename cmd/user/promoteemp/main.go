package main

import (
	"context"
	"fmt"

	"github.com/Yash-Watchguard/Tasknest/internal/config"
	"github.com/Yash-Watchguard/Tasknest/internal/middleware"
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
    fmt.Println("Extracted email:", email)

    if email == "" {

		return response.LambdaErrorResponse(nil,"missing user_email",400,400),nil
       
    }

    // Call your repository
    err := userService.PromoteEmployee(email)
    if err != nil {
    
       return response.LambdaErrorResponse(nil,"missing user_email",500,500),nil
    }

   return response.LambdaSuccessResponse(nil,"ptomoted",200,200),nil
}

