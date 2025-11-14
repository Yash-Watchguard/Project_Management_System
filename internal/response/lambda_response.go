package response

import (
	"encoding/json"
	

	"github.com/aws/aws-lambda-go/events"
)

func LambdaSuccessResponse(data any, message string, code int,statusCode  int) (events.APIGatewayProxyResponse){
	response,_:=json.Marshal(response{
		Status: "Success",
		Message: message,
		Data: data,
	})

	return events.APIGatewayProxyResponse{
		StatusCode:statusCode,
		Body: string(response),
	}

}
func LambdaErrorResponse(data any, message string, code int,statusCode  int) (events.APIGatewayProxyResponse){
	response,_:=json.Marshal(response{
		Status: "fail",
		Message: message,
		Data: data,
	})

	return events.APIGatewayProxyResponse{
		StatusCode:statusCode,
		Body: string(response),
	}

}