package main

import (
	"context"
	"log"
	"time"

	"github.com/antklim/go-dynamodb/dynamo"
	"github.com/antklim/go-dynamodb/invoice"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
)

func main() {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String("ap-southeast-2")}))
	client := dynamodb.New(sess)
	repo := dynamo.NewRepository(client, "invoices")
	service := invoice.NewService(repo)

	{
		// 1. Create invoice
		now := time.Now()
		inv := invoice.Invoice{
			ID:           uuid.New().String(),
			Number:       "123",
			CustomerName: "John Doe",
			Status:       "NEW",
			Date:         now,
			Items:        nil,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		ctx := context.Background()
		err := service.StoreInvoice(ctx, inv)
		log.Println(inv)
		log.Println(err)
	}

	// {
	// 	ctx := context.Background()
	// 	inv, err := service.GetInvoice(ctx, "123")
	// 	log.Println(inv, err)
	// }
}
