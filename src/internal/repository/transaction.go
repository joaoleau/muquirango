package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/joaoleau/muquirango/internal/model"
)

type TransactionRepo struct {
	db        *dynamodb.Client
	tableName string
}

func NewTransactionRepository(db *dynamodb.Client, tableName string) *TransactionRepo {
	return &TransactionRepo{
		db:        db,
		tableName: tableName,
	}
}

func (r *TransactionRepo) NewTransaction(ctx context.Context, Transaction *model.Transaction) (*model.Transaction, error) {
	item, err := attributevalue.MarshalMap(Transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Transaction: %w", err)
	}

	_, err = r.db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add Transaction: %w", err)
	}

	return Transaction, nil
}

func (r *TransactionRepo) UpdateTransaction(ctx context.Context, Transaction *model.Transaction) (*model.Transaction, error) {
	var updateTransaction *model.Transaction

	update := expression.Set(expression.Name("type"), expression.Value(Transaction.Type))
	update = update.Set(expression.Name("description"), expression.Value(Transaction.Description))

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build update expression: %w", err)
	}

	response, err := r.db.UpdateItem(
		ctx,
		&dynamodb.UpdateItemInput{
			TableName:                 &r.tableName,
			Key:                       Transaction.GetKey(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
			ReturnValues:              types.ReturnValueUpdatedNew,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update Transaction with ID '%s': %w", Transaction.ID, err)
	}

	err = attributevalue.UnmarshalMap(response.Attributes, &updateTransaction)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal updated Transaction: %w", err)
	}

	return updateTransaction, nil
}

func (r *TransactionRepo) ListTransactions(ctx context.Context) (*[]model.Transaction, error) {
	var entries *[]model.Transaction
	var allItems []map[string]types.AttributeValue
	var lastEvaluatedKey map[string]types.AttributeValue
	var err error

	for {
		input := &dynamodb.ScanInput{
			TableName:         aws.String(r.tableName),
			ExclusiveStartKey: lastEvaluatedKey,
		}

		response, err := r.db.Scan(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to scan entries from table: %w", err)
		}
		allItems = append(allItems, response.Items...)

		if response.LastEvaluatedKey == nil {
			break
		}
		lastEvaluatedKey = response.LastEvaluatedKey
	}

	err = attributevalue.UnmarshalListOfMaps(allItems, &entries)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal entries list: %w", err)
	}

	return entries, nil
}

func (r *TransactionRepo) GetTransactionByID(ctx context.Context, id string) (*model.Transaction, error) {
	result := &model.Transaction{ID: id}

	response, err := r.db.GetItem(
		ctx,
		&dynamodb.GetItemInput{
			TableName: aws.String(r.tableName),
			Key:       result.GetKey(),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get Transaction with ID '%s': %w", id, err)
	}

	if len(response.Item) == 0 {
		return nil, fmt.Errorf("Transaction with ID '%s' not found", id)
	}

	err = attributevalue.UnmarshalMap(response.Item, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal Transaction with ID '%s': %w", id, err)
	}

	return result, nil
}

func (r *TransactionRepo) DeleteTransaction(ctx context.Context, Transaction *model.Transaction) (*model.Transaction, error) {
	var deleteTransaction *model.Transaction

	response, err := r.db.DeleteItem(
		ctx,
		&dynamodb.DeleteItemInput{
			TableName: aws.String(r.tableName),
			Key:       Transaction.GetKey(),
			ReturnValues: types.ReturnValueAllOld,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to delete Transaction with ID '%s': %w", Transaction.ID, err)
	}

	if len(response.Attributes) == 0 {
		return nil, fmt.Errorf("no Transaction found to delete with ID '%s'", Transaction.ID)
	}

	err = attributevalue.UnmarshalMap(response.Attributes, &deleteTransaction)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal deleted Transaction: %w", err)
	}

	return deleteTransaction, nil
}
