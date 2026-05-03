# Contributing to SSHH

First off, thanks for taking the time to contribute! 🎉

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
- [Development Setup](#development-setup)
- [Coding Standards](#coding-standards)
- [Commit Messages](#commit-messages)
- [Pull Request Process](#pull-request-process)

## Code of Conduct

This project and everyone participating in it is governed by our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## How Can I Contribute?

### 🐛 Reporting Bugs

Before creating a bug report, please check existing issues to avoid duplicates.

When filing a bug, include:
- **Go version** (`go version`)
- **OS and version**
- **Steps to reproduce**
- **Expected vs actual behavior**
- **Relevant logs or error messages**

### 💡 Suggesting Features

Feature requests are welcome! Please open an issue with:
- **Use case** — What problem does this solve?
- **Proposed solution** — How should it work?
- **Alternatives considered** — What other approaches did you think of?

### 🔧 Submitting Changes

1. Fork the repo
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run checks (see below)
5. Commit with a clear message
6. Push to your fork
7. Open a Pull Request

## Development Setup

### Prerequisites

- **Go 1.21+**
- **golangci-lint** (optional, for linting)
- **Git**

### Getting Started

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/sshh.git
cd sshh

# Install dependencies
go mod download

# Build
go build -o sshh.exe .

# Run tests
go test -v ./...

# Set up git hooks
.\scripts\setup-hooks.ps1
```

### Using the Makefile

```bash
make build      # Build the binary
make test       # Run tests
make lint       # Run linters
make fmt        # Format code
make vet        # Run go vet
make all        # Format + lint + test + build
```

## Coding Standards

### Go Style

- Follow standard [Go conventions](https://go.dev/doc/effective_go)
- Run `gofmt` and `goimports` before committing
- All exported functions must have doc comments
- Keep functions small and focused
- Handle all errors explicitly — no `_ = err`

### Project Structure

```
internal/       # Private packages (not importable by others)
  crypto/       # Encryption/decryption (AES-256-GCM + Argon2)
  vault/        # Credential vault operations
  sshclient/    # SSH connection logic
  installer/    # Windows PATH installer
  ui/           # Terminal UI helpers
cmd/            # CLI command definitions (cobra)
```

### Security Rules

- **NEVER** log or print passwords, keys, or master passwords
- **NEVER** store secrets in plaintext
- **ALWAYS** zero-out sensitive byte slices after use
- **ALWAYS** use constant-time comparison for secrets
- Test crypto changes thoroughly

## Commit Messages

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

### Types

| Type | Description |
|------|-------------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation only |
| `style` | Formatting, no code change |
| `refactor` | Code change that neither fixes a bug nor adds a feature |
| `perf` | Performance improvement |
| `test` | Adding or updating tests |
| `chore` | Build process or auxiliary tool changes |
| `ci` | CI/CD changes |

### Examples

```
feat(vault): add support for SSH key passphrases
fix(crypto): handle empty master password gracefully
docs(readme): add installation instructions for Linux
chore(ci): add GitHub Actions workflow
```

## Pull Request Process

1. **Update docs** — If your change affects usage, update the README
2. **Add tests** — New features need tests, bug fixes need regression tests
3. **Pass CI** — All checks must pass (lint, test, build)
4. **One feature per PR** — Keep PRs focused and small
5. **Describe your changes** — Use the PR template and explain *why*

### PR Title Format

Same as commit messages:
```
feat(vault): add credential tags support
```

### Review Process

- At least 1 approval required
- Maintainers may request changes
- Once approved, a maintainer will merge

---

## Questions?

Feel free to open an issue or start a discussion. We're happy to help! 🚀
