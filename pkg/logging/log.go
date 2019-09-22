package logging

import (
	"os"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/rickbassham/example-go/pkg/env"
)

// Initialize creates a new JSON zap logger.
func Initialize(c env.Config) *zap.Logger {
	logEnc := zap.NewProductionEncoderConfig()
	logEnc.EncodeTime = zapcore.ISO8601TimeEncoder
	logEnc.TimeKey = "timestamp"

	log := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(logEnc),
		zapcore.Lock(os.Stdout),
		zap.NewAtomicLevel(),
	))

	log = log.With(
		zap.String("app_name", c.AppName),
		zap.String("environment", c.Environment),
		zap.String("build_git_hash", c.BuildGitHash),
		zap.String("build_git_tag", c.BuildGitTag),
		zap.Time("build_date", c.BuildDate),
		zap.String("team", c.TeamName),
		zap.String("run_id", uuid.New().String()),
		zap.Time("start_time", time.Now()),
	)

	return log
}
