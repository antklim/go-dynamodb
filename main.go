package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/antklim/go-dynamodb/dynamo"
	"github.com/antklim/go-dynamodb/invoice"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
)

// TODO: Clean DB before and after script run
// TODO: Add flags to control DB clean

func main() {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String("ap-southeast-2")}))
	client := dynamodb.New(sess)
	repo := dynamo.NewRepository(client, "invoices")
	service := invoice.NewService(repo)

	/* Load invoice and items from JSON */
	// invs, err := getInvoices()
	// if err != nil {
	// 	log.Panic(err)
	// }

	// items, err := getInvoiceItems()
	// if err != nil {
	// 	log.Panic(err)
	// }

	// inv := invs[0]
	// {
	// 	// 1. Store invoice
	// 	log.Println(inv)
	// 	ctx := context.Background()
	// 	err := service.StoreInvoice(ctx, inv)
	// 	log.Println(err)
	// }

	// {
	// 	// 2. Replace all invoce's items
	// 	log.Println(items)
	// 	for _, item := range items {
	// 		ctx := context.Background()
	// 		err := service.AddItem(ctx, item)
	// 		if err != nil {
	// 			log.Println(err)
	// 		}
	// 	}
	// }
	/********************************/

	/* Create new invoice and items */
	inv := createInvoice()
	{
		// 1. Store invoice
		log.Println("1. Store invoice =================================")
		ctx := context.Background()
		err := service.StoreInvoice(ctx, inv)
		log.Println(inv)
		log.Println(err)
	}

	{
		// 2. Get all new items
		log.Println("2. Get all new items =============================")
		ctx := context.Background()
		items, err := service.GetItemsByStatus(ctx, "NEW")
		log.Printf("%+v\n", items)
		log.Println(err)
	}

	{
		// 3. Get NEW items of the invoice
		log.Println("3. Get NEW items of the invoice ==================")
		ctx := context.Background()
		items, err := service.GetInvoiceItemsByStatus(ctx, inv.ID, "NEW")
		log.Printf("%+v\n", items)
		log.Println(err)
	}

	{
		// 4. Update all invoce's items status
		log.Println("4. Update all invoce's items status ==============")
		ctx := context.Background()
		err := service.UpdateInvoiceItemsStatus(ctx, inv.ID, "CANCELLED")
		log.Println(err)
	}

	{
		// 5. Replace all invoce's items in status NEW with the new set of items
		log.Println("5. Replace all NEW items =========================")
		ctx := context.Background()
		now := time.Now()
		newItems := []invoice.Item{
			{
				ID:        uuid.NewString(),
				InvoiceID: inv.ID,
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
				InvoiceID: inv.ID,
				SKU:       "301",
				Name:      "Drum sticks",
				Price:     4200,
				Qty:       2,
				Status:    "NEW",
				CreatedAt: now,
				UpdatedAt: now,
			},
		}
		err := service.ReplaceItems(ctx, inv.ID, newItems)
		log.Println(err)
	}

	{
		// 6. Get invoice
		log.Println("6. Get invoice ===================================")
		ctx := context.Background()
		inv, err := service.GetInvoice(ctx, inv.ID)
		log.Println("invoiceID", inv.ID)
		log.Println(err)
		log.Println("invoice", inv)
	}

	{
		// 7. Get item
		log.Println("6. Get item ======================================")
		ctx := context.Background()
		item, err := service.GetItem(ctx, inv.ID, inv.Items[0].ID)
		log.Println("item", item)
		log.Println(err)
	}

	{
		// 8. Get item product
		log.Println("6. Get product ===================================")
		ctx := context.Background()
		product, err := service.GetItemProduct(ctx, inv.ID, inv.Items[0].ID)
		log.Println("product", product)
		log.Println(err)
	}
}

func getInvoices() ([]invoice.Invoice, error) {
	raw, err := ioutil.ReadFile("./data/invoices.json")
	if err != nil {
		return nil, err
	}

	var invoices []invoice.Invoice
	err = json.Unmarshal(raw, &invoices)
	return invoices, err
}

func getInvoiceItems() ([]invoice.Item, error) {
	raw, err := ioutil.ReadFile("./data/items.json")
	if err != nil {
		return nil, err
	}

	var items []invoice.Item
	err = json.Unmarshal(raw, &items)
	return items, err
}

func createInvoice() invoice.Invoice {
	invoiceID := uuid.NewString()
	now := time.Now()
	inv := invoice.Invoice{
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

	return inv
}
