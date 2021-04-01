package invoice

import "context"

// Repository interface defines invoces repository methods
type Repository interface {
	AddInvoice(context.Context, Invoice) error
	GetInvoice(context.Context, string) (*Invoice, error) // gets invoice and all its items
	AddItem(context.Context, Item) error                  // adds invoice's item
	GetItem(ctx context.Context, invoiceID, itemID string) (*Item, error)
	GetItemProduct(ctx context.Context, invoiceID, itemID string) (*Product, error)
	GetItemsByStatus(context.Context, Status) ([]Item, error)
	// TODO: implement
	// GetInvoiceItems(context.Context, string) ([]Item, error)
	GetInvoiceItemsByStatus(context.Context, string, Status) ([]Item, error)
	UpdateInvoiceItemStatus(ctx context.Context, invoiceID, itemID string, status Status) error
	UpdateInvoiceItemsStatus(ctx context.Context, invoiceID string, itemIDs []string, status Status) error
	ReplaceItems(context.Context, string, []Item) error // cancells all invoice items and adds new items
}
