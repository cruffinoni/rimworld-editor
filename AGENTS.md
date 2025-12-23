# Repository Guidelines

## Project Structure & Module Organization
- `cmd/`: entry points for binaries (`app`, `cycle`, `gen_go`, `gen_xml`).
- `internal/`: core implementation packages (not for external use).
- `pkg/`: reusable libraries (for example `logging` and `runtimepath`).
- `generated/`: Go files produced by the code/XML generators.
- `config/`: example configuration files like `config.example.yaml`.
- `docs/guidelines/`: naming, testing, comments, and commit conventions.

## Build, Test, and Development Commands
- `make check`: quick compile of all packages (`go build ./...`).
- `make build` / `make build-all`: build all binaries into `bin/`.
- `make build-app|build-cycle|build-gen-go|build-gen-xml`: build a single binary.
- `make test`: run unit tests (`go test ./...`).
- `make vendor`: vendor dependencies (`go mod vendor`).
- `make clean`: remove `bin/`.

## Coding Style & Naming Conventions
- Format with `gofmt -s` and organize imports with `goimports`.
- Go identifiers use `MixedCaps`/`mixedCaps`; avoid `snake_case` and `ALL_CAPS`.
- Package names are lowercase and single word; avoid `util`/`common`.
- Prefer clear, minimal names and avoid `Get` prefixes for simple accessors.
- Add doc comments for exported packages/types/functions per `docs/guidelines/COMMENTS.md`.

## Testing Guidelines
- Tests live alongside code as `*_test.go`; run with `go test ./...`.
- Use Testify (`assert`, `require`, `mock`) for assertions and mocks.
- Follow `TestXxx` naming, descriptive `t.Run` subtests, and table-driven cases.
- CI-style run: `go test -race -cover ./...`; use `-short` for slow tests.

## Commit & Pull Request Guidelines
- Use Conventional Commits: `type(scope): description` (imperative, <=50 chars).
- Body is optional; wrap at 72 chars and explain what/why, not how.
- Keep one logical change per commit.
- PRs should describe the change, mention tests run, and link issues when applicable.
