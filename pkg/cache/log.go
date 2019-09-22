package cache

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-redis/redis"
	"github.com/rickbassham/example-go/pkg/logging"
	"go.uber.org/zap"
)

// LoggerHook is used to log all calls to redis.
type LoggerHook struct {
}

// BeforeProcess is called before the call to redis for a single command.
func (h LoggerHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	l := logging.FromContext(ctx)

	var full []string

	full = append(full, cmd.Name())

	for _, arg := range cmd.Args()[1:] {
		full = append(full, fmt.Sprintf("%#v", arg))
	}

	l.Info("starting redis call", zap.String("cmd", strings.Join(full, " ")))

	return ctx, nil
}

// AfterProcess is called after the call to redis for a single command.
func (h LoggerHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	l := logging.FromContext(ctx)

	var full []string

	full = append(full, cmd.Name())

	for _, arg := range cmd.Args()[1:] {
		full = append(full, fmt.Sprintf("%#v", arg))
	}

	l.Info("redis call complete", zap.String("cmd", strings.Join(full, " ")))

	return nil
}

// BeforeProcessPipeline is called before the call to redis for a group of pipelined commands.
func (h LoggerHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	l := logging.FromContext(ctx)

	var p []string

	for _, cmd := range cmds {
		var full []string

		full = append(full, cmd.Name())

		for _, arg := range cmd.Args()[1:] {
			full = append(full, fmt.Sprintf("%#v", arg))
		}

		p = append(p, strings.Join(full, " "))
	}

	l.Info("starting redis call", zap.String("cmds", strings.Join(p, "\n")))

	return ctx, nil
}

// AfterProcessPipeline is called after the call to redis for a group of pipelined commands.
func (h LoggerHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	l := logging.FromContext(ctx)

	var p []string

	for _, cmd := range cmds {
		var full []string

		full = append(full, cmd.Name())

		for _, arg := range cmd.Args()[1:] {
			full = append(full, fmt.Sprintf("%#v", arg))
		}

		p = append(p, strings.Join(full, " "))
	}

	l.Info("redis call complete", zap.String("cmds", strings.Join(p, "\n")))

	return nil
}
