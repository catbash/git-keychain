# git-keychain

A CLI tool for managing multiple git identities across platforms. Switch between accounts interactively or non-interactively. Overwrites your `~/.gitconfig` file with backup.

---

## Install

Clone the repo and build:

```sh
git clone <repo-url>
cd git-keychain
go build -o git-keychain .
```

Place the binary and `conf.yaml` in the same directory.

---

## Configuration

Create a `conf.yaml` file next to the binary. A schema file (`conf.schema.json`) is included for editor autocompletion.

```yaml
# yaml-language-server: $schema=conf.schema.json

accounts:
  - alias: test-github
    username: test
    email: test@test.com
    host: github.com
    # sshkey: test        # optional — defaults to username
    # note: ""            # optional — free-form description

  - alias: example-codeberg
    username: example
    email: example@example.com
    host: codeberg.org
```

### Account fields

| Field | Required | Description |
|-------|----------|-------------|
| `alias` | yes | Unique identifier. Alphanumeric, hyphens, underscores. |
| `username` | yes | Git display name written to `~/.gitconfig`. |
| `email` | yes | Git commit email written to `~/.gitconfig`. |
| `host` | yes | Full hostname of the platform (e.g. `github.com`, not `github`). |
| `sshkey` | no | SSH private key filename in `~/.ssh/`. Defaults to `username`. |
| `note` | no | Free-form description shown in the details view. |

### Global fields

| Field | Required | Description |
|-------|----------|-------------|
| `accounts` | yes | List of accounts (see above). |
| `color_active` | no | Hex color for active/highlighted elements. |
| `color_muted` | no | Hex color for muted/secondary text. |

---

## SSH setup

When an account is applied, git-keychain writes to `~/.ssh/catbash/git-keychain.conf`. Add the following line to your main SSH config (`~/.ssh/config`) so it takes effect:

```
Include catbash/git-keychain.conf
```

SSH keys are expected to live in `~/.ssh/`. The tool will error before writing any config if the key file does not exist.

---

## Usage

```
git-keychain [flags]
```

| Flag | Description |
|------|-------------|
| `-h`, `--help` | Show help text |
| `-m`, `--mode <mode>` | Launch a specific mode: `lite` or `details` |
| `-a`, `--alias <alias>` | Switch to an account non-interactively |
| `-c`, `--config <path>` | Path to config file (default: `conf.yaml` next to the binary) |

Running with no arguments launches the lite inline picker (default). Use `--mode details` for the full interactive TUI.

### Non-interactive use

Switch accounts from scripts or shell aliases without opening a TUI:

```sh
git-keychain --alias work
```

---

## Notes

- Accounts are sorted alphabetically by alias.
- Duplicate aliases are detected and blocked — the tool will refuse to apply a duplicate account.
- Applying an account overwrites `~/.gitconfig` and `~/.ssh/catbash/git-keychain.conf`. Both files are backed up to `*.bak` before each write.
