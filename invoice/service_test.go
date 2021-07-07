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

}

func TestServicePropagatesStorageErrors(t *testing.T) {
	if useInMemoryDB {
		t.Skip("skipping for in memory DB")
	}

	t.Log("implements test")
}
