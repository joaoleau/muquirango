package dto

type CreateUserInput struct {
	Name string `json:"name" dynamodbav:"name"`
	Email string `json:"email" dynamodbav:"email"`
	Password string `json:"password" dynamodbav:"password"`
}