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
[![CI](https://img.shields.io/github/actions/workflow/status/wolfsouldev/ssh/ci.yml?branch=main&style=flat-square&label=CI)](https://github.com/wolfsouldev/ssh/actions)
[![Release](https://img.shields.io/github/v/release/wolfsouldev/ssh?style=flat-square&color=blue)](https://github.com/wolfsouldev/ssh/releases)
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
| 🖥️ **Interactive mode**  | Run `sshh` for a full menu: search, connect, manage      |
| 🔍 **Search**            | Type part of a name to find connections instantly        |
| 📋 **Full management**   | List, add, edit, delete credentials from the console     |
| 📦 **Export / Import**   | Backup and restore your encrypted vault                  |
| 🔄 **Self-update**       | `sshh update` updates to the latest version              |
| 🪟 **Cross-platform**    | Linux, macOS, and Windows support                        |
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

### Interactive mode

Just type `sshh` with no arguments for a full interactive experience:

```bash
sshh                    # 🖥️  Interactive menu: search, connect, manage
```

- **Type a number** to connect directly
- **Type a name** to search connections (e.g., `prox` finds `proxy`, `proxy-2`, etc.)
- **Press `a`** to add, **`e`** to edit, **`d`** to delete, **`q`** to quit

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

### Version & Update

```bash
sshh version            # ℹ️  Show version, commit, build date
sshh update             # 🔄 Self-update to the latest release
sshh update --check     # 🔍 Check for updates without installing
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
- No telemetry, no external network calls (except SSH connections and `sshh update`)

---

## 🏗️ Project Structure

```
sshh/
├── 📄 main.go                      # Entry point
├── 📁 cmd/
│   ├── root.go                     # Root command & connect logic
│   ├── interactive.go              # Interactive mode (menu, search)
│   ├── list.go                     # sshh list
│   ├── add.go                      # sshh add
│   ├── edit.go                     # sshh edit
│   ├── delete.go                   # sshh delete
│   ├── master.go                   # sshh master
│   ├── install.go                  # sshh install / uninstall
│   ├── export_import.go            # sshh export / import
│   ├── version.go                  # sshh version
│   └── update.go                   # sshh update (self-update)
├── 📁 internal/
│   ├── 🔐 crypto/crypto.go        # AES-256-GCM + Argon2id
│   ├── 🗄️ vault/vault.go          # Credential CRUD operations
│   ├── 🌐 sshclient/client.go     # SSH connection (password & key)
│   ├── 📦 installer/installer.go   # Windows PATH installer
│   ├── 🖥️ ui/prompt.go            # Terminal UI, ASCII art, colors
│   ├── 🖥️ ui/flush_unix.go        # Stdin flush (Linux/macOS)
│   ├── 🖥️ ui/flush_windows.go     # Stdin flush (Windows)
│   └── 📌 version/version.go      # Version info (ldflags injected)
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
cd ssh
go mod download
make build       # Builds with version info embedded
./sshh version   # Verify it works
```

### Makefile commands

```bash
make build      # 🔨 Build binary (with version/commit/date)
make test       # 🧪 Run tests with coverage
make lint       # 🔍 Run golangci-lint
make fmt        # 💅 Format code
make vet        # 🩺 Run go vet
make install    # 📦 Build and install to /usr/local/bin (Linux/macOS)
make uninstall  # �️  Remove from /usr/local/bin
make all        # �🚀 Format + lint + test + build
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

## � Updating

### Self-update (recommended)

```bash
sshh update              # Downloads and installs the latest release
sshh update --check      # Just check, don't install
```

On Linux, if sshh is installed in `/usr/local/bin`, you may need:

```bash
sudo sshh update
```

### Manual update

<details>
<summary><b>🐧 Linux</b></summary>

```bash
# Download the latest release
curl -sL https://github.com/wolfsouldev/ssh/releases/latest/download/sshh_$(curl -s https://api.github.com/repos/wolfsouldev/ssh/releases/latest | grep tag_name | cut -d'"' -f4 | tr -d v)_linux_amd64.tar.gz | tar xz

# Replace the binary
sudo mv sshh /usr/local/bin/sshh
```

</details>

<details>
<summary><b>🪟 Windows (PowerShell)</b></summary>

```powershell
# Download latest release from:
# https://github.com/wolfsouldev/ssh/releases/latest
# Extract and replace sshh.exe in your PATH location
```

</details>

<details>
<summary><b>🍎 macOS</b></summary>

```bash
curl -sL https://github.com/wolfsouldev/ssh/releases/latest/download/sshh_$(curl -s https://api.github.com/repos/wolfsouldev/ssh/releases/latest | grep tag_name | cut -d'"' -f4 | tr -d v)_darwin_arm64.tar.gz | tar xz
sudo mv sshh /usr/local/bin/sshh
```

</details>

### From source

```bash
git pull origin main
make build
sudo make install   # Linux/macOS
```

Your vault data is **never affected** by updates — it stays safely in its encrypted file.

---

## �🗺️ Roadmap

- [x] Core vault with AES-256-GCM encryption
- [x] Password & private key storage
- [x] Interactive credential management (add/edit/delete/list)
- [x] Auto-connect with saved credentials
- [x] Export/Import vault
- [x] Windows PATH installer
- [x] Git hooks & CI pipeline
- [x] Hacker-themed UI with ASCII art & ANSI colors
- [x] Interactive mode with connection search
- [x] Self-update from GitHub releases
- [x] Cross-platform support (Linux, macOS, Windows)
- [ ] SSH agent forwarding support
- [ ] Credential groups & tags
- [ ] SSH config file (`~/.ssh/config`) import
- [ ] Multi-vault support
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
