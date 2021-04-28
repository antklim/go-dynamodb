package dynamo

import (
	"context"
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
	pk := invoicePartitionKey(inv.ID)
	sk := invoiceSortKey(inv.ID)

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
	InvoiceID string    `dynamodbav:"invoiceId"`
	SKU       string    `dynamodbav:"sku"`
	Name      string    `dynamodbav:"name"`
	Price     uint      `dynamodbav:"price"`
	Qty       uint      `dynamodbav:"qty"`
	Status    string    `dynamodbav:"status"`
	CreatedAt time.Time `dynamodbav:"createdAt"`
	UpdatedAt time.Time `dynamodbav:"updatedAt"`
}

// NewItem creates an instance of DynamoDB item from invoice.Item.
func NewItem(item invoice.Item) Item {
	pk := itemPartitionKey(item.InvoiceID)
	sk := itemSortKey(item.ID)

	return Item{
		PK:        pk,
		SK:        sk,
		ID:        item.ID,
		InvoiceID: item.InvoiceID,
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
		InvoiceID: item.InvoiceID,
		SKU:       item.SKU,
		Name:      item.Name,
		Price:     item.Price,
		Qty:       item.Qty,
		Status:    invoice.Status(item.Status),
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}

// Product describes product properties of Item
type Product struct {
	SKU   string `dynamodbav:"sku"`
	Name  string `dynamodbav:"name"`
	Price uint   `dynamodbav:"price"`
}

func (p *Product) ToProduct() *invoice.Product {
	return &invoice.Product{
		SKU:   p.SKU,
		Name:  p.Name,
		Price: p.Price,
	}
}

type Repository struct {
	client dynamodbiface.DynamoDBAPI
	table  *string
}

// NewRepository ...
func NewRepository(client dynamodbiface.DynamoDBAPI, table string) *Repository {
	return &Repository{client: client, table: aws.String(table)}
}

func (r *Repository) AddInvoice(ctx context.Context, inv invoice.Invoice) error {
	dbinv := NewInvoice(inv)
	putInvoiceItem, err := dynamodbattribute.MarshalMap(dbinv)
	if err != nil {
		return err
	}

	transactItems := []*dynamodb.TransactWriteItem{}
	transactItems = append(transactItems, &dynamodb.TransactWriteItem{Put: &dynamodb.Put{
		TableName: r.table,
		Item:      putInvoiceItem,
	}})

	for _, item := range inv.Items {
		dbitem := NewItem(item)
		putInvoiceItemItem, err := dynamodbattribute.MarshalMap(dbitem)
		if err != nil {
			return err
		}

		transactItems = append(transactItems, &dynamodb.TransactWriteItem{Put: &dynamodb.Put{
			TableName: r.table,
			Item:      putInvoiceItemItem,
		}})
	}

	transaction := &dynamodb.TransactWriteItemsInput{TransactItems: transactItems}
	if err := transaction.Validate(); err != nil {
		return err
	}

	_, err = r.client.TransactWriteItemsWithContext(ctx, transaction)
	return err
}

func (r *Repository) GetInvoice(ctx context.Context, invoiceID string) (*invoice.Invoice, error) {
	pk, err := invoicePrimaryKey(invoiceID)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.GetItemInput{
		TableName: r.table,
		Key:       pk,
	}

	result, err := r.client.GetItemWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	return toInvoice(result.Item)
}

func (r *Repository) AddItem(ctx context.Context, item invoice.Item) error {
	dbitem := NewItem(item)
	putItem, err := dynamodbattribute.MarshalMap(dbitem)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: r.table,
		Item:      putItem,
	}

	_, err = r.client.PutItemWithContext(ctx, input)
	return err
}

func (r *Repository) GetItem(ctx context.Context, invoiceID, itemID string) (*invoice.Item, error) {
	pk, err := itemPrimaryKey(invoiceID, itemID)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.GetItemInput{
		TableName: r.table,
		Key:       pk,
	}

	result, err := r.client.GetItemWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	return toItem(result.Item)
}

func (r *Repository) GetItemProduct(ctx context.Context, invoiceID, itemID string) (*invoice.Product, error) {
	pk, err := itemPrimaryKey(invoiceID, itemID)
	if err != nil {
		return nil, err
	}

	proj := expression.NamesList(expression.Name("sku"), expression.Name("name"), expression.Name("price"))
	expr, err := expression.NewBuilder().WithProjection(proj).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.GetItemInput{
		TableName:                r.table,
		Key:                      pk,
		ExpressionAttributeNames: expr.Names(),
		ProjectionExpression:     expr.Projection(),
	}

	result, err := r.client.GetItemWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	product := Product{}
	if err := dynamodbattribute.UnmarshalMap(result.Item, &product); err != nil {
		return nil, err
	}

	return product.ToProduct(), nil
}

