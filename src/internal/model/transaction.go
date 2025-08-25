package model

import (
	"fmt"
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
	PK          string          `dynamodbav:"PK"`
	SK          string          `dynamodbav:"SK"`
	ID          string          `json:"id" dynamodbav:"id"`
	Type        TransactionType `json:"type" dynamodbav:"type"`
	Title       string          `json:"title" dynamodbav:"title"`
	Description *string         `json:"description,omitempty" dynamodbav:"description,omitempty"`
	Amount      int             `json:"amount" dynamodbav:"amount"`
	CreatedAt   time.Time       `json:"created_at" dynamodbav:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" dynamodbav:"updated_at"`
}

func (t *Transaction) GetKey() map[string]types.AttributeValue {
    return map[string]types.AttributeValue{
        "PK": &types.AttributeValueMemberS{Value: t.PK},
        "SK": &types.AttributeValueMemberS{Value: t.SK},
    }
}

func (e *Transaction) SetKeys() {
	e.PK = "TRANSACTION"
	e.SK = fmt.Sprintf("CREATEDAT#%s#%s", e.CreatedAt.Format("2006-01-02"), e.ID)
}
