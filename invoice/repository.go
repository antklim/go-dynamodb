package invoice

// Repository ...
type Repository interface {
	AddInvoice(Invoice) error
	GetInvoice(string) (*Invoice, error) // gets invoice and all its items
	CancelInvoice(string) error          // cancels invoice and all its items status
	AddItem(string, Item) error          // adds item to invoice
	GetItem(string) (*Item, error)
	GetItemsByStatus(Status) ([]*Item, error)
	GetInvoiceItemsByStatus(string, Status) ([]*Item, error)
	UpdateInvoiceItemsStatus(string, Status) error
	ReplaceItems(string, []Item) error // cancells all invoice items and adds new items
}
