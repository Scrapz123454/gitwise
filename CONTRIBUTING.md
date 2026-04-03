# Contributing to gitwise

Thanks for your interest in contributing! This document covers everything you need to get started.

## Development Setup

### Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- [Git](https://git-scm.com/)
- An LLM provider for testing:
  - [Ollama](https://ollama.ai/) (recommended — free, local)
  - Or an API key for OpenAI/Anthropic/Gemini

### Clone and Build

```bash
git clone https://github.com/aymenhmaidiwastaken/gitwise.git
cd gitwise
go mod download
make build
./gitwise version
```

### Run Tests

```bash
make test
```

### Lint

```bash
# Install golangci-lint: https://golangci-lint.run/welcome/install-locally/
make lint
```

## Project Structure

```
gitwise/
├── main.go                      # Entry point
├── cmd/                         # CLI commands (Cobra)
│   ├── root.go                  #   Root command + global flags
│   ├── commit.go                #   gitwise commit
│   ├── pr.go                    #   gitwise pr
│   ├── diff.go                  #   gitwise diff
│   ├── review.go                #   gitwise review
│   ├── changelog.go             #   gitwise changelog
│   ├── lint.go                  #   gitwise lint
│   ├── amend.go                 #   gitwise amend
│   ├── hook.go                  #   gitwise hook install/uninstall
│   ├── setup.go                 #   gitwise setup (wizard)
│   ├── config.go                #   gitwise config show/init
│   ├── version.go               #   gitwise version
│   └── helpers.go               #   Shared utilities
├── internal/
│   ├── config/                  # YAML config + env vars
│   ├── git/                     # Git operations (diff, log, branch, hooks, monorepo)
│   ├── llm/                     # LLM providers (ollama, openai, anthropic, gemini, custom)
│   ├── generator/               # Prompt pipelines (commit, PR, review, changelog, lint)
│   └── ui/                      # Terminal UI (spinner, TUI picker, interactive prompts)
├── npm/                         # npm wrapper package
├── .github/workflows/           # CI + release automation
├── .goreleaser.yaml             # Cross-platform release builds
└── Makefile                     # Build commands
```

## How to Add a New LLM Provider

1. Create `internal/llm/yourprovider.go` implementing the `Provider` interface:

```go
type Provider interface {
    Generate(ctx context.Context, prompt string) (string, error)
    Name() string
}
```

2. Optionally implement `StreamingProvider` for token-by-token streaming:

```go
type StreamingProvider interface {
    Provider
    GenerateStream(ctx context.Context, prompt string, callback StreamCallback) (string, error)
}
```

3. Register your provider in `internal/llm/provider.go` → `NewProvider()` switch.
4. Add any config fields to `internal/config/config.go`.
5. Add tests in `internal/llm/yourprovider_test.go`.

## How to Add a New Command

1. Create `cmd/yourcommand.go`:

```go
var yourCmd = &cobra.Command{
    Use:   "yourcommand",
    Short: "One-line description",
    RunE:  runYourCommand,
}

func init() {
    rootCmd.AddCommand(yourCmd)
}

func runYourCommand(cmd *cobra.Command, args []string) error {
    cfg, err := loadConfig(cmd)
    // ...
}
```

2. Use `loadConfig(cmd)` from `helpers.go` to get config with CLI flag overrides.
3. Use `ui.NewSpinner()` for loading states.

## Code Style

- Follow standard Go conventions (`go fmt`, `go vet`)
- Error messages: lowercase, no trailing period, wrap with `%w`
- Exported functions need doc comments
- Tests go in `_test.go` files in the same package
- Use table-driven tests where possible

## Commit Messages

We use [Conventional Commits](https://www.conventionalcommits.org/):

```
feat(scope): add new feature
fix(scope): resolve bug
docs: update README
test: add unit tests
refactor(scope): restructure code
chore: update dependencies
```

Or just use `gitwise commit` to generate them!

## Pull Requests

1. Fork the repo
2. Create a feature branch: `git checkout -b feat/your-feature`
3. Make your changes
4. Run tests: `make test`
5. Push and open a PR

## Reporting Issues

- Use GitHub Issues
- Include: Go version, OS, provider used, steps to reproduce
- For LLM output quality issues, include the diff that produced bad results

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
