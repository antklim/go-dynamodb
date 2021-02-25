package dynamo

import (
	"errors"
	"time"

	"github.com/antklim/go-dynamodb/invoice"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// Invoice describes dynamodb representation of invoice.Invoice
type Invoice struct {
	PK           string    `dynamodbav:"pk"`
	SK           string    `dynamodbav:"sk"`
	ID           string    `dynamodbav:"id"`
	Number       string    `dynamodbav:"number"`
	CustomerName string    `dynamodbav:"customerName"`
	Status       string    `dynamodbav:"status"`
	Date         string    `dynamodbav:"date"` // YYYYMMDD
	CreatedAt    time.Time `dynamodbav:"createdAt"`
	UpdatedAt    time.Time `dynamodbav:"updatedAt"`
}

// NewInvoice creates an instance of DynamoDB invoice from invoice.Invoice.
func NewInvoice(inv invoice.Invoice) *Invoice {
	return nil
}

// ToInvoice creates an instance of invoice.Invoice from DynamoDB invoice.
func (inv *Invoice) ToInvoice() (*invoice.Invoice, error) {
	return nil, nil
}

// Item describes dynamodb representation of invoice.Item
type Item struct {
	PK        string    `dynamodbav:"pk"`
	SK        string    `dynamodbav:"sk"`
	ID        string    `dynamodbav:"id"`
	SKU       string    `dynamodbav:"sku"`
	Name      string    `dynamodbav:"name"`
	Price     uint      `dynamodbav:"price"`
	Qty       uint      `dynamodbav:"qty"`
	Status    string    `dynamodbav:"createdAt"`
	CreatedAt time.Time `dynamodbav:"createdAt"`
	UpdatedAt time.Time `dynamodbav:"updatedAt"`
}

// NewItem creates an instance of DynamoDB item from invoice.Item.
func NewItem(item invoice.Invoice) *Item {
	return nil
}

// ToItem creates an instance of invoice.Item from DynamoDB item.
func (item *Item) ToItem() (*invoice.Item, error) {
	return nil, nil
}

type repository struct {
	client dynamodbiface.DynamoDBAPI
}

// NewRepository ...
func NewRepository(client dynamodbiface.DynamoDBAPI) invoice.Repository {
	return &repository{client: client}
}

func (r *repository) AddInvoice(invoice.Invoice) error {
	return errors.New("not implemented")
}

func (r *repository) GetInvoice(string) (*invoice.Invoice, error) {
	return nil, errors.New("not implemented")
}

func (r *repository) CancelInvoice(string) error {
	return errors.New("not implemented")
}

func (r *repository) AddItem(string, invoice.Item) error {
	return errors.New("not implemented")
}

func (r *repository) GetItem(string) (*invoice.Item, error) {
	return nil, errors.New("not implemented")
}

func (r *repository) GetItemsByStatus(invoice.Status) ([]*invoice.Item, error) {
	return nil, errors.New("not implemented")
}

func (r *repository) GetInvoiceItemsByStatus(string, invoice.Status) ([]*invoice.Item, error) {
	return nil, errors.New("not implemented")
}

func (r *repository) UpdateInvoiceItemsStatus(string, invoice.Status) error {
	return errors.New("not implemented")
}

func (r *repository) ReplaceItems(string, []invoice.Item) error {
	return errors.New("not implemented")
}
