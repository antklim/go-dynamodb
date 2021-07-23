package memory

import (
	"errors"
	"io"

	"github.com/antklim/go-dynamodb/invoice"
)

var (
	errEndOfTable = errors.New("end of table")
)

type itemsReader struct {
	t map[string]invoice.Item
	i int // current reading index
}

func (r *itemsReader) Read(b []invoice.Item) (n int, err error) {
	if r.i >= len(r.t) {
		return 0, errEndOfTable
	}

	tt := make([]invoice.Item, len(r.t)-r.i)
	for _, v := range r.t {
		tt = append(tt, v)
	}
	n = copy(b, tt)
	r.i += n
	return
}

func newInvoiceReader(t map[string]invoice.Item) *itemsReader {
	return &itemsReader{t, 0}
}

// TODO: implement filtering using itemsFilter
type itemsFilter struct {
	r io.Reader
	f itemFilter
	b []invoice.Item
}

func newItemsFilter(r io.Reader, f itemFilter) *itemFilter {
	return nil
}
