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
	"github.com/joaoleau/muquirango/internal/config/logger"
	"go.uber.org/zap"
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

func (r *TransactionRepo) NewTransaction(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error) {
	logger.Info("Attempting to create new transaction", zap.String("transaction_id", transaction.ID))

	item, err := attributevalue.MarshalMap(transaction)
	if err != nil {
		logger.Error("Failed to marshal transaction", err, zap.String("transaction_id", transaction.ID))
		return nil, fmt.Errorf("failed to marshal transaction: %w", err)
	}

	_, err = r.db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	if err != nil {
		logger.Error("Failed to add transaction to DynamoDB", err, zap.String("transaction_id", transaction.ID))
		return nil, fmt.Errorf("failed to add transaction: %w", err)
	}

	logger.Info("Transaction successfully created", zap.String("transaction_id", transaction.ID))
	return transaction, nil
}

func (r *TransactionRepo) UpdateTransaction(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error) {
	logger.Info("Attempting to update transaction", zap.String("transaction_id", transaction.ID))

	var updateTransaction *model.Transaction

	update := expression.Set(expression.Name("type"), expression.Value(transaction.Type))
	update = update.Set(expression.Name("title"), expression.Value(transaction.Title))
	update = update.Set(expression.Name("description"), expression.Value(transaction.Description))
	update = update.Set(expression.Name("amount"), expression.Value(transaction.Amount))

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		logger.Error("Failed to build update expression", err, zap.String("transaction_id", transaction.ID))
		return nil, fmt.Errorf("failed to build update expression: %w", err)
	}

	response, err := r.db.UpdateItem(
		ctx,
		&dynamodb.UpdateItemInput{
			TableName:                 &r.tableName,
			Key:                       transaction.GetKey(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
			ReturnValues:              types.ReturnValueUpdatedNew,
		},
	)
	if err != nil {
		logger.Error("Failed to update transaction in DynamoDB", err, zap.String("transaction_id", transaction.ID))
		return nil, fmt.Errorf("failed to update transaction with ID '%s': %w", transaction.ID, err)
	}

	err = attributevalue.UnmarshalMap(response.Attributes, &updateTransaction)
	if err != nil {
		logger.Error("Failed to unmarshal updated transaction", err, zap.String("transaction_id", transaction.ID))
		return nil, fmt.Errorf("failed to unmarshal updated transaction: %w", err)
	}

	logger.Info("Transaction successfully updated", zap.String("transaction_id", transaction.ID))
	return updateTransaction, nil
}

func (r *TransactionRepo) ListTransactions(ctx context.Context, startDate string, endDate string) (*[]model.Transaction, error) {
    logger.Info("Attempting to list transactions", zap.String("startDate", startDate), zap.String("endDate", endDate))

    startSK := fmt.Sprintf("CREATEDAT#%s#", startDate)
    endSK := fmt.Sprintf("CREATEDAT#%s#z", endDate)

    input := &dynamodb.QueryInput{
        TableName: aws.String(r.tableName),
        KeyConditionExpression: aws.String("PK = :pk AND SK BETWEEN :start AND :end"),
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":pk":    &types.AttributeValueMemberS{Value: "TRANSACTION"},
            ":start": &types.AttributeValueMemberS{Value: startSK},
            ":end":   &types.AttributeValueMemberS{Value: endSK},
        },
    }

    resp, err := r.db.Query(ctx, input)
    if err != nil {
        logger.Error("Failed to query transactions from DynamoDB", err)
        return nil, fmt.Errorf("failed to query entries from table: %w", err)
    }

    var entries []model.Transaction
    if err := attributevalue.UnmarshalListOfMaps(resp.Items, &entries); err != nil {
        logger.Error("Failed to unmarshal transactions list", err)
        return nil, fmt.Errorf("failed to unmarshal entries list: %w", err)
    }

    logger.Info("Transactions successfully retrieved", zap.Int("count", len(entries)))
    return &entries, nil
}

func (r *TransactionRepo) GetTransactionByID(ctx context.Context, id string, createdAt string) (*model.Transaction, error) {
    logger.Info("Attempting to fetch transaction",
        zap.String("transaction_id", id),
        zap.String("created_at", createdAt),
    )

    input := &dynamodb.QueryInput{
        TableName: aws.String(r.tableName),
        KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :id)"),
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":pk": &types.AttributeValueMemberS{Value: "TRANSACTION"},
            ":id": &types.AttributeValueMemberS{Value: fmt.Sprintf("CREATEDAT#%s#%s", createdAt, id)},
        },
        Limit: aws.Int32(1),
    }

    resp, err := r.db.Query(ctx, input)
    if err != nil {
        logger.Error("Failed to fetch transaction from DynamoDB", err,
            zap.String("transaction_id", id),
            zap.String("created_at", createdAt),
        )
        return nil, fmt.Errorf("failed to get transaction with ID '%s': %w", id, err)
    }

    if len(resp.Items) == 0 {
		err := fmt.Errorf("transaction not found: created_at %s && transaction_id %s", createdAt, id)
		logger.Error("Transaction not found", err,
			zap.String("transaction_id", id),
			zap.String("created_at", createdAt),
		)
		return nil, fmt.Errorf("transaction with ID '%s' not found", id)
    }

    var transaction model.Transaction
    if err := attributevalue.UnmarshalMap(resp.Items[0], &transaction); err != nil {
        logger.Error("Failed to unmarshal transaction", err,
            zap.String("transaction_id", id),
            zap.String("created_at", createdAt),
        )
        return nil, fmt.Errorf("failed to unmarshal transaction with ID '%s': %w", id, err)
    }

    logger.Info("Transaction successfully retrieved",
        zap.String("transaction_id", id),
        zap.String("created_at", createdAt),
    )
    return &transaction, nil
}

func (r *TransactionRepo) DeleteTransaction(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error) {
	logger.Info("Attempting to delete transaction", zap.String("transaction_id", transaction.ID))

	var deleteTransaction *model.Transaction

	response, err := r.db.DeleteItem(
		ctx,
		&dynamodb.DeleteItemInput{
			TableName:    aws.String(r.tableName),
			Key:          transaction.GetKey(),
			ReturnValues: types.ReturnValueAllOld,
		},
	)
	if err != nil {
		logger.Error("Failed to delete transaction from DynamoDB", err, zap.String("transaction_id", transaction.ID))
		return nil, fmt.Errorf("failed to delete transaction with ID '%s': %w", transaction.ID, err)
	}

	if len(response.Attributes) == 0 {
		logger.Error("No transaction found to delete", fmt.Errorf("not found"), zap.String("transaction_id", transaction.ID))
		return nil, fmt.Errorf("no transaction found to delete with ID '%s'", transaction.ID)
	}

	err = attributevalue.UnmarshalMap(response.Attributes, &deleteTransaction)
	if err != nil {
		logger.Error("Failed to unmarshal deleted transaction", err, zap.String("transaction_id", transaction.ID))
		return nil, fmt.Errorf("failed to unmarshal deleted transaction: %w", err)
	}

	logger.Info("Transaction successfully deleted", zap.String("transaction_id", transaction.ID))
	return deleteTransaction, nil
}
