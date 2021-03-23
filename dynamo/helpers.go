package dynamo

import (
	"strings"

	"github.com/antklim/go-dynamodb/invoice"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

func invoicePartitionKey(invoiceID string) string {
	elems := []string{invoicePkPrefix, invoiceID}
	return strings.Join(elems, keySeparator)
}

func invoiceSortKey(invoiceID string) string {
	elems := []string{invoiceSkPrefix, invoiceID}
	return strings.Join(elems, keySeparator)
}

func invoicePrimaryKey(invoiceID string) (map[string]*dynamodb.AttributeValue, error) {
	primaryKey := map[string]string{
		"pk": invoicePartitionKey(invoiceID),
		"sk": invoiceSortKey(invoiceID),
	}

	return dynamodbattribute.MarshalMap(primaryKey)
}

func itemPk(invoiceID string) string {
	elems := []string{itemPkPrefix, invoiceID}
	return strings.Join(elems, keySeparator)
}

func itemSk(itemID string) string {
	elems := []string{itemSkPrefix, itemID}
	return strings.Join(elems, keySeparator)
}

func toInvoice(rawItem map[string]*dynamodb.AttributeValue) (*invoice.Invoice, error) {
	if rawItem == nil {
		return nil, nil
	}

	dbInvoice := Invoice{}
	if err := dynamodbattribute.UnmarshalMap(rawItem, &dbInvoice); err != nil {
		return nil, err
	}

	return dbInvoice.ToInvoice()
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

func invoiceItemsToUpdates(items []invoice.Item, table *string, expr expression.Expression) []*dynamodb.Update {
	updates := make([]*dynamodb.Update, len(items))

	for idx, item := range items {
		pk := aws.String(itemPk(item.InvoiceID))
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

func invoiceItemsToPuts(items []invoice.Item, table *string) ([]*dynamodb.Put, error) {
	putItems := make([]*dynamodb.Put, len(items))

	for idx, item := range items {
		dbitem := NewItem(item)
		putItem, err := dynamodbattribute.MarshalMap(dbitem)
		if err != nil {
			return nil, err
		}

		putItems[idx] = &dynamodb.Put{
			TableName: table,
			Item:      putItem,
		}
	}

	return putItems, nil
}
