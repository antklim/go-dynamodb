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

	var inv invoice.Invoice
	{
		// 1. Create invoice
		now := time.Now()
		inv = invoice.Invoice{
			ID:           uuid.NewString(),
			Number:       "123",
			CustomerName: "John Doe",
			Status:       "NEW",
			Date:         now,
			Items: []invoice.Item{
				{
					ID:        uuid.NewString(),
					SKU:       "100",
					Name:      "Guitar",
					Price:     75000,
					Qty:       1,
					Status:    "NEW",
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        uuid.NewString(),
					SKU:       "101",
					Name:      "Guitar strings",
					Price:     8300,
					Qty:       3,
					Status:    "PENDING",
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        uuid.NewString(),
					SKU:       "102",
					Name:      "Pick",
					Price:     1000,
					Qty:       2,
					Status:    "NEW",
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			CreatedAt: now,
			UpdatedAt: now,
		}
		ctx := context.Background()
		err := service.StoreInvoice(ctx, inv)
		log.Println(inv)
		log.Println(err)
	}

	{
		// 2. Get all new items
		ctx := context.Background()
		items, err := service.GetItemsByStatus(ctx, "NEW")
		log.Printf("%+v\n", items)
		log.Println(err)
	}

	{
		// 3. Get new items of the invoice
		ctx := context.Background()
		items, err := service.GetInvoiceItemsByStatus(ctx, inv.ID, "NEW")
		log.Printf("%+v\n", items)
		log.Println(err)
	}
}
