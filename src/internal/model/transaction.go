package model

import (
	"time"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type TransactionType string

const (
	TransactionTypePurchase TransactionType = "PURCHASE"
	TransactionTypeIncome TransactionType = "INCOME"
	TransactionTypeInvestment TransactionType = "INVESTMENT"
)

type Transaction struct {
	ID          string    `json:"id" dynamodbav:"id"`
	Type        TransactionType `json:"type" dynamodbav:"type"`
	Description string    `json:"description" dynamodbav:"description"`
	CreatedAt   time.Time `json:"created_at" dynamodbav:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" dynamodbav:"updated_at"`
} 

func (e Transaction) GetKey() (map[string]types.AttributeValue) {
	return map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{Value: e.ID},
	}
}