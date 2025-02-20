package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"lambda-dropin/internal/jwt"
)

func ValidateJWTMiddleware(next func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		// extract the headers from our token
		tokenString := extractTokenFromHeaders(request.Headers)
		if tokenString == "" {
			return events.APIGatewayProxyResponse{
				Body:       "Missing Auth token",
				StatusCode: http.StatusUnauthorized,
			}, nil
		}

		claims, err := jwt.VerifyToken(tokenString)

		if err != nil {
			return events.APIGatewayProxyResponse{
				Body:       "Invalid token",
				StatusCode: http.StatusUnauthorized,
			}, nil
		}

		expires := int64(claims["expires"].(float64))
		if time.Now().Unix() > expires {
			return events.APIGatewayProxyResponse{
				Body:       "token expired",
				StatusCode: http.StatusUnauthorized,
			}, nil
		}

		return next(request)
	}

}

func extractTokenFromHeaders(headers map[string]string) string {
	tokenString, ok := headers["Authorization"]

	if !ok {
		return ""
	}

	tokenString = tokenString[len("Bearer "):]
	fmt.Println(tokenString)
	return tokenString
}
