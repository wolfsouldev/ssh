<div align="center">

```
   _____ _____ _    _ _    _
  / ____/ ____| |  | | |  | |
 | (___| (___ | |__| | |__| |
  \___ \\___ \|  __  |  __  |
  ____) |___) | |  | | |  | |
 |_____/_____/|_|  |_|_|  |_|
```

### Secure SSH Credential Manager

_One master password. All your servers._

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat-square&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](LICENSE)
[![CI](https://img.shields.io/github/actions/workflow/status/astro/sshh/ci.yml?branch=main&style=flat-square&label=CI)](https://github.com/astro/sshh/actions)
[![Release](https://img.shields.io/github/v/release/astro/sshh?style=flat-square&color=blue)](https://github.com/astro/sshh/releases)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen?style=flat-square)](CONTRIBUTING.md)

[Features](#-features) · [Install](#-installation) · [Usage](#-usage) · [Security](#-security) · [Contributing](#-contributing)

</div>

---

## Why SSHH?

Tired of remembering dozens of SSH passwords? `sshh` replaces the standard `ssh` command with an encrypted vault that stores all your credentials locally, protected by a single master password. Just type `sshh user@host` and you're in.

```
$ sshh admin@production-server

  Connecting to admin@production-server:22 ...
  🔑 Master password: ********
  ✓ Credentials found for admin@production-server
  → Authenticating with password...

  Welcome to Ubuntu 22.04 LTS
  admin@production-server:~$
```

---

## ✨ Features

| Feature                  | Description                                              |
| ------------------------ | -------------------------------------------------------- |
| 🔐 **Master Password**   | All credentials encrypted with AES-256-GCM + Argon2id    |
| 🔑 **Passwords & Keys**  | Store SSH passwords or private keys (`.pem`, `.key`)     |
| ⚡ **Auto-connect**      | Saved credentials? Just enter master password to connect |
| 💾 **Save on first use** | New connections prompt to save after successful login    |
| 📋 **Full management**   | List, add, edit, delete credentials from the console     |
| 📦 **Export / Import**   | Backup and restore your encrypted vault                  |
| 🪟 **Windows installer** | `sshh install` adds the binary to your PATH              |
| 🛡️ **Security first**    | Zero plaintext storage, memory-only decryption           |

---

## 📦 Installation

### From releases

Download the latest binary for your platform from [**Releases**](https://github.com/wolfsouldev/ssh/releases):

| Platform        | Archive                          |
| --------------- | -------------------------------- |
| Windows (x64)   | `sshh_x.x.x_windows_amd64.zip`   |
| Windows (ARM64) | `sshh_x.x.x_windows_arm64.zip`   |
| Linux (x64)     | `sshh_x.x.x_linux_amd64.tar.gz`  |
| Linux (ARM64)   | `sshh_x.x.x_linux_arm64.tar.gz`  |
| macOS (x64)     | `sshh_x.x.x_darwin_amd64.tar.gz` |
| macOS (ARM64)   | `sshh_x.x.x_darwin_arm64.tar.gz` |

### From source

<details>
<summary><b>🪟 Windows (PowerShell)</b></summary>

```powershell
git clone https://github.com/wolfsouldev/ssh.git
cd sshh
go build -o sshh.exe .
```

</details>

<details>
<summary><b>🐧 Linux / 🍎 macOS (Bash)</b></summary>

```bash
git clone https://github.com/wolfsouldev/ssh.git
cd sshh
go build -o sshh .
```

</details>

### Self-install to PATH (Windows)

```powershell
.\sshh.exe install
# ✓ Installed to C:\Users\You\AppData\Local\sshh\bin\sshh.exe
# ✓ Added to user PATH (restart terminal to take effect)
```

### Add to PATH manually (Linux / macOS)

```bash
sudo cp sshh /usr/local/bin/
# — or —
mv sshh ~/.local/bin/   # Make sure ~/.local/bin is in your $PATH
```

---

## 🚀 Usage

### Connect to a server

```bash
sshh user@hostname           # Default port 22
sshh user@hostname -p 2222   # Custom port
sshh user@hostname:2222      # Port in target
```

### Connection flow

```
┌─────────────────────────┐
│   sshh user@hostname    │
└───────────┬─────────────┘
            │
     ┌──────▼──────┐
     │ Vault exists?│
     └──┬───────┬───┘
     Yes│       │No
        │       │
 ┌──────▼───┐   │  ┌──────────────┐
 │ Ask for  │   └──► Ask password │
 │ master pw│      │ & connect   │
 └──────┬───┘      └──────┬──────┘
        │                 │
  ┌─────▼─────┐    ┌──────▼──────┐
  │ Credential │    │ Save creds? │
  │  found?    │    │  Create     │
  └─┬──────┬──┘    │  vault?     │
  Yes│    No│       └─────────────┘
     │      │
┌────▼─┐ ┌──▼──────────┐
│ Auto │ │ Ask password │
│connect│ │ & offer save│
└──────┘ └─────────────┘
```

### Manage credentials

```bash
sshh list               # 📋 List all saved credentials
sshh add                # ➕ Add a new credential interactively
sshh edit <alias>       # ✏️  Edit an existing credential
sshh delete <alias>     # 🗑️  Delete a credential
```

### Master password

```bash
sshh master             # 🔄 Change your master password
```

### Export & Import

```bash
sshh export backup.vault   # 📤 Export vault (encrypted)
sshh import backup.vault   # 📥 Import vault from file
```

### Install & Uninstall (Windows)

```powershell
sshh install            # ⬇️  Install to PATH
sshh uninstall          # ⬆️  Remove from PATH
```

---

## 🔒 Security

| Layer                | Technology                 | Purpose                                     |
| -------------------- | -------------------------- | ------------------------------------------- |
| **Encryption**       | AES-256-GCM                | Authenticated symmetric encryption          |
| **Key Derivation**   | Argon2id                   | Memory-hard, GPU-resistant password hashing |
| **File Permissions** | 0600                       | Owner-only access to vault file             |
| **Storage (Win)**    | `%APPDATA%\sshh\vault.enc` | Encrypted vault, never plaintext            |
| **Storage (Linux)**  | `~/.config/sshh/vault.enc` | Same format, same encryption                |

**Principles:**

- Master password is **never** stored on disk
- Credentials are only decrypted **in memory** when needed
- Sensitive byte slices are zeroed after use
- No telemetry, no external network calls (except SSH connections)

---

## 🏗️ Project Structure

```
sshh/
├── 📄 main.go                      # Entry point
├── 📁 cmd/
│   ├── root.go                     # Root command & connect logic
│   ├── list.go                     # sshh list
│   ├── add.go                      # sshh add
│   ├── edit.go                     # sshh edit
│   ├── delete.go                   # sshh delete
│   ├── master.go                   # sshh master
│   ├── install.go                  # sshh install / uninstall
│   └── export_import.go            # sshh export / import
├── 📁 internal/
│   ├── 🔐 crypto/crypto.go        # AES-256-GCM + Argon2id
│   ├── 🗄️ vault/vault.go          # Credential CRUD operations
│   ├── 🌐 sshclient/client.go     # SSH connection (password & key)
│   ├── 📦 installer/installer.go   # Windows PATH installer
│   └── 🖥️ ui/prompt.go            # Terminal prompts & formatting
├── 📁 scripts/
│   └── setup-hooks.ps1             # Git hooks installer
├── 📁 .github/
│   ├── workflows/ci.yml            # CI pipeline
│   ├── PULL_REQUEST_TEMPLATE.md    # PR template
│   └── ISSUE_TEMPLATE/             # Bug & feature templates
├── .golangci.yml                   # Linter config
├── .goreleaser.yml                 # Release automation
├── .editorconfig                   # Editor settings
├── Makefile                        # Build commands
├── go.mod
└── go.sum
```

---

## 🛠️ Development

### Prerequisites

- **Go 1.25+**
- **Git**
- **golangci-lint** _(optional)_

### Quick start

```bash
git clone https://github.com/wolfsouldev/ssh.git
cd sshh
go mod download

# Windows
go build -o sshh.exe .

# Linux / macOS
go build -o sshh .
```

### Makefile commands

```bash
make build      # 🔨 Build the binary
make test       # 🧪 Run tests with coverage
make lint       # 🔍 Run golangci-lint
make fmt        # 💅 Format code (gofmt + goimports)
make vet        # 🩺 Run go vet
make all        # 🚀 Format + lint + test + build
```

### Git hooks

Set up pre-commit hooks (formatting, vetting, tests, secret scanning) and commit message validation:

```powershell
# Windows (PowerShell)
.\scripts\setup-hooks.ps1
```

---

## 🤝 Contributing

Contributions are what make the open source community amazing! Any contributions you make are **greatly appreciated**.

Please read our [**Contributing Guide**](CONTRIBUTING.md) for details on:

- Setting up the development environment
- Coding standards & Go style
- Commit message conventions (Conventional Commits)
- Pull request process

See the [open issues](https://github.com/wolfsouldev/ssh/issues) for a list of proposed features and known issues.

---

## 📄 License

Distributed under the **MIT License**. See [`LICENSE`](LICENSE) for more information.

---

## 🗺️ Roadmap

- [x] Core vault with AES-256-GCM encryption
- [x] Password & private key storage
- [x] Interactive credential management (add/edit/delete/list)
- [x] Auto-connect with saved credentials
- [x] Export/Import vault
- [x] Windows PATH installer
- [x] Git hooks & CI pipeline
- [ ] SSH agent forwarding support
- [ ] Credential groups & tags
- [ ] SSH config file (`~/.ssh/config`) import
- [ ] Multi-vault support
- [ ] Biometric unlock (Windows Hello)
- [ ] Linux & macOS full support
- [ ] TUI dashboard with [bubbletea](https://github.com/charmbracelet/bubbletea)

---

## 🙏 Acknowledgments

- [spf13/cobra](https://github.com/spf13/cobra) — CLI framework
- [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto) — Argon2id & SSH
- [golang.org/x/term](https://pkg.go.dev/golang.org/x/term) — Terminal handling
- [Contributor Covenant](https://www.contributor-covenant.org/) — Code of Conduct

---

<div align="center">

**[⬆ Back to top](#)**

Made with ❤️ by [wolfsouldev](https://github.com/wolfsouldev)

</div>
