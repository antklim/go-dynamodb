package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/antklim/go-dynamodb/invoice"
)

type invoices struct {
	mu    sync.Mutex
	table map[string]invoice.Invoice
}

type items struct {
	mu    sync.Mutex
	table map[string]invoice.Item
}

type Repository struct {
	invs invoices
	itms items
}

var errNotImplemented = errors.New("not implemented")

// NewRepository creates in memory implementation of the repository
func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) AddInvoice(ctx context.Context, inv invoice.Invoice) error {
	return errNotImplemented
}

func (r *Repository) GetInvoice(ctx context.Context, invoiceID string) (*invoice.Invoice, error) {
	return nil, errNotImplemented
}

func (r *Repository) AddItem(ctx context.Context, item invoice.Item) error {
	return errNotImplemented
}

func (r *Repository) GetItem(ctx context.Context, invoiceID, itemID string) (*invoice.Item, error) {
	return nil, errNotImplemented
}

func (r *Repository) GetItemProduct(ctx context.Context, invoiceID, itemID string) (*invoice.Product, error) {
	return nil, errNotImplemented
}

func (r *Repository) DeleteItem(ctx context.Context, invoiceID, itemID string) error {
	return errNotImplemented
}

func (r *Repository) GetItemsByStatus(context.Context, invoice.Status) ([]invoice.Item, error) {
	return nil, errNotImplemented
}

func (r *Repository) GetInvoiceItems(context.Context, string) ([]invoice.Item, error) {
	return nil, errNotImplemented
}

func (r *Repository) GetInvoiceItemsByStatus(context.Context, string, invoice.Status) ([]invoice.Item, error) {
	return nil, errNotImplemented
}

func (r *Repository) UpdateInvoiceItemStatus(ctx context.Context, invoiceID, itemID string, status invoice.Status) error {
	return errNotImplemented
}

func (r *Repository) UpdateInvoiceItemsStatus(
	ctx context.Context, invoiceID string, itemIDs []string, status invoice.Status) error {
	return errNotImplemented
}

func (r *Repository) ReplaceItems(context.Context, string, []invoice.Item) error {
	return errNotImplemented
}
