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

func (i *invoices) create(inv invoice.Invoice) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if i.table == nil {
		i.table = make(map[string]invoice.Invoice)
	}

	i.table[inv.ID] = inv
	return nil
}

func (i *invoices) get(invoiceID string) (*invoice.Invoice, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if inv, ok := i.table[invoiceID]; ok {
		return &inv, nil
	}
	return nil, nil
}

type items struct {
	mu    sync.Mutex
	table map[string]invoice.Item
}

func (i *items) create(item invoice.Item) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if i.table == nil {
		i.table = make(map[string]invoice.Item)
	}

	i.table[item.ID] = item
	return nil
}

func (i *items) get(itemID string) (*invoice.Item, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if item, ok := i.table[itemID]; ok {
		return &item, nil
	}
	return nil, nil
}

func (i *items) del(itemID string) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	delete(i.table, itemID)
	return nil
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
	return r.invs.create(inv)
}

func (r *Repository) GetInvoice(ctx context.Context, invoiceID string) (*invoice.Invoice, error) {
	return r.invs.get(invoiceID)
}

func (r *Repository) AddItem(ctx context.Context, item invoice.Item) error {
	return r.itms.create(item)
}

func (r *Repository) GetItem(ctx context.Context, invoiceID, itemID string) (*invoice.Item, error) {
	return r.itms.get(itemID)
}

func (r *Repository) GetItemProduct(ctx context.Context, invoiceID, itemID string) (*invoice.Product, error) {
	item, err := r.itms.get(itemID)
	if item == nil || err != nil {
		return nil, err
	}

	p := invoice.Product{
		SKU:   item.SKU,
		Name:  item.Name,
		Price: item.Price,
	}

	return &p, nil
}

func (r *Repository) DeleteItem(ctx context.Context, invoiceID, itemID string) error {
	return r.itms.del(itemID)
}

func (r *Repository) GetItemsByStatus(ctx context.Context, status invoice.Status) ([]invoice.Item, error) {

	return nil, errNotImplemented
}

func (r *Repository) GetInvoiceItems(ctx context.Context, invoiceID string) ([]invoice.Item, error) {
	return nil, errNotImplemented
}

func (r *Repository) GetInvoiceItemsByStatus(
	ctx context.Context, invoiceID string, status invoice.Status) ([]invoice.Item, error) {

	return nil, errNotImplemented
}

func (r *Repository) UpdateInvoiceItemStatus(
	ctx context.Context, invoiceID, itemID string, status invoice.Status) error {

	return errNotImplemented
}

func (r *Repository) UpdateInvoiceItemsStatus(
	ctx context.Context, invoiceID string, itemIDs []string, status invoice.Status) error {

	return errNotImplemented
}

func (r *Repository) ReplaceItems(ctx context.Context, invoiceID string, newItems []invoice.Item) error {
	return errNotImplemented
}
