package model

import (
	"time"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type User struct {
	ID string `json:"id" dynamodbav:"id"`
	Name string `json:"name" dynamodbav:"name"`
	Email string `json:"email" dynamodbav:"email"`
	Password string `json:"password" dynamodbav:"password"`
	CreatedAt   time.Time `json:"created_at" dynamodbav:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" dynamodbav:"updated_at"`
}


func (e User) GetKey() (map[string]types.AttributeValue) {
	return map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{Value: e.ID},
	}
}