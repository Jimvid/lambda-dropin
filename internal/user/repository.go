package user

import (
	"fmt"
	"lambda-dropin/internal/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type UserRepository struct {
	db *dynamodb.DynamoDB
}

const USER_TABLE_NAME = "userTable"

func NewUserRepository(db *dynamodb.DynamoDB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u UserRepository) DoesUserExist(username string) (bool, error) {
	result, err := u.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(USER_TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	})

	if err != nil {
		return true, err

	}

	if result.Item == nil {
		return false, nil
	}

	return true, nil
}

func (u UserRepository) InsertUser(user model.User) error {
	item := &dynamodb.PutItemInput{
		TableName: aws.String(USER_TABLE_NAME),
		Item: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(user.Username),
			},
			"password": {
				S: aws.String(user.PasswordHash),
			},
		},
	}
	_, err := u.db.PutItem(item)
	if err != nil {
		return err
	}

	return nil
}

func (u UserRepository) GetUser(username string) (model.User, error) {
	var user model.User

	result, err := u.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(USER_TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	})

	if err != nil {
		return user, err

	}

	if result.Item == nil {
		return user, fmt.Errorf("user not found")
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		return user, err
	}

	return user, nil
}
