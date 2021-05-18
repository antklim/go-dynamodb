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

type Service struct {
	repo Repository
}

// NewService creates a new instance of invoice service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) StoreInvoice(ctx context.Context, inv Invoice) error {
	return s.repo.AddInvoice(ctx, inv)
}

func (s *Service) GetInvoice(ctx context.Context, invoiceID string) (*Invoice, error) {
	return s.repo.GetInvoice(ctx, invoiceID)
}

func (s *Service) CancelInvoice(context.Context, string) error {
	// TODO: implement
	return errors.New("not implemented")
}

func (s *Service) AddItem(ctx context.Context, item Item) error {
	return s.repo.AddItem(ctx, item)
}

func (s *Service) GetItem(ctx context.Context, invoiceID, itemID string) (*Item, error) {
	return s.repo.GetItem(ctx, invoiceID, itemID)
}

func (s *Service) DeleteItem(ctx context.Context, invoiceID, itemID string) error {
	return s.repo.DeleteItem(ctx, invoiceID, itemID)
}

func (s *Service) GetItemProduct(ctx context.Context, invoiceID, itemID string) (*Product, error) {
	return s.repo.GetItemProduct(ctx, invoiceID, itemID)
}

func (s *Service) GetItemsByStatus(ctx context.Context, status Status) ([]Item, error) {
	return s.repo.GetItemsByStatus(ctx, status)
}

func (s *Service) GetInvoiceItemsByStatus(ctx context.Context, invoiceID string, status Status) ([]Item, error) {
	return s.repo.GetInvoiceItemsByStatus(ctx, invoiceID, status)
}

func (s *Service) UpdateInvoiceItemsStatus(ctx context.Context, invoiceID string, status Status) error {
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

func (s *Service) ReplaceItems(ctx context.Context, invoiceID string, newItems []Item) error {
	return s.repo.ReplaceItems(ctx, invoiceID, newItems)
}

func (s *Service) CancelInvoiceItem(ctx context.Context, invoiceID, itemID string) error {
	return s.repo.UpdateInvoiceItemStatus(ctx, invoiceID, itemID, Cancelled)
}
