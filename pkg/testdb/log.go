package testdb

import (
	"context"

	"go.uber.org/zap"

	"github.com/rickbassham/example-go/pkg/logging"
)

type Logger struct {
}

func (l Logger) Before(ctx context.Context, name, statement string, args ...interface{}) (context.Context, error) {
	log := logging.FromContext(ctx)
	log.Info("executing sql statement", zap.String("sql_statement_name", name))

	return ctx, nil
}

func (l Logger) After(ctx context.Context, err error, name, statement string, args ...interface{}) error {
	log := logging.FromContext(ctx)
	log.Info("sql statement complete", zap.String("sql_statement_name", name))

	return nil
}
