package testdb

import (
	"context"

	newrelic "github.com/newrelic/go-agent"
)

type contextKey string

func (k contextKey) String() string {
	return "context key: " + string(k)
}

var (
	newrelicKey = contextKey("newrelicsegment")
)

type NewRelic struct {
}

func (l NewRelic) Before(ctx context.Context, name, statement string, args ...interface{}) (context.Context, error) {
	txn := newrelic.FromContext(ctx)

	s := &newrelic.DatastoreSegment{
		Product:            newrelic.DatastoreMySQL,
		Operation:          name,
		ParameterizedQuery: statement,
		StartTime:          newrelic.StartSegmentNow(txn),
	}

	ctx = context.WithValue(ctx, newrelicKey, s)

	return ctx, nil
}

func (l NewRelic) After(ctx context.Context, err error, name, statement string, args ...interface{}) error {
	s := ctx.Value(newrelicKey).(*newrelic.DatastoreSegment)

	s.End() // nolint

	return nil
}
