# Naming Convention Guide

## 1 Principles

* Aim for **clarity first**, then simplicity, concision, maintainability, and consistency ([Google GitHub][1])
* A name should be **clear** (it tells readers what it *is*) and **precise** (it also tells them what it *is not*) ([Google Testing Blog][2])

## 2 Identifier Case

* Use `MixedCaps` (exported) or `mixedCaps` (unexported) for **all multi-word identifiers**.
* Never use `snake_case`, `ALL_CAPS`, or `kCamelCase` ([Google GitHub][1], [Google GitHub][3])

```go
MaxLength      // exported constant
maxLength      // unexported variable
```

## 3 Underscores: when they are allowed

1. **Package names** that exist *only* for generated code
2. **Test, Benchmark, Example** functions inside `*_test.go` files
3. Low-level libraries that mirror OS or cgo identifiers (rare) ([Google GitHub][3])

Source-file names themselves may contain underscores because filenames are *not* Go identifiers.

## 4 Package Names

* All lowercase, one word, no underscores: `auth`, `bigqueryclient` ([Google GitHub][3])
* Choose words unlikely to collide with common locals (`usercount` > `count`).
* Avoid generic names like `util`, `common`, `models`.
* If you rename an imported package (`foopb "path/to/foo_go_proto"`), use a lowercase name that matches these rules, and be **consistent across files**.

## 5 File-Scoped Test Packages

Use the `_test` suffix to isolate black-box tests, e.g. `linkedlist_test` ([Google GitHub][3])

## 6 Receiver Names

* One or two letters, an **abbreviation of the type**:

```go
type Scanner struct{}
func (s *Scanner) Scan() error { ... }   // good
```

Be consistent across every method on the same type ([Google GitHub][3])

## 7 Constant Names

* MixedCaps like everything else.
* Name for the **meaning**, not the value, and avoid the `K` prefix. ([Google GitHub][3])

```go
const MaxPacketSize = 512           // good
const MAX_PACKET_SIZE = 512         // bad
const kMaxPacketSize = 512          // bad
```

## 8 Initialisms & Acronyms

Keep each initialism exactly as it appears in prose; examples:

| Scope      | Correct  | Incorrect |
|------------|----------|-----------|
| exported   | `URLAPI` | `UrlApi`  |
| unexported | `xmlAPI` | `xmlApi`  |
| exported   | `GRPC`   | `Grpc`    |

* If the initialism starts with a lowercase letter in prose (`gRPC`, `iOS`), use the same case unless you must change the first letter to export the name ([Google GitHub][3])

## 9 Function and Method Names

* **Avoid a `Get` prefix** unless “get” is inherent (HTTP GET, etc.).

    * Simple accessors: `Size()` not `GetSize()`.
    * Expensive operations: choose verbs that hint at work, e.g. `Fetch`, `Compute`. ([Google GitHub][3])
* Boolean results: start with `Is`, `Has`, `Contains`.

## 10 Variable Names

### 10.1 Length vs Scope

| Scope (≈lines) | Suggested length | Example           |
|----------------|------------------|-------------------|
| 1-7            | 1–3 chars        | `i`, `p`, `db`    |
| 8-15           | 1 word           | `count`           |
| 15-25          | 1-2 words        | `userCount`       |
| >25            | descriptive      | `activeUserCount` |
Guideline, not a rule ([Google GitHub][3])


### 10.2 Single-Letter Identifiers

* Always okay for **receiver** (`s`, `w`) and for familiar patterns (`i` loop index, `x`,`y` coordinates).
* Use only when the expanded word would be obvious and repetitive ([Google GitHub][3])

### 10.3 Drop Needless Words

Follow the blog’s checklist ([Google Testing Blog][2])

* Omit the type: `user` not `userString`.
* Omit irrelevant detail: `boss` not `finalBattleMostDangerousBossMonster`.
* Omit context already given: inside `AnnualHolidaySale`, prefer `rebate` over `annualSaleRebate`.
* Cut vague words that add no signal: `data`, `object`, `manager`, etc.

## 11 Avoid Repetition

* The **package name** already qualifies exported symbols (`db.Load`, not `db.LoadFromDatabase`) ([Google GitHub][3])
* The **type** already tells readers what a variable is (`users`, not `usersSlice`).
* Names inside deeply nested packages shouldn’t restate the full path (`report.Name()`, not `ProjectName()`) ([Google GitHub][3])

## 12 Named Result Parameters

* Use names only when they clarify multiple results or caller actions (`ctx, cancel`).
* Don’t add names just to allow naked returns or to repeat the function name ([Google GitHub][3])

## 13 Quick Checklist

* Identifiers use MixedCaps / mixedCaps, never underscores or ALL\_CAPS.
* Package name is lowercase, single word, no underscores or “util”.
* Receiver variable is 1–2 letters (type abbreviation).
* Constant names explain the *role*, not the value.
* Functions avoid a `Get` prefix unless unavoidable.
* Variable names are proportional to scope and omit redundant words, types, and context.
* Tests and examples may use underscores; code identifiers do not.
* Acronyms keep a consistent case (`ID`, `URL`, `gRPC`).
* All code is `gofmt`-formatted (enforced separately) ([Google GitHub][1])

Keep this guide close when writing or reviewing Go code to ensure every identifier carries its weight—no more, no less.

[1]: https://google.github.io/styleguide/go/guide.html "styleguide | Style guides for Google-originated open-source projects"
[2]: https://testing.googleblog.com/2017/10/code-health-identifiernamingpostforworl.html "Google Testing Blog: Code Health: IdentifierNamingPostForWorldWideWebBlog"
[3]: https://google.github.io/styleguide/go/decisions.html "styleguide | Style guides for Google-originated open-source projects"
