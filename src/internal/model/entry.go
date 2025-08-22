package model

import (
	"time"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type EntryType string

const (
	EntryTypePurchase EntryType = "PURCHASE"
	EntryTypeIncome EntryType = "INCOME"
	EntryTypeInvestment EntryType = "INVESTMENT"
)

type Entry struct {
	ID          string    `json:"id" dynamodbav:"id"`
	UserID      string    `json:"user_id" dynamodbav:"user_id"`
	Type        EntryType `json:"type" dynamodbav:"type"`
	Description string    `json:"description" dynamodbav:"description"`
	CreatedAt   time.Time `json:"created_at" dynamodbav:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" dynamodbav:"updated_at"`
} 

func (e Entry) GetKey() (map[string]types.AttributeValue) {
	return map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{Value: e.ID},
	}
}