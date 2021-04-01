package invoice

import (
	"context"
	"errors"
	"time"
)

// Status ...
type Status string

const (
	New       Status = "NEW"
	Pending   Status = "PENDING"
	Cancelled Status = "CANCELLED"
)

// Invoice ...
type Invoice struct {
	ID           string // unique identifier, uuid format
	Number       string // sequential invoice number
	CustomerName string
	Status       Status
	Date         time.Time
	Items        []Item
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Item ...
type Item struct {
	ID        string // unique identifier, uuid format
	InvoiceID string
	SKU       string
	Name      string
	Price     uint
	Qty       uint
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Product ...
type Product struct {
	SKU   string
	Name  string
	Price uint
}

type Service interface {
	StoreInvoice(context.Context, Invoice) error
	GetInvoice(context.Context, string) (*Invoice, error) // gets invoice and all its items
	CancelInvoice(context.Context, string) error          // cancels invoice and all its items
	AddItem(context.Context, Item) error                  // adds invoice's item
	GetItem(ctx context.Context, invoiceID, itemID string) (*Item, error)
	GetItemProduct(ctx context.Context, invoiceID, itemID string) (*Product, error)
	GetItemsByStatus(context.Context, Status) ([]Item, error)
	GetInvoiceItemsByStatus(context.Context, string, Status) ([]Item, error)
	UpdateInvoiceItemsStatus(context.Context, string, Status) error
	ReplaceItems(context.Context, string, []Item) error                    // cancells all invoice items and adds new items
	CancelInvoiceItem(ctx context.Context, invoiceID, itemID string) error // cancells invoice item
}

type service struct {
	repo Repository
}

// NewService creates a new instance of invoice service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) StoreInvoice(ctx context.Context, inv Invoice) error {
	return s.repo.AddInvoice(ctx, inv)
}

func (s *service) GetInvoice(ctx context.Context, invoiceID string) (*Invoice, error) {
	return s.repo.GetInvoice(ctx, invoiceID)
}

func (s *service) CancelInvoice(context.Context, string) error {
	// TODO: implement
	return errors.New("not implemented")
}

func (s *service) AddItem(ctx context.Context, item Item) error {
	return s.repo.AddItem(ctx, item)
}

func (s *service) GetItem(ctx context.Context, invoiceID, itemID string) (*Item, error) {
	return s.repo.GetItem(ctx, invoiceID, itemID)
}

func (s *service) GetItemProduct(ctx context.Context, invoiceID, itemID string) (*Product, error) {
	return s.repo.GetItemProduct(ctx, invoiceID, itemID)
}

func (s *service) GetItemsByStatus(ctx context.Context, status Status) ([]Item, error) {
	return s.repo.GetItemsByStatus(ctx, status)
}

func (s *service) GetInvoiceItemsByStatus(ctx context.Context, invoiceID string, status Status) ([]Item, error) {
	return s.repo.GetInvoiceItemsByStatus(ctx, invoiceID, status)
}

func (s *service) UpdateInvoiceItemsStatus(ctx context.Context, invoiceID string, status Status) error {
	items, err := s.repo.GetInvoiceItems(ctx, invoiceID)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		return nil
	}

	itemIDs := make([]string, len(items))
	for idx, item := range items {
		itemIDs[idx] = item.ID
	}

	return s.repo.UpdateInvoiceItemsStatus(ctx, invoiceID, itemIDs, status)
}

func (s *service) ReplaceItems(ctx context.Context, invoiceID string, newItems []Item) error {
	return s.repo.ReplaceItems(ctx, invoiceID, newItems)
}

func (s *service) CancelInvoiceItem(ctx context.Context, invoiceID, itemID string) error {
	return s.repo.UpdateInvoiceItemStatus(ctx, invoiceID, itemID, Cancelled)
}
