package invoice

import (
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
	StoreInvoice(Invoice) error
	GetInvoice(string) (*Invoice, error) // gets invoice and all its items
	CancelInvoice(string) error          // cancels invoice and all its items
	AddItem(string, Item) error          // adds item to invoice
	GetItem(string) (*Item, error)
	GetItemsByStatus(Status) ([]*Item, error)
	GetInvoiceItemsByStatus(string, Status) ([]*Item, error)
	UpdateInvoiceItemsStatus(string, Status) error
	ReplaceItems(string, []Item) error // cancells all invoice items and adds new items
}

type service struct {
	repo Repository
}

// NewService creates a new instance of invoice service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) StoreInvoice(Invoice) error {
	return errors.New("not implemented")
}

func (s *service) GetInvoice(string) (*Invoice, error) {
	return nil, errors.New("not implemented")
}

func (s *service) CancelInvoice(string) error {
	return errors.New("not implemented")
}

func (s *service) AddItem(string, Item) error {
	return errors.New("not implemented")
}

func (s *service) GetItem(string) (*Item, error) {
	return nil, errors.New("not implemented")
}

func (s *service) GetItemsByStatus(Status) ([]*Item, error) {
	return nil, errors.New("not implemented")
}

func (s *service) GetInvoiceItemsByStatus(string, Status) ([]*Item, error) {
	return nil, errors.New("not implemented")
}

func (s *service) UpdateInvoiceItemsStatus(string, Status) error {
	return errors.New("not implemented")
}

func (s *service) ReplaceItems(string, []Item) error {
	return errors.New("not implemented")
}
