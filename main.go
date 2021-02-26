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
	invoiceID := uuid.NewString()
	{
		// 1. Create invoice
		now := time.Now()
		inv = invoice.Invoice{
			ID:           invoiceID,
			Number:       "123",
			CustomerName: "John Doe",
			Status:       "NEW",
			Date:         now,
			Items: []invoice.Item{
				{
					ID:        uuid.NewString(),
					InvoiceID: invoiceID,
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
					InvoiceID: invoiceID,
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
					InvoiceID: invoiceID,
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
		// 3. Get NEW items of the invoice
		ctx := context.Background()
		items, err := service.GetInvoiceItemsByStatus(ctx, invoiceID, "NEW")
		log.Printf("%+v\n", items)
		log.Println(err)
	}

	{
		// 4. Update all invoce's items status
		ctx := context.Background()
		err := service.UpdateInvoiceItemsStatus(ctx, invoiceID, "CANCELLED")
		log.Println(err)
	}

	{
		// 5. Replace all invoce's items in status NEW with the new set of items
		ctx := context.Background()
		now := time.Now()
		newItems := []invoice.Item{
			{
				ID:        uuid.NewString(),
				InvoiceID: invoiceID,
				SKU:       "300",
				Name:      "Drums set",
				Price:     132000,
				Qty:       1,
				Status:    "NEW",
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        uuid.NewString(),
				InvoiceID: invoiceID,
				SKU:       "301",
				Name:      "Drum sticks",
				Price:     4200,
				Qty:       2,
				Status:    "NEW",
				CreatedAt: now,
				UpdatedAt: now,
			},
		}
		err := service.ReplaceItems(ctx, invoiceID, newItems)
		log.Println(err)
	}
}
