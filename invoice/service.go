package invoice

import "time"

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

func NewService(repo Repository) Service {
	return &service{repo: repo}
}
