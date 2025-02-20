package jwt

import (
	"context"
	"fmt"
	// "os"
	"time"

	"lambda-dropin/internal/model"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/golang-jwt/jwt/v5"
)

func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	secret, err := GetJWTSecret()
	if err != nil {
		return jwt.MapClaims{}, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("token of unrecognized type - unauthorized")
	}

	return claims, nil
}

func CreateToken(user model.User) (string, error) {
	now := time.Now()
	validUntil := now.Add(time.Hour * 1).Unix()
	secret, err := GetJWTSecret()

	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"user":    user.Username,
		"expires": validUntil,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims, nil)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func GetJWTSecret() (string, error) {
	secretName := "lambdadropin/jwt-secret"
	region := "eu-north-1"

	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return "", err
	}

	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		return "", err
	}

	var secretString string = *result.SecretString
	return secretString, nil

}
