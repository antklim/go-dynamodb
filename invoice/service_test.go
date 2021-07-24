package invoice_test

import (
	"os"
	"testing"

	"github.com/antklim/go-dynamodb/dynamo"
	"github.com/antklim/go-dynamodb/invoice"
	"github.com/antklim/go-dynamodb/memory"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Test against in memory DB
// Test against mock - use mock to test dynamodb integration
// Test against real instance

var useInMemoryDB = false

func initRepo() invoice.Repository {
	dburl := os.Getenv("TEST_DB_URL")
	dbtable := os.Getenv("TEST_DB_TABLE")
	if dburl == "" || dbtable == "" {
		useInMemoryDB = true
		return memory.NewRepository()
	}

	cfg := &aws.Config{}
	cfg.WithEndpoint(dburl).WithRegion("ap-southeast-2")
	sess := session.Must(session.NewSession(cfg))
	dbapi := dynamodb.New(sess)
	return dynamo.NewRepository(dbapi, dbtable)
}

func TestService(t *testing.T) {
	t.Run("given an inovice does not exist", func(t *testing.T) {
		t.Run("when call GetInvoice then expect nothing to be returned", func(t *testing.T) {})
		t.Run("when call CancelInvoice then expect nothing to be returned", func(t *testing.T) {})
		t.Run("when call AddItem then expect error to be returned", func(t *testing.T) {})
		t.Run("when GetItem then expect nothing to be returned", func(t *testing.T) {})
	})

	t.Run("given an existing active invoice", func(t *testing.T) {
		t.Run("when call GetInvoice then expect invoice to be returned", func(t *testing.T) {})
		t.Run("when call AddItem then item to be added to invoice", func(t *testing.T) {})
		// TODO: GetItem will be called in the previous test
		// t.Run("when GetItem then expect nothing returned", func(t *testing.T) {})
		t.Run("when call CancelInvoice then expect invoice to be cancelled", func(t *testing.T) {})
	})

	t.Run("given a cancelled invoice", func(t *testing.T) {
		t.Run("when call GetInvoice then expect invoice to be returned", func(t *testing.T) {})
		t.Run("when call AddItem then expect error to be returned", func(t *testing.T) {})
		t.Run("when GetItem then expect invoice item to be returned", func(t *testing.T) {})
		t.Run("when call CancelInvoice then expect error to be returned", func(t *testing.T) {})
	})
}

func TestServicePropagatesStorageErrors(t *testing.T) {
	if useInMemoryDB {
		t.Skip("skipping for in memory DB")
	}

	t.Log("implements test")
}
