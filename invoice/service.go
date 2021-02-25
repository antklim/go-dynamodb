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

func (s *service) GetInvoice(context.Context, string) (*Invoice, error) {
	return nil, errors.New("not implemented")
}

func (s *service) CancelInvoice(context.Context, string) error {
	return errors.New("not implemented")
}

func (s *service) AddItem(context.Context, string, Item) error {
	return errors.New("not implemented")
}

func (s *service) GetItem(context.Context, string) (*Item, error) {
	return nil, errors.New("not implemented")
}

func (s *service) GetItemsByStatus(context.Context, Status) ([]Item, error) {
	return nil, errors.New("not implemented")
}

func (s *service) GetInvoiceItemsByStatus(context.Context, string, Status) ([]Item, error) {
	return nil, errors.New("not implemented")
}

func (s *service) UpdateInvoiceItemsStatus(context.Context, string, Status) error {
	return errors.New("not implemented")
}

func (s *service) ReplaceItems(context.Context, string, []Item) error {
	return errors.New("not implemented")
}
