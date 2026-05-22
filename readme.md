A bash script to manage git accounts across multiple platforms.

```
keychain.sh

--file=/path/to/file       specify conf file to use
--alias=<alias>            set by the specified alias
--help                     show this help
```

Run with no `--alias` to list all configured accounts.

## Dependencies

Requires `yq` on `$PATH`

## Configuration

Accounts are configured in a YAML file (default: `./conf.yaml`, override with `--file=`). Each account requires `alias`, `username`, `email`, `host`, and `sshkey`. `note` is optional.

Note that the host for any account _must_ include the full TLD for the platform. For example:
- `github.com` -> OK
- `github` -> NO
- `codeberg.org` -> OK
- `codeberg` -> NO

`sshkey` is the key filename in `~/.ssh/`.

See `conf.example.yaml` for an example config file.

Schema is in `conf.schema.json` for editor validation.

## On Selection

When `--alias=<alias>` matches:

- `~/.gitconfig` is overwritten with the account's `username` and `email`. Existing file is backed up to `~/.gitconfig.<timestamp>.bak`.
- `~/.ssh/catbash/git-keychain.conf` is overwritten with a Host block for the account. Existing file is backed up to `~/.ssh/catbash/git-keychain.conf.<timestamp>.bak`. Directory is created if missing.

Include this line in `~/.ssh/config` to make sure all keychain switches are immediately applied to ssh:

```
Include catbash/git-keychain.conf
```
