package env

import (
	"time"

	"github.com/caarlos0/env"
)

// Config represents the common environment variables needed for all apps.
type Config struct {
	AppName         string    `env:"APP_NAME,required"`
	TeamName        string    `env:"TEAM_NAME,required"`
	Environment     string    `env:"APP_ENV,required"`
	BuildDate       time.Time `env:"BUILD_DATE"`
	BuildGitHash    string    `env:"BUILD_GIT_HASH,required"`
	BuildGitTag     string    `env:"BUILD_GIT_TAG,required"`
	NewRelicLicense string    `env:"NEW_RELIC_LICENSE,required"`
}

// Load will bind the environment variables to the given config.
func Load(c interface{}) error {
	return env.Parse(c)
}
