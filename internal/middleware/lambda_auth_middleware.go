package middleware

import (
	"context"
	"net/http"
	"strings"

	ContextKey "github.com/Yash-Watchguard/Tasknest/internal/model/context_key"
	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
	"github.com/Yash-Watchguard/Tasknest/internal/response"
	"github.com/Yash-Watchguard/Tasknest/internal/util"
	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt/v5"
)

func LambdaAuthMiddleWare(
	fn func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error),
) func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

		authorizationHeader := req.Headers["Authorization"]
		if authorizationHeader == "" {
			return response.LambdaErrorResponse(nil, "Missing token", 1002, http.StatusUnauthorized), nil
		}

		if !strings.HasPrefix(authorizationHeader, "Bearer ") {
			return response.LambdaErrorResponse(nil, "Invalid token format", 1002, http.StatusUnauthorized), nil
		}

		tokenString := strings.TrimPrefix(authorizationHeader, "Bearer ")

		token, err := util.VarifyJwt(tokenString)
		if err != nil {
			return response.LambdaErrorResponse(nil, "Invalid token", 1002, http.StatusUnauthorized), nil
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return response.LambdaErrorResponse(nil, "Invalid token", 1002, http.StatusUnauthorized), nil
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			return response.LambdaErrorResponse(nil, "Invalid token", 1002, http.StatusUnauthorized), nil
		}

		roleFloat, ok := claims["role"].(float64)
		if !ok {
			return response.LambdaErrorResponse(nil, "Invalid token", 1002, http.StatusUnauthorized), nil
		}
		role := roles.Role(int(roleFloat))

		ctx = context.WithValue(ctx, ContextKey.UserId, userID)
		ctx = context.WithValue(ctx, ContextKey.UserRole, role)

		return fn(ctx, req)
	}
}
