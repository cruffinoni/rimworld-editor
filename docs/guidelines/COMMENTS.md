# Go Documentation and Comment Guidelines

This document defines standards for writing clear, consistent, and maintainable documentation in Go projects.  
It focuses on effective use of Go’s native documentation system and plain-text code comments.

## 1. Core Principles

### 1.1 Use Go-Native Documentation

Use Go’s built-in documentation system (`go doc`) as the main source of API documentation.  
This ensures compatibility with standard tools and keeps documentation close to the code.

### 1.2 Separate Documentation by Purpose

- **Go doc comments** — API contracts, usage examples, and expected behaviors
- **README files** — project overview, setup, and workflows
- **Inline comments** — clarify complex logic or edge cases within functions

### 1.3 Practical Focus

Document what developers need to understand and use the code effectively.  
Not every exported identifier needs documentation—prioritize complex or externally used APIs.

## 2. Package-Level Documentation

Each package must include a top-level comment that describes its purpose.

**Do**

```go
// Package cache provides in-memory caching with eviction policies.
package cache
````

**Rationale**

* Start with “Package [name] provides…”
* Describe key responsibilities and behaviors
* Keep concise; architecture details belong in README files
* Include short examples for complex packages

## 3. Type Documentation

Document exported types that have non-obvious behavior or are important for the public API.

**Example**

```go
// Coordinator manages multiple platform adapters and routes events.
//
// All Coordinator methods are safe for concurrent use.
type Coordinator struct {
    bus    *EventBus
    mu     sync.RWMutex
    closed bool
}
```

**Guidelines**

* Explain the type’s purpose and key behaviors
* Mention concurrency or thread-safety guarantees
* Skip simple data containers or obvious types

## 4. Function Documentation

Document exported functions when their behavior or purpose is not immediately clear.

**Always document:**

* Public API functions
* Functions that return errors
* Constructors and setup functions
* Functions with side effects or non-obvious logic

**Example**

```go
// Register adds a handler to the event bus.
//
// Returns an error if the handler already exists.
func (b *Bus) Register(name string, h Handler) error
```

**Optional:**

* Simple getters or one-line accessors may omit comments if self-explanatory.

## 5. Method Documentation

Document methods that have significant side effects, implement interface contracts, or can fail.

**Example**

```go
// Publish sends an event to all subscribers of the topic.
//
// If a subscriber channel is full, the message is dropped
// to prevent blocking other subscribers.
func (b *Bus) Publish(topic string, e any) int
```

**Rationale**
Highlight side effects, error conditions, or special concurrency behaviors.

## 6. Configuration and Error Types

Document configuration structs and exported error variables that are part of public APIs.

**Example**

```go
// Config contains the application's runtime configuration.
type Config struct {
    Logging LoggingConfig `mapstructure:"logging"`
    Server  ServerConfig  `mapstructure:"server"`
}

// ErrNotFound is returned when the requested item does not exist.
var ErrNotFound = errors.New("not found")
```

**Guidelines**

* Clarify how configuration is loaded or overridden
* Explain common error scenarios and recovery expectations

## 7. Documentation Tools

### 7.1 Command-Line Tools

Use standard Go tooling for documentation:

```bash
go doc ./...
godoc -http=:6060
```

### 7.2 Validation

```bash
go vet ./...
go test -run=Example
```

These commands check for missing documentation and validate that example code compiles.

## 8. Style Guidelines

### 8.1 Language and Tone

* Use **present tense**: “Registers a handler” instead of “Will register”
* Use **active voice**: “Start begins” instead of “Operations are begun”
* Be **concise and precise**: focus on behavior, not implementation
* Avoid filler words or vague terms

### 8.2 Formatting Rules

1. Keep lines under 80 characters
2. Indent example code with tabs
3. Align bullet points and lists consistently
4. Do **not** use Markdown formatting in Go comments
   (no `**bold**`, `[links]`, or headers inside comments)

## 9. Example Documentation Patterns

### Constructor Example

```go
// Example shows a typical setup for a new service.
//
//  svc := NewService(config)
//  if err := svc.Start(); err != nil {
//      log.Fatal(err)
//  }
//  defer svc.Stop()
```

### Behavioral Example

```go
// Start begins background processing.
//
// Start is idempotent—calling it multiple times has no additional effect.
```

## 10. Documentation Priorities

### High Priority

* Package comments for all packages
* Exported types and interfaces
* Public methods that return errors
* Configuration and error types
* Critical entry points or initialization code

### Medium Priority

* Functions with complex control flow
* Methods with side effects or concurrency
* Validation or transformation helpers

### Low Priority

* Simple accessors or obvious data fields
* Internal helper functions used within a single file

## 11. Documentation Anti-Patterns

**Avoid**

```go
// GetName returns the name.
func (p *Platform) GetName() string
```

**Avoid**

```go
// **Important**: Uses RWMutex for safety.
```

**Avoid**

```go
// platforms is a map of platform names to Platform instances.
```

**Instead**

```go
// Name returns the unique identifier for the platform.
func (a *Adapter) Name() string
```