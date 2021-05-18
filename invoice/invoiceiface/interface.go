package invoiceiface

import (
	"context"

	"github.com/antklim/go-dynamodb/invoice"
)

type Service interface {
	StoreInvoice(context.Context, invoice.Invoice) error
	GetInvoice(context.Context, string) (*invoice.Invoice, error) // gets invoice and all its items
	CancelInvoice(context.Context, string) error                  // cancels invoice and all its items
	AddItem(context.Context, invoice.Item) error                  // adds invoice's item
	GetItem(ctx context.Context, invoiceID, itemID string) (*invoice.Item, error)
	DeleteItem(ctx context.Context, invoiceID, itemID string) error
	GetItemProduct(ctx context.Context, invoiceID, itemID string) (*invoice.Product, error)
	GetItemsByStatus(context.Context, invoice.Status) ([]invoice.Item, error)
	GetInvoiceItemsByStatus(context.Context, string, invoice.Status) ([]invoice.Item, error)
	UpdateInvoiceItemsStatus(context.Context, string, invoice.Status) error
	ReplaceItems(context.Context, string, []invoice.Item) error            // cancells all invoice items and adds new items
	CancelInvoiceItem(ctx context.Context, invoiceID, itemID string) error // cancells invoice item
}