func (r *Repository) DeleteItem(ctx context.Context, invoiceID, itemID string) error {
	pk, err := itemPrimaryKey(invoiceID, itemID)
	if err != nil {
		return err
	}

	input := &dynamodb.DeleteItemInput{
		TableName: r.table,
		Key:       pk,
	}

	_, err = r.client.DeleteItemWithContext(ctx, input)
	return err
}

func (r *Repository) GetItemsByStatus(ctx context.Context, status invoice.Status) ([]invoice.Item, error) {
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

func (r *Repository) GetInvoiceItems(
	ctx context.Context, invoiceID string) ([]invoice.Item, error) {

	pk := itemPartitionKey(invoiceID)
	keyCond := expression.KeyAnd(
		expression.Key("pk").Equal(expression.Value(pk)),
		expression.Key("sk").BeginsWith(itemSkPrefix+keySeparator),
	)

	expr, err := expression.NewBuilder().
		WithKeyCondition(keyCond).
		Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.QueryInput{
		TableName:                 r.table,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	result, err := r.client.QueryWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	return toInvoiceItems(result.Items)
}

func (r *Repository) GetInvoiceItemsByStatus(
	ctx context.Context, invoiceID string, status invoice.Status) ([]invoice.Item, error) {

	pk := itemPartitionKey(invoiceID)
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

func (r *Repository) UpdateInvoiceItemStatus(
	ctx context.Context, invoiceID, itemID string, status invoice.Status) error {

	pk, err := itemPrimaryKey(invoiceID, itemID)
	if err != nil {
		return err
	}

	upd := expression.
		Set(expression.Name("status"), expression.Value(status)).
		Set(expression.Name("updatedAt"), expression.Value(time.Now()))
	expr, err := expression.NewBuilder().WithUpdate(upd).Build()
	if err != nil {
		return err
	}

	input := &dynamodb.UpdateItemInput{
		TableName:                 r.table,
		Key:                       pk,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	}

	_, err = r.client.UpdateItemWithContext(ctx, input)
	return err
}

func (r *Repository) UpdateInvoiceItemsStatus(
	ctx context.Context, invoiceID string, itemIDs []string, status invoice.Status) error {

	if len(itemIDs) == 0 {
		return nil
	}

	upd := expression.
		Set(expression.Name("status"), expression.Value(status)).
		Set(expression.Name("updatedAt"), expression.Value(time.Now()))
	expr, err := expression.NewBuilder().WithUpdate(upd).Build()
	if err != nil {
		return err
	}

	transactItems := make([]*dynamodb.TransactWriteItem, len(itemIDs))
	for idx, itemID := range itemIDs {
		pk, err := itemPrimaryKey(invoiceID, itemID)
		if err != nil {
			return err
		}

		update := &dynamodb.Update{
			TableName:                 r.table,
			Key:                       pk,
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
		}
		transactItems[idx] = &dynamodb.TransactWriteItem{Update: update}
	}

	transaction := &dynamodb.TransactWriteItemsInput{TransactItems: transactItems}
	if err := transaction.Validate(); err != nil {
		return err
	}

	_, err = r.client.TransactWriteItemsWithContext(ctx, transaction)
	return err
}

// TODO: add a list of old items IDs to replace
func (r *Repository) ReplaceItems(
	ctx context.Context, invoiceID string, newItems []invoice.Item) error {

	items, err := r.GetInvoiceItemsByStatus(ctx, invoiceID, invoice.New)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}

	upd := expression.Set(expression.Name("status"), expression.Value(invoice.Cancelled))
	expr, err := expression.NewBuilder().WithUpdate(upd).Build()
	if err != nil {
		return err
	}

	updates, err := invoiceItemsToUpdates(items, r.table, expr)
	if err != nil {
		return err
	}

	puts, err := invoiceItemsToPuts(newItems, r.table)
	if err != nil {
		return err
	}

	transactionItems := []*dynamodb.TransactWriteItem{}
	for _, update := range updates {
		transactionItems = append(transactionItems, &dynamodb.TransactWriteItem{Update: update})
	}
	for _, put := range puts {
		transactionItems = append(transactionItems, &dynamodb.TransactWriteItem{Put: put})
	}

	transaction := &dynamodb.TransactWriteItemsInput{TransactItems: transactionItems}
	if err := transaction.Validate(); err != nil {
		return err
	}

	_, err = r.client.TransactWriteItemsWithContext(ctx, transaction)
	return err
}
