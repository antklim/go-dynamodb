package memory_test

import (
	"context"
	"testing"

	"github.com/antklim/go-dynamodb/invoice"
	"github.com/antklim/go-dynamodb/memory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var repo = memory.NewRepository()

func TestInvoiceGet(t *testing.T) {
	inv1, err := repo.GetInvoice(context.Background(), "")
	require.NoError(t, err)
	assert.Nil(t, inv1)

	inv2 := invoice.Invoice{
		ID: uuid.NewString(),
	}
	err = repo.AddInvoice(context.Background(), inv2)
	require.NoError(t, err)

	inv3, err := repo.GetInvoice(context.Background(), inv2.ID)
	require.NoError(t, err)
	assert.Equal(t, inv2, *inv3)
}

func TestItemGet(t *testing.T) {
	item1, err := repo.GetItem(context.Background(), "", "")
	require.NoError(t, err)
	assert.Nil(t, item1)

	item2 := invoice.Item{
		ID: uuid.NewString(),
	}
	err = repo.AddItem(context.Background(), item2)
	require.NoError(t, err)

	item3, err := repo.GetItem(context.Background(), "", item2.ID)
	require.NoError(t, err)
	assert.Equal(t, item2, *item3)

	err = repo.DeleteItem(context.Background(), "", item2.ID)
	require.NoError(t, err)

	item4, err := repo.GetItem(context.Background(), "", item2.ID)
	require.NoError(t, err)
	assert.Nil(t, item4)
}

func TestItemsScanners(t *testing.T) {
	itms := []invoice.Item{
		{
			ID:        uuid.NewString(),
			InvoiceID: "1",
			Status:    invoice.New,
		},
		{
			ID:        uuid.NewString(),
			InvoiceID: "1",
			Status:    invoice.Cancelled,
		},
		{
			ID:        uuid.NewString(),
			InvoiceID: "2",
			Status:    invoice.New,
		},
	}

	for _, item := range itms {
		err := repo.AddItem(context.Background(), item)
		require.NoError(t, err)
	}

	t.Run("returns items by status", func(t *testing.T) {
		items, err := repo.GetItemsByStatus(context.Background(), invoice.New)
		require.NoError(t, err)
		assert.Len(t, items, 2)

		items, err = repo.GetItemsByStatus(context.Background(), invoice.Pending)
		require.NoError(t, err)
		assert.Empty(t, items)
	})

	t.Run("returns invoice items", func(t *testing.T) {
		items, err := repo.GetInvoiceItems(context.Background(), "1")
		require.NoError(t, err)
		assert.Len(t, items, 2)

		items, err = repo.GetInvoiceItems(context.Background(), "2")
		require.NoError(t, err)
		assert.Len(t, items, 1)

		items, err = repo.GetInvoiceItems(context.Background(), "3")
		require.NoError(t, err)
		assert.Empty(t, items)
	})

	t.Run("returns invoice items by status", func(t *testing.T) {
		items, err := repo.GetInvoiceItemsByStatus(context.Background(), "1", invoice.New)
		require.NoError(t, err)
		assert.Len(t, items, 1)

		items, err = repo.GetInvoiceItemsByStatus(context.Background(), "2", invoice.Pending)
		require.NoError(t, err)
		assert.Empty(t, items)
	})
}
