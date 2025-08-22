package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/joaoleau/muquirango/model"
)

type EntryRepo struct {
	db        *dynamodb.Client
	tableName string
}

func NewEntryRepository(db *dynamodb.Client, tableName string) *EntryRepo {
	return &EntryRepo{
		db:        db,
		tableName: tableName,
	}
}

func (r *EntryRepo) Add(ctx context.Context, entry *model.Entry) (*model.Entry, error) {
	item, err := attributevalue.MarshalMap(entry)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal entry: %w", err)
	}

	_, err = r.db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add entry: %w", err)
	}

	return entry, nil
}

func (r *EntryRepo) Update(ctx context.Context, entry *model.Entry) (*model.Entry, error) {
	var updateEntry *model.Entry

	update := expression.Set(expression.Name("type"), expression.Value(entry.Type))
	update = update.Set(expression.Name("description"), expression.Value(entry.Description))

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build update expression: %w", err)
	}

	response, err := r.db.UpdateItem(
		ctx,
		&dynamodb.UpdateItemInput{
			TableName:                 &r.tableName,
			Key:                       entry.GetKey(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
			ReturnValues:              types.ReturnValueUpdatedNew,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update entry with ID '%s': %w", entry.ID, err)
	}

	err = attributevalue.UnmarshalMap(response.Attributes, &updateEntry)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal updated entry: %w", err)
	}

	return updateEntry, nil
}

func (r *EntryRepo) GetAll(ctx context.Context) (*[]model.Entry, error) {
	var entries *[]model.Entry
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

func (r *EntryRepo) GetById(ctx context.Context, id string) (*model.Entry, error) {
	result := &model.Entry{ID: id}

	response, err := r.db.GetItem(
		ctx,
		&dynamodb.GetItemInput{
			TableName: aws.String(r.tableName),
			Key:       result.GetKey(),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get entry with ID '%s': %w", id, err)
	}

	if len(response.Item) == 0 {
		return nil, fmt.Errorf("entry with ID '%s' not found", id)
	}

	err = attributevalue.UnmarshalMap(response.Item, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal entry with ID '%s': %w", id, err)
	}

	return result, nil
}

func (r *EntryRepo) Delete(ctx context.Context, entry *model.Entry) (*model.Entry, error) {
	var deleteEntry *model.Entry

	response, err := r.db.DeleteItem(
		ctx,
		&dynamodb.DeleteItemInput{
			TableName: aws.String(r.tableName),
			Key:       entry.GetKey(),
			ReturnValues: types.ReturnValueAllOld,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to delete entry with ID '%s': %w", entry.ID, err)
	}

	if len(response.Attributes) == 0 {
		return nil, fmt.Errorf("no entry found to delete with ID '%s'", entry.ID)
	}

	err = attributevalue.UnmarshalMap(response.Attributes, &deleteEntry)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal deleted entry: %w", err)
	}

	return deleteEntry, nil
}
