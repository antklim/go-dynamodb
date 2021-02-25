package invoice

import "context"

// Repository ...
type Repository interface {
	AddInvoice(context.Context, Invoice) error
	GetInvoice(context.Context, string) (*Invoice, error) // gets invoice and all its items
	CancelInvoice(context.Context, string) error          // cancels invoice and all its items status
	AddItem(context.Context, string, Item) error          // adds item to invoice
	GetItem(context.Context, string) (*Item, error)
	GetItemsByStatus(context.Context, Status) ([]Item, error)
	GetInvoiceItemsByStatus(context.Context, string, Status) ([]Item, error)
	UpdateInvoiceItemsStatus(context.Context, string, Status) error
	ReplaceItems(context.Context, string, []Item) error // cancells all invoice items and adds new items
}
