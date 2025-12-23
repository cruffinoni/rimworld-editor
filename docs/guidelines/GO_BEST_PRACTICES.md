# GO_BEST_PRACTICES.md

This document defines coding and design guidelines for Go projects.  
Each section presents preferred ("Do") and discouraged ("Don’t") patterns, followed by a short rationale.

## 1. Imports and Formatting

**Do**

```sh
gofmt -s -w ./...
goimports -w ./...
````

**Don’t**

```go
import . "fmt"
```

**Rationale**
Always format code using `gofmt` and organize imports with `goimports`.
Avoid `import .` as it pollutes the namespace and reduces readability.

## 2. Package Design and Documentation

**Do**

```go
// Package cache provides in-memory caching.
package cache
```

**Don’t**

```go
package util
```

**Rationale**
Use small, focused packages with clear purpose and documentation comments.
Avoid vague package names such as `util` or `common`.

## 3. Errors

**Do**

```go
if err != nil {
    return fmt.Errorf("open config: %w", err)
}
if errors.Is(err, os.ErrNotExist) { /* ... */ }
```

**Don’t**

```go
panic(err)
return fmt.Errorf("Open Config Failed: %v", err)
```

**Rationale**
Treat errors as values. Use `%w` for wrapping.
Keep error messages lowercase and without punctuation.
Avoid `panic` in normal program flow.

## 4. Context Usage

**Do**

```go
func Fetch(ctx context.Context, url string) error {
    ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel()
    // ...
}
```

**Don’t**

```go
type Client struct{ ctx context.Context }
func Fetch(url string) error
```

**Rationale**
Pass `context.Context` as the first argument.
Use it only for cancellation, timeouts, and request-scoped values.
Do not store contexts in structs.

## 5. Concurrency

**Do**

```go
g, ctx := errgroup.WithContext(ctx)
g.Go(func() error { /* ... */ return nil })
if err := g.Wait(); err != nil { /* ... */ }
```

**Don’t**

```go
var wg sync.WaitGroup
go func() { /* ignore errors */ }()
wg.Wait()
```

**Rationale**
Prefer structured concurrency using `errgroup`.
Ensure cancellation on error and avoid fire-and-forget goroutines.

## 6. HTTP and I/O

**Do**

```go
req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
resp, err := http.DefaultClient.Do(req)
if err != nil { return err }
defer resp.Body.Close()
```

**Don’t**

```go
resp, _ := http.Get(u)
body, _ := io.ReadAll(resp.Body)
```

**Rationale**
Always attach a `Context` and close response bodies.
Reuse a single `http.Client` for efficiency.

## 7. Interfaces

**Do**

```go
type Reader interface { Read(p []byte) (int, error) }
func NewStore(db *sql.DB) *Store
```

**Don’t**

```go
type BigInterface interface { /* many methods */ }
func NewStore(i BigInterface) Interface
```

**Rationale**
Keep interfaces small and defined by the consumer.
Return concrete types instead of interfaces.

## 8. Defer, Panic, and Recover

**Do**

```go
f, err := os.Open(name)
if err != nil { return err }
defer f.Close()
```

**Don’t**

```go
defer f.Close() // before checking err
panic("unexpected error")
```

**Rationale**
Call `defer` only after a successful resource acquisition.
Use `panic/recover` only at process boundaries.

## 9. Slices, Maps, and Strings

**Do**

```go
out := make([]T, 0, n)
var s []T
var b strings.Builder
b.WriteString("x")
```

**Don’t**

```go
out := []T{}
var b strings.Builder = *otherB
```

**Rationale**
Preallocate when size is known.
Nil slices are valid.
Do not copy initialized `strings.Builder` values.

## 10. JSON and Zero Values

**Do**

```go
type X struct { Items []int `json:"items,omitempty"` }
```

**Don’t**

```go
// Expecting nil slice to encode as []
```

**Rationale**
Nil slices encode as `null`.
Use `omitempty` or initialize to empty slices when necessary.

## 11. Logging

**Do**

```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
slog.SetDefault(logger)
slog.Info("fetch", "url", u, "attempt", n)
```

**Don’t**

```go
log.Printf("fetch %s %d", u, n)
```

**Rationale**
Use structured logging with key-value pairs.
Avoid unstructured `Printf`-style logs.

## 12. Tooling and Modules

**Do**

```sh
go vet ./...
go test -race ./...
```

```go
// go.mod
module example.com/project
go 1.23
```

**Don’t**

```go
// Missing go directive
```

**Rationale**
Run static analysis (`vet`) and the race detector.
Keep `go.mod` up to date and follow semantic versioning.

## 13. Timeouts, Tickers, and Timers

**Do**

```go
ctx, cancel := context.WithTimeout(ctx, d)
defer cancel()

t := time.NewTicker(d)
defer t.Stop()
```

**Don’t**

```go
time.AfterFunc(d, fn)
ticker := time.NewTicker(d)
```

**Rationale**
Use contexts for timeouts.
Always stop tickers and timers to avoid leaks.

## 14. Pipelines and Cancellation

**Do**

```go
for v := range in {
    select {
    case out <- f(v):
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

**Don’t**

```go
for v := range in {
    out <- f(v)
}
```

**Rationale**
Every stage in a pipeline must handle cancellation.
Senders should close output channels when done.

## 15. Testing

Keep tests deterministic and isolated.
Use table-driven tests and subtests.
Avoid relying on non-deterministic behavior such as map iteration order.

## 16. Usage of `init()`

**Rule**
Avoid using `init()` except for trivial or deterministic initialization.

**Rationale**

* Hidden execution before `main()` complicates reasoning and testing.
* Explicit constructors or setup functions make dependencies clearer.

**Exception**
Use `init()` only for simple, deterministic registration logic.

## 17. Variable Declarations

**Do**

```go
var (
    name    string
    age     int
    enabled bool
)
```

**Don’t**

```go
var name string
var age int
var enabled bool
```

**Rationale**
Group related declarations for readability and consistency.

# Summary

* Prefer explicitness and small, composable functions.
* Handle errors explicitly.
* Manage resources and cancellations carefully.
* Keep interfaces and packages minimal.
* Write clear, structured, and deterministic code.