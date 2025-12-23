# Go Unit Testing Guide

This guide defines consistent testing conventions for Go projects.  
All examples are self-contained and focus on clarity, determinism, and reproducibility.

## 1. Core Principles

- Tests must be **clear, isolated, and deterministic**.
- Each test must be runnable **independently** and in any order.
- Unit tests should be **fast**; mark integration or slow tests with `testing.Short()`.

## 2. Test File Structure

- Place tests in files named `*_test.go`, next to the source code.
- Run all tests with `go test ./...`.
- Continuous integration should always include:
```sh
  go test -race -cover ./...
```

## 3. Naming Conventions

**Test Functions**

* Follow the `TestXxx` pattern.
  Example:

```go
  func TestParser_RejectsEmptyInput(t *testing.T) {}
  func TestStore_Save_RetryOnConflict(t *testing.T) {}
```

**Subtests**

* Use descriptive names:

```go
  t.Run("valid input", ...)
  t.Run("empty string returns error", ...)
```

**Table-Driven Tests**

* The map key is the subtest name.
  Use meaningful, stable names (e.g., `"valid/email"`, `"error/timeout"`).

## 4. Table-Driven Tests

* Define cases as `map[string]struct{...}`.
* Each entry should include inputs and expected outputs.
* Use `t.Run` for each case.
* Apply `t.Parallel()` only when test cases are independent.

**Example:**

```go
cases := map[string]struct {
    in   string
    want string
}{
    "valid/uppercase": {in: "a", want: "A"},
    "error/empty":     {in: "", want: ""},
}

for name, tc := range cases {
    tc := tc
    t.Run(name, func(t *testing.T) {
        t.Parallel()
        got := ToUpper(tc.in)
        require.Equal(t, tc.want, got)
    })
}
```

**Rationale**
Table-driven tests encourage consistency, reduce duplication, and simplify test maintenance.

## 5. Assertions

* Use the **Testify** library for assertions:

    * `require` for critical conditions (stop test on failure)
    * `assert` for multiple validations within one test

**Example:**

```go
require.NoError(t, err)
assert.Equal(t, expected, actual)
```

## 6. Mocking

* Use **testify/mock** for mocking interfaces.
* Define expected method calls and return values.
* Always verify mock expectations at the end of the test.

**Example:**

```go
m.On("Save", mock.Anything, item).Return(nil).Once()
require.NoError(t, svc.SaveItem(ctx, item))
m.AssertExpectations(t)
```

**Rationale**
Mocks help isolate dependencies and ensure the test covers only the target logic.

## 7. HTTP Testing

* Use `httptest.NewRequest` and `httptest.NewRecorder` for handler testing.
* Use `httptest.NewServer` for external API stubs or integration-style tests.

**Example:**

```go
req := httptest.NewRequest(http.MethodGet, "/ping", nil)
w := httptest.NewRecorder()
handler.ServeHTTP(w, req)
require.Equal(t, http.StatusOK, w.Code)
```

## 8. Fuzzing and Examples

* Add **fuzz tests** (`FuzzXxx`) for parsing, serialization, or validation logic.
* Provide **example tests** (`ExampleXxx`) with `// Output:` comments to verify documentation examples.

**Example:**

```go
func ExampleAdd() {
    fmt.Println(Add(1, 2))
    // Output: 3
}
```

## 9. Coverage and Continuous Integration

* Measure coverage:

```sh
  go test -coverprofile=coverage.out ./...
  go tool cover -html=coverage.out
```
* Use the `-short` flag to skip slow tests:

```sh
  go test -short ./...
```

**Rationale**
Track coverage for insight, but optimize for quality and relevance over raw percentage.

## 10. Quick Reference

* Tests live in `*_test.go`.
* Each test is independent and order-agnostic.
* Use table-driven patterns (`map[string]struct`) with descriptive keys.
* Use `Testify` for assertions and mocks (`assert`, `require`, `mock`).
* Use `httptest` for HTTP tests.
* Add fuzz and example tests for critical components.
* Run CI with `-race`, `-cover`, and `-short`.

## 11. Summary

* Keep tests deterministic, simple, and readable.
* Mock external dependencies.
* Prefer table-driven and subtest structures.
* Write examples and fuzz tests to improve reliability and documentation.
