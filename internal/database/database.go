package database

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const USER_TABLE_NAME = "userTable"

func NewDynamoDB() *dynamodb.DynamoDB {
	dbSession := session.Must(session.NewSession())
	db := dynamodb.New(dbSession)

	return db
}
