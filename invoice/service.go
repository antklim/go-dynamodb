package invoice

import (
	"context"
	"errors"
	"time"
)

// Status ...
type Status string

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

type Service interface {
	StoreInvoice(context.Context, Invoice) error
	GetInvoice(context.Context, string) (*Invoice, error) // gets invoice and all its items
	CancelInvoice(context.Context, string) error          // cancels invoice and all its items
	AddItem(context.Context, string, Item) error          // adds item to invoice
	GetItem(context.Context, string) (*Item, error)
	GetItemsByStatus(context.Context, Status) ([]Item, error)
	GetInvoiceItemsByStatus(context.Context, string, Status) ([]Item, error)
	UpdateInvoiceItemsStatus(context.Context, string, Status) error
	ReplaceItems(context.Context, string, []Item) error // cancells all invoice items and adds new items
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

func (s *service) AddItem(context.Context, string, Item) error {
	// TODO: implement
	return errors.New("not implemented")
}

func (s *service) GetItem(context.Context, string) (*Item, error) {
	// TODO: implement
	return nil, errors.New("not implemented")
}

func (s *service) GetItemsByStatus(ctx context.Context, status Status) ([]Item, error) {
	return s.repo.GetItemsByStatus(ctx, status)
}

func (s *service) GetInvoiceItemsByStatus(ctx context.Context, invoiceID string, status Status) ([]Item, error) {
	return s.repo.GetInvoiceItemsByStatus(ctx, invoiceID, status)
}

func (s *service) UpdateInvoiceItemsStatus(ctx context.Context, invoiceID string, status Status) error {
	return s.repo.UpdateInvoiceItemsStatus(ctx, invoiceID, status)
}

func (s *service) ReplaceItems(ctx context.Context, invoiceID string, newItems []Item) error {
	return s.repo.ReplaceItems(ctx, invoiceID, newItems)
}
