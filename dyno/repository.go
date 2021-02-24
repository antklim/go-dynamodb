package dyno

import (
	"github.com/antklim/go-dynamodb/invoice"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type repository struct {
	client dynamodbiface.DynamoDBAPI
}

func NewRepository(client dynamodbiface.DynamoDBAPI) invoice.Repository {
	return &repository{client: client}
}
