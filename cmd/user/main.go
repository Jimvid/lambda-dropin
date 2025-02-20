package main

import (
	"lambda-dropin/internal/database"
	"lambda-dropin/internal/middleware"
	"lambda-dropin/internal/user"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Username string `json:"username"`
}

func ProtectedHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       "Hello there, this is protected",
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	db := database.NewDynamoDB()
	userService := user.NewUserService(db)
	userHandler := user.NewUserHandler(userService)

	lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch request.Path {
		case "/register":
			return userHandler.RegisterUserHandler(request)

		case "/login":
			return userHandler.LoginUserHandler(request)

		case "/protected":
			return middleware.ValidateJWTMiddleware(ProtectedHandler)(request)

		default:
			return events.APIGatewayProxyResponse{
				Body:       "Not Found",
				StatusCode: http.StatusNotFound,
			}, nil
		}
	})
}
