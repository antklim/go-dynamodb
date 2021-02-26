package dynamo

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/antklim/go-dynamodb/invoice"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const (
	keySeparator    = "#"
	invoicePkPrefix = "INVOICE"
	invoiceSkPrefix = "INVOICE"
	itemPkPrefix    = "INVOICE" // invoice items are in the same partition as the invoice
	itemSkPrefix    = "ITEM"
	yyyymmddFormat  = "20060102"
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
func NewInvoice(inv invoice.Invoice) Invoice {
	pk := invoicePk(inv.ID)
	sk := invoiceSk(inv.ID)

	return Invoice{
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
	Status    string    `dynamodbav:"status"`
	CreatedAt time.Time `dynamodbav:"createdAt"`
	UpdatedAt time.Time `dynamodbav:"updatedAt"`
}

// NewItem creates an instance of DynamoDB item from invoice.Item.
func NewItem(inv invoice.Invoice, item invoice.Item) Item {
	pk := itemPk(inv.ID)
	sk := itemSk(item.ID)

	return Item{
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
func (item *Item) ToItem() invoice.Item {
	return invoice.Item{
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

func invoicePk(invoiceID string) string {
	elems := []string{invoicePkPrefix, invoiceID}
	return strings.Join(elems, keySeparator)
}

func invoiceSk(invoiceID string) string {
	elems := []string{invoiceSkPrefix, invoiceID}
	return strings.Join(elems, keySeparator)
}

func itemPk(invoiceID string) string {
	elems := []string{itemPkPrefix, invoiceID}
	return strings.Join(elems, keySeparator)
}

func itemSk(itemID string) string {
	elems := []string{itemSkPrefix, itemID}
	return strings.Join(elems, keySeparator)
}

type repository struct {
	client dynamodbiface.DynamoDBAPI
	table  *string
}

// NewRepository ...
func NewRepository(client dynamodbiface.DynamoDBAPI, table string) invoice.Repository {
	return &repository{client: client, table: aws.String(table)}
}

func (r *repository) AddInvoice(ctx context.Context, inv invoice.Invoice) error {
	dbinv := NewInvoice(inv)
	putInvoiceItem, err := dynamodbattribute.MarshalMap(dbinv)
	if err != nil {
		return err
	}

	putItems := []*dynamodb.TransactWriteItem{}
	putItems = append(putItems, &dynamodb.TransactWriteItem{Put: &dynamodb.Put{
		TableName: r.table,
		Item:      putInvoiceItem,
	}})

	for _, item := range inv.Items {
		dbitem := NewItem(inv, item)
		putInvoiceItemItem, err := dynamodbattribute.MarshalMap(dbitem)
		if err != nil {
			return err
		}

		putItems = append(putItems, &dynamodb.TransactWriteItem{Put: &dynamodb.Put{
			TableName: r.table,
			Item:      putInvoiceItemItem,
		}})
	}

	transaction := &dynamodb.TransactWriteItemsInput{TransactItems: putItems}
	if err := transaction.Validate(); err != nil {
		return err
	}

	_, err = r.client.TransactWriteItemsWithContext(ctx, transaction)

	return err
}

func (r *repository) GetInvoice(ctx context.Context, invoiceID string) (*invoice.Invoice, error) {
	return nil, errors.New("not implemented")
}

func (r *repository) CancelInvoice(ctx context.Context, invoiceID string) error {
	return errors.New("not implemented")
}

func (r *repository) AddItem(ctx context.Context, invoiceID string, item invoice.Item) error {
	return errors.New("not implemented")
}

func (r *repository) GetItem(ctx context.Context, itemID string) (*invoice.Item, error) {
	return nil, errors.New("not implemented")
}

func (r *repository) GetItemsByStatus(ctx context.Context, status invoice.Status) ([]invoice.Item, error) {
	filt := expression.And(
		expression.Name("sk").BeginsWith(itemSkPrefix+keySeparator),
		expression.Name("status").Equal(expression.Value(status)),
	)
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.ScanInput{
		TableName:                 r.table,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	}

	result, err := r.client.ScanWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	return toInvoiceItems(result.Items)
}

func (r *repository) GetInvoiceItemsByStatus(
	ctx context.Context, invoiceID string, status invoice.Status) ([]invoice.Item, error) {

	pk := itemPk(invoiceID)
	keyCond := expression.KeyAnd(
		expression.Key("pk").Equal(expression.Value(pk)),
		expression.Key("sk").BeginsWith(itemSkPrefix+keySeparator),
	)

	filt := expression.Name("status").Equal(expression.Value(status))

	expr, err := expression.NewBuilder().
		WithKeyCondition(keyCond).
		WithFilter(filt).
		Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.QueryInput{
		TableName:                 r.table,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
	}

	result, err := r.client.QueryWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	return toInvoiceItems(result.Items)
}

func (r *repository) UpdateInvoiceItemsStatus(
	ctx context.Context, invoiceID string, status invoice.Status) error {

	items, err := r.GetInvoiceItemsByStatus(ctx, invoiceID, "CANCELLED")
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}

	upd := expression.Set(expression.Name("status"), expression.Value(status))
	expr, err := expression.NewBuilder().WithUpdate(upd).Build()
	if err != nil {
		return err
	}

	updates := invoiceItemsToUpdates(items, invoiceID, r.table, expr)
	updateItems := make([]*dynamodb.TransactWriteItem, len(updates))
	for idx, update := range updates {
		updateItems[idx] = &dynamodb.TransactWriteItem{Update: update}
	}

	transaction := &dynamodb.TransactWriteItemsInput{TransactItems: updateItems}
	if err := transaction.Validate(); err != nil {
		return err
	}

	_, err = r.client.TransactWriteItemsWithContext(ctx, transaction)

	return err
}

func (r *repository) ReplaceItems(ctx context.Context, invoiceID string, items []invoice.Item) error {
	return errors.New("not implemented")
}

func toInvoiceItems(rawItems []map[string]*dynamodb.AttributeValue) ([]invoice.Item, error) {
	if len(rawItems) == 0 {
		return nil, nil
	}

	dbItems := []*Item{}
	if err := dynamodbattribute.UnmarshalListOfMaps(rawItems, &dbItems); err != nil {
		return nil, err
	}

	invoiceItems := make([]invoice.Item, len(dbItems))
	for idx, dbItem := range dbItems {
		invoiceItems[idx] = dbItem.ToItem()
	}

	return invoiceItems, nil
}

func invoiceItemsToUpdates(
	items []invoice.Item,
	invoiceID string,
	table *string,
	expr expression.Expression,
) []*dynamodb.Update {

	pk := aws.String(itemPk(invoiceID))

	updates := make([]*dynamodb.Update, len(items))

	for idx, item := range items {
		sk := aws.String(itemSk(item.ID))
		updates[idx] = &dynamodb.Update{
			TableName: table,
			Key: map[string]*dynamodb.AttributeValue{
				"pk": {
					S: pk,
				},
				"sk": {
					S: sk,
				},
			},
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			ConditionExpression:       expr.Condition(),
			UpdateExpression:          expr.Update(),
		}
	}

	return updates
}
