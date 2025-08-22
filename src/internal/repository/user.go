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

type UserRepo struct {
	db        *dynamodb.Client
	tableName string
}

func NewUserRepository(db *dynamodb.Client, tableName string) *UserRepo {
	return &UserRepo{
		db:        db,
		tableName: tableName,
	}
}

func (r *UserRepo) Add(ctx context.Context, user *model.User) (*model.User, error) {
	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal User: %w", err)
	}

	_, err = r.db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add User: %w", err)
	}

	return user, nil
}

func (r *UserRepo) UpdateMe(ctx context.Context, user *model.User) (*model.User, error) {
	var updateUser *model.User

	update := expression.Set(expression.Name("name"), expression.Value(user.Name))
	update = update.Set(expression.Name("email"), expression.Value(user.Email))

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build update expression: %w", err)
	}

	response, err := r.db.UpdateItem(
		ctx,
		&dynamodb.UpdateItemInput{
			TableName:                 &r.tableName,
			Key:                       user.GetKey(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
			ReturnValues:              types.ReturnValueUpdatedNew,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update User with ID '%s': %w", user.ID, err)
	}

	err = attributevalue.UnmarshalMap(response.Attributes, &updateUser)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal updated User: %w", err)
	}

	return updateUser, nil
}

func (r *UserRepo) GetMe(ctx context.Context, id string) (*model.User, error) {
	result := &model.User{ID: id}

	response, err := r.db.GetItem(
		ctx,
		&dynamodb.GetItemInput{
			TableName: aws.String(r.tableName),
			Key:       result.GetKey(),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get User with ID '%s': %w", id, err)
	}

	if len(response.Item) == 0 {
		return nil, fmt.Errorf("User with ID '%s' not found", id)
	}

	err = attributevalue.UnmarshalMap(response.Item, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal User with ID '%s': %w", id, err)
	}

	return result, nil
}

func (r *UserRepo) DeleteMe(ctx context.Context, User *model.User) (*model.User, error) {
	var deleteUser *model.User

	response, err := r.db.DeleteItem(
		ctx,
		&dynamodb.DeleteItemInput{
			TableName: aws.String(r.tableName),
			Key:       User.GetKey(),
			ReturnValues: types.ReturnValueAllOld,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to delete User with ID '%s': %w", User.ID, err)
	}

	if len(response.Attributes) == 0 {
		return nil, fmt.Errorf("no User found to delete with ID '%s'", User.ID)
	}

	err = attributevalue.UnmarshalMap(response.Attributes, &deleteUser)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal deleted User: %w", err)
	}

	return deleteUser, nil
}
