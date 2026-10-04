[![Test and Build](https://github.com/SayHeyD/sops-age-manager/actions/workflows/test-and-build.yaml/badge.svg?branch=dev)](https://github.com/SayHeyD/sops-age-manager/actions/workflows/test-and-build.yaml) [![Lint](https://github.com/SayHeyD/sops-age-manager/actions/workflows/lint.yaml/badge.svg)](https://github.com/SayHeyD/sops-age-manager/actions/workflows/lint.yaml?branch=dev)

# sops-age-manager (sam)

**sam** is a cross-platform tool and system tray utility designed to seamlessly manage your [sops](https://github.com/getsops/sops) configuration when working with multiple [age](https://github.com/FiloSottile/age) encryption keys.

It eliminates the friction of manually passing `--age` recipient keys or managing `SOPS_AGE_KEY` environment variables when switching between multiple environments, teams, or Kubernetes namespaces.

---

## Table of Contents
- [Why isn't SOPS enough?](#why-isnt-sops-enough)
- [What does SAM do?](#what-does-sam-do)
- [Features](#features)
- [System Tray UI](#system-tray-ui)
- [User Guide](#user-guide)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Key Directory Setup](#key-directory-setup)
- [CLI Commands & Usage](#cli-commands--usage)
  - [Basic SOPS Execution](#basic-sops-execution)
  - [Key Management Commands](#key-management-commands)
  - [Configuration Commands](#configuration-commands)
  - [Command Documentation](#command-documentation)
- [Configuration](#configuration)

---

## Why isn't SOPS enough?

With standard SOPS tooling, switching age keys requires either passing the full public recipient key (`age1...`) as a command-line argument for every operation or setting the `SOPS_AGE_KEY` environment variable with the corresponding private key. 

Both methods are cumbersome, error-prone, and slow down everyday workflows when switching frequently between clusters, namespaces, or repositories.

---

## What does SAM do?

`sam` provides a smart, configurable layer on top of SOPS:
- **CLI Wrapper**: Intercepts commands passed after `--`, automatically setting `SOPS_AGE_KEY` for decryption and injecting `--age <public-key>` for encryption commands.
- **Native System Tray UI**: Runs in the background on macOS, Windows, and Linux to let you switch active keys, copy public recipient keys to the clipboard, or enter pass-through mode in a single click.
- **Live Watchers**: Automatically detects changes to your key directory (`~/.age/`) and configuration file, keeping all CLI and UI instances synchronized in real time.

---

## Features

- 🚀 **Desktop System Tray**: Fast key selection directly from the menu bar / system tray.
  - Currently only for MacOS and Windows 
- 🔄 **Live Dynamic Watchers**: Adding, removing, or modifying key files in `~/.age/` updates the UI immediately without restarting.
- ⌨️ **Shell Autocompletion**: Native dynamic tab completion for key names (`sam key use <TAB>`, `sam key copy <TAB>`), subcommands, and flags in Bash, Zsh, PowerShell, and Fish.
- 📋 **One-Click Clipboard Copy**: Copy public recipient keys (`age1...`), private keys, or key names to your clipboard from the CLI (`sam key copy`) or tray menu.
- 🔀 **Independent Key Selection**: Choose separate keys for Encryption, Decryption, or link them for Both.
- ⚡ **Pass-Through / Clear Mode**: Clear active keys (`sam key clear`) to revert to default SOPS behavior without closing SAM.
- 🪟 **Native Windows & macOS Integration**:
  - **Windows**: Dual-binary distribution with `sam.exe` (synchronous CLI console) and `samw.exe` (silent background tray launcher).
  - **macOS**: Menu bar application with dockless background execution.
- 📁 **Quick Directory Navigation**: Open configuration and key directories directly in Finder, File Explorer, or XDG.

---

## System Tray UI

Launch the UI by running `sam` without subcommands or flags:

```bash
sam
```

When launched, SAM sits in your system tray / menu bar with the following menu items:
- **Keys Submenu**: Lists all detected keys from `~/.age/` (including subdirectories).
  - **Both**: Set key for both encryption and decryption.
  - **Encryption**: Set key for encryption only (passed as `--age`).
  - **Decryption**: Set key for decryption only (passed as `SOPS_AGE_KEY`).
  - **Copy**: Copy key name, public key (`age1...`), or private key to clipboard.
- **Clear active keys**: Clears selected keys to operate in pass-through mode.
- **Open Config Directory**: Opens `~/.sops-age-manager/` in your native file manager.
- **Open Key Directory**: Opens `~/.age/` in your native file manager.

### Windows: `sam.exe` vs `samw.exe`
- **`sam.exe` (Console Executable)**: Use for command-line workflows in PowerShell, Command Prompt, Git Bash, and scripts. When run with no arguments, it launches the UI attached to your terminal session.
- **`samw.exe` (GUI Executable)**: Use for desktop shortcuts, Start Menu entries, or Windows Startup (`shell:startup`). It launches the system tray icon completely silently with no terminal window flashing.

---

## User Guide

### Prerequisites
- [sops](https://github.com/getsops/sops) installed and available in `$PATH`.
- Existing [age](https://github.com/FiloSottile/age) keys. (SAM manages existing keys; it does not generate new keys).

### Installation
Download the latest pre-compiled binary for your operating system and architecture from the [Releases](https://github.com/SayHeyD/sops-age-manager/releases) page on GitHub.

### Key Directory Setup
By default, SAM scans for age keys in `$HOME/.age/`.
- Key files must use the `.txt` extension (e.g. `production.txt`, `staging.txt`).
- The base filename without extension becomes the key identifier in SAM.
- Nested folders are supported (e.g. `$HOME/.age/k8s/dev.txt` is identified as `k8s/dev`).

A standard age key file contains comments, public key, and secret key:
```text
# created: 2026-01-01T00:00:00Z
# public key: age1ql3...
AGE-SECRET-KEY-1...
```

---

## CLI Commands & Usage

### Basic SOPS Execution

Use the `--` delimiter to wrap any command. SAM sets the required environment variables and arguments for the active key:

```bash
# Decrypt a file using the active decryption key:
sam -- sops -d super-secret.enc.yaml

# Encrypt a file using the active encryption key:
sam -- sops -e -i secret.yaml

# Works with any tool that consumes SOPS or SOPS_AGE_KEY:
sam -- terraform plan
```

### Key Management Commands

```bash
# List all available age keys and see which keys are active:
sam key list

# Select active keys for both encryption and decryption:
sam key use production-cluster

# Select a key for encryption only:
sam key use production-cluster -e

# Select a key for decryption only:
sam key use staging-cluster -d

# Copy the public key (age1...) recipient string to your clipboard:
sam key copy production-cluster

# Copy private key or key name to clipboard:
sam key copy production-cluster --private
sam key copy production-cluster --name

# Clear active keys to enter pass-through mode:
sam key clear
sam key clear -e    # Clear encryption key only
sam key clear -d    # Clear decryption key only
```

### Configuration Commands

```bash
# Display the path to the active configuration file:
sam config path

# Output the contents of the configuration file:
sam config dump
```

### Shell Autocompletion

SAM supports full shell autocompletion for subcommands, flags, and dynamic key name suggestions (`sam key use <TAB>`, `sam key copy <TAB>`).

Generate the autocompletion script for your preferred shell:

#### Bash
```bash
# Load in current session:
source <(sam completion bash)

# Load automatically for new sessions:
# Linux:
sam completion bash > /etc/bash_completion.d/sam
# macOS:
sam completion bash > $(brew --prefix)/etc/bash_completion.d/sam
```

#### Zsh
```bash
# Load in current session:
source <(sam completion zsh)

# Load automatically for new sessions:
sam completion zsh > "${fpath[1]}/_sam"
```

#### PowerShell (Windows / macOS / Linux)
```powershell
# Load in current session:
sam completion powershell | Out-String | Invoke-Expression

# Load automatically for new sessions (add to $PROFILE):
Add-Content -Path $PROFILE -Value "sam completion powershell | Out-String | Invoke-Expression"
```

#### Fish
```fish
# Load in current session:
sam completion fish | source

# Load automatically for new sessions:
sam completion fish > ~/.config/fish/completions/sam.fish
```

### Command Documentation

Detailed CLI documentation is available in the [`docs/`](./docs) directory:
- [SAM Base Command](./docs/sam.md)
  - [Completion](./docs/sam_completion.md)
  - [Config](./docs/sam_config.md)
    - [Dump](./docs/sam_config_dump.md)
    - [Path](./docs/sam_config_path.md)
  - [Key](./docs/sam_key.md)
    - [Clear](./docs/sam_key_clear.md)
    - [Copy](./docs/sam_key_copy.md)
    - [List](./docs/sam_key_list.md)
    - [Use](./docs/sam_key_use.md)

---

## Configuration

The default configuration file is created automatically at `$HOME/.sops-age-manager/config.yaml`.

```yaml
encryptionKey: "production-cluster"
decryptionKey: "production-cluster"
keyDir: ""
```

### Fields

- `encryptionKey`: The name of the key to use for encryption (injected as `--age <public-key>` into SOPS invocations). If empty, no `--age` flag is injected.
- `decryptionKey`: The name of the key to use for decryption (set as `SOPS_AGE_KEY` in the environment). If empty, `SOPS_AGE_KEY` is not set.
- `keyDir`: Custom absolute path to the directory containing `.txt` age key files. If left empty, defaults to `$HOME/.age/`.