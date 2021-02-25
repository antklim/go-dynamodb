package dynamo

import (
	"errors"
	"strings"
	"time"

	"github.com/antklim/go-dynamodb/invoice"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

const (
	keySeparator    = "#"
	invoicePkPrefix = "INVOICE"
	invoiceSkPrefix = "INVOICE"
	itemPkPrefix    = "INVOICE" // invoice items are in the same partition as the invoice
	itemSkPrefix    = "ITEM"
	yyyymmddFormat  = "2006-01-02"
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
	pk := invoicePk(inv)
	sk := invoiceSk(inv)

	return &Invoice{
		PK:           pk,
		SK:           sk,
		ID:           inv.ID,
		Number:       inv.Number,
		CustomerName: inv.CustomerName,
		Status:       string(inv.Status),
		Date:         inv.Date.Format(yyyymmddFormat),
		CreatedAt:    inv.CreatedAt,
		UpdatedAt:    inv.UpdatedAt,
	}
}

// ToInvoice creates an instance of invoice.Invoice from DynamoDB invoice.
func (inv *Invoice) ToInvoice() (*invoice.Invoice, error) {
	date, err := time.Parse(yyyymmddFormat, inv.Date)
	if err != nil {
		return nil, err
	}

	return &invoice.Invoice{
		ID:           inv.ID,
		Number:       inv.Number,
		CustomerName: inv.CustomerName,
		Status:       invoice.Status(inv.Status),
		Date:         date,
		Items:        nil,
		CreatedAt:    inv.CreatedAt,
		UpdatedAt:    inv.UpdatedAt,
	}, nil
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
func NewItem(inv invoice.Invoice, item invoice.Item) *Item {
	pk := itemPk(inv, item)
	sk := itemSk(inv, item)

	return &Item{
		PK:        pk,
		SK:        sk,
		ID:        item.ID,
		SKU:       item.SKU,
		Name:      item.Name,
		Price:     item.Price,
		Qty:       item.Qty,
		Status:    string(item.Status),
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}

// ToItem creates an instance of invoice.Item from DynamoDB item.
func (item *Item) ToItem() *invoice.Item {
	return &invoice.Item{
		ID:        item.ID,
		SKU:       item.SKU,
		Name:      item.Name,
		Price:     item.Price,
		Qty:       item.Qty,
		Status:    invoice.Status(item.Status),
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}

func invoicePk(inv invoice.Invoice) string {
	elems := []string{invoicePkPrefix, inv.ID}
	return strings.Join(elems, keySeparator)
}

func invoiceSk(inv invoice.Invoice) string {
	elems := []string{invoiceSkPrefix, inv.ID}
	return strings.Join(elems, keySeparator)
}

func itemPk(inv invoice.Invoice, item invoice.Item) string {
	elems := []string{itemPkPrefix, inv.ID}
	return strings.Join(elems, keySeparator)
}

func itemSk(inv invoice.Invoice, item invoice.Item) string {
	elems := []string{itemSkPrefix, item.ID}
	return strings.Join(elems, keySeparator)
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
