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
