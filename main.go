package main

import (
	"context"
	"log"

	"github.com/antklim/go-dynamodb/dynamo"
	"github.com/antklim/go-dynamodb/invoice"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String("ap-southeast-2")}))
	client := dynamodb.New(sess)
	repo := dynamo.NewRepository(client)
	service := invoice.NewService(repo)

	ctx := context.Background()
	inv, err := service.GetInvoice(ctx, "123")
	log.Println(inv, err)
}
