package main

import (
	"context"
	"encoding/json"

	"net/http"

	"github.com/Yash-Watchguard/Tasknest/internal/config"
	"github.com/Yash-Watchguard/Tasknest/internal/model/user"
	"github.com/Yash-Watchguard/Tasknest/internal/response"
	"github.com/Yash-Watchguard/Tasknest/internal/util"

	"github.com/Yash-Watchguard/Tasknest/internal/repository"
	"github.com/Yash-Watchguard/Tasknest/internal/service1"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	
	
)

var authService service1.AuthServiceInterface

func init() {
	//  get the db instance and services
	

	dynmoDbClient := config.GetDyanoDbCliebt()

	userRepo := repository.NewUserRepo(dynmoDbClient, "TaskNest")

	authService = service1.NewAuthService(userRepo)

}
func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var err error = nil
	var userRequest struct {
		Email    string
		Password string
	}

	if err = json.Unmarshal([]byte(event.Body), &userRequest); err != nil {
		return response.LambdaErrorResponse(nil,"Invalid Request body",1001,http.StatusBadRequest),nil
	}

	person := &user.User{}

	person, err = authService.Login("yash", userRequest.Email, userRequest.Password)

	if err != nil {
		return response.LambdaErrorResponse(nil,"Invalid email and password",1002,http.StatusInternalServerError), nil
	}

	person.Status = user.Active

	

	var jwtTokenString string

	jwtTokenString, _ = util.GenerateJwt(person.Id, person.Role)
	data := map[string]any{"token": jwtTokenString, "name": person.Name, "UserId": person.Id}
	return response.LambdaSuccessResponse(data,"Token generated Successfully",1001,http.StatusAccepted),nil
}
