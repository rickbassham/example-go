# cache
--
    import "."


## Usage

#### type Cache

```go
type Cache struct {
}
```

Cache is a logged and instrumented wrapper around a redis client.

#### func  New

```go
func New(client Client) *Cache
```
New creates a new Cache.

#### func (*Cache) GetValue

```go
func (c *Cache) GetValue(ctx context.Context) (string, error)
```
GetValue retrieves the value from redis.

#### type Client

```go
type Client interface {
	AddHook(redis.Hook)
	WithContext(context.Context) *redis.Client
}
```

Client represents the functions needed for this wrapper.

#### type LoggerHook

```go
type LoggerHook struct {
}
```

LoggerHook is used to log all calls to redis.

#### func (LoggerHook) AfterProcess

```go
func (h LoggerHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error
```
AfterProcess is called after the call to redis for a single command.

#### func (LoggerHook) AfterProcessPipeline

```go
func (h LoggerHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error
```
AfterProcessPipeline is called after the call to redis for a group of pipelined
commands.

#### func (LoggerHook) BeforeProcess

```go
func (h LoggerHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error)
```
BeforeProcess is called before the call to redis for a single command.

#### func (LoggerHook) BeforeProcessPipeline

```go
func (h LoggerHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error)
```
BeforeProcessPipeline is called before the call to redis for a group of
pipelined commands.

#### type NewRelicHook

```go
type NewRelicHook struct {
}
```

NewRelicHook is used to instrument all calls to redis using a newrelic segment.

#### func (NewRelicHook) AfterProcess

```go
func (h NewRelicHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error
```
AfterProcess is called after the call to redis for a single command.

#### func (NewRelicHook) AfterProcessPipeline

```go
func (h NewRelicHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error
```
AfterProcessPipeline is called after the call to redis for a group of pipelined
commands.

#### func (NewRelicHook) BeforeProcess

```go
func (h NewRelicHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error)
```
BeforeProcess is called before the call to redis for a single command.

#### func (NewRelicHook) BeforeProcessPipeline

```go
func (h NewRelicHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error)
```
BeforeProcessPipeline is called before the call to redis for a group of
pipelined commands.
