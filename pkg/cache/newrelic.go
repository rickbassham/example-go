package cache

import (
	"context"
	"strings"

	"github.com/go-redis/redis/v7"
	newrelic "github.com/newrelic/go-agent"
)

var (
	newrelicKey = contextKey("newrelicsegment")
)

// NewRelicHook is used to instrument all calls to redis using a newrelic segment.
type NewRelicHook struct {
}

// BeforeProcess is called before the call to redis for a single command.
func (h NewRelicHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	txn := newrelic.FromContext(ctx)

	s := &newrelic.DatastoreSegment{
		Product:   newrelic.DatastoreRedis,
		Operation: cmd.Name(),
		StartTime: newrelic.StartSegmentNow(txn),
	}

	ctx = context.WithValue(ctx, newrelicKey, s)

	return ctx, nil
}

// AfterProcess is called after the call to redis for a single command.
func (h NewRelicHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	s := ctx.Value(newrelicKey).(*newrelic.DatastoreSegment)

	s.End() // nolint

	return nil
}

// BeforeProcessPipeline is called before the call to redis for a group of pipelined commands.
func (h NewRelicHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	var cmd []string

	for _, c := range cmds {
		cmd = append(cmd, c.Name())
	}

	txn := newrelic.FromContext(ctx)

	s := &newrelic.DatastoreSegment{
		Product:   newrelic.DatastoreRedis,
		Operation: strings.Join(cmd, " "),
		StartTime: newrelic.StartSegmentNow(txn),
	}

	ctx = context.WithValue(ctx, newrelicKey, s)

	return ctx, nil
}

// AfterProcessPipeline is called after the call to redis for a group of pipelined commands.
func (h NewRelicHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	s := ctx.Value(newrelicKey).(*newrelic.DatastoreSegment)

	s.End() // nolint

	return nil
}
