# env
--
    import "github.com/rickbassham/example-go/pkg/env"


## Usage

#### func  Load

```go
func Load(c interface{}) error
```
Load will bind the environment variables to the given config.

#### type Config

```go
type Config struct {
	AppName         string    `env:"APP_NAME,required"`
	Environment     string    `env:"APP_ENV,required"`
	BuildDate       time.Time `env:"BUILD_DATE"`
	BuildGitHash    string    `env:"BUILD_GIT_HASH"`
	BuildGitTag     string    `env:"BUILD_GIT_TAG"`
	NewRelicLicense string    `env:"NEW_RELIC_LICENSE,required"`
}
```

Config represents the common environment variables needed for all apps.
