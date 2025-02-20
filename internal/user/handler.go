package user

import (
	"encoding/json"
	"fmt"
	"lambda-dropin/internal/jwt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type UserHandler struct {
	userService *UserService
}

func NewUserHandler(userService *UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) RegisterUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var userRequest UserRequest
	err := json.Unmarshal([]byte(request.Body), &userRequest)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid request",
			StatusCode: http.StatusBadRequest,
		}, err

	}

	if userRequest.Username == "" || userRequest.Password == "" {
		return events.APIGatewayProxyResponse{
			Body:       "invalid request, the fields cannot be empty",
			StatusCode: http.StatusBadRequest,
		}, fmt.Errorf("Invalid request, fields can not be empty %w", err)
	}

	userExist, err := h.userService.DoesUserExist(userRequest.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	if userExist {
		return events.APIGatewayProxyResponse{
			Body:       "User already exists",
			StatusCode: http.StatusConflict,
		}, fmt.Errorf("User already exists %w: ", err)
	}

	newUser, err := h.userService.NewUserWithHashedPassword(userRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error",
			StatusCode: http.StatusConflict,
		}, fmt.Errorf("Could not create new user %w: ", err)
	}

	// we know that this user does not exist
	err = h.userService.InsertUser(newUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("Error inserting user %w", err)
	}

	return events.APIGatewayProxyResponse{
		Body:       "Success",
		StatusCode: http.StatusOK,
	}, nil
}

func (h *UserHandler) LoginUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var loginRequest UserRequest
	err := json.Unmarshal([]byte(request.Body), &loginRequest)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid request",
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	newUser, err := h.userService.GetUser(loginRequest.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("Could not get user %w", err)

	}

	if !h.userService.ValidatePassword(newUser.PasswordHash, loginRequest.Password) {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid user credentials",
			StatusCode: http.StatusBadGateway,
		}, nil

	}

	accessToken, err := jwt.CreateToken(newUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("Could not create token %w", err)
	}

	successMsg := fmt.Sprintf(`{"access_token": "%s"}`, accessToken)

	return events.APIGatewayProxyResponse{
		Body:       successMsg,
		StatusCode: http.StatusOK,
	}, nil
}
