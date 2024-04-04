A cli written in tcl to manage git accounts across multiple platforms.

```
commands:
help        display this help text
list        list all available accounts
set <index> list all available accounts and then set account by index
            (you may optionally specify the account index to skip
            listing all accounts)
```

Accounts are configured in the `accounts.conf` file. This file must be placed in the same directory as the tcl script. Accounts are configured with a username, email, and host. 

Note that the host for any account _must_ include the full hostname for the platform. For example:
- `github.com` -> OK
- `github` -> NO
- `codeberg.org` -> OK
- `codeberg` -> NO

Accounts may optionally include an ssh key filename if the name of your key file does not match the username of your account. The script assumes all ssh keys are in the `~/.ssh/` directory.

Accounts may also optionally include a note.

The script will overwrite the `~/.ssh/config` and `~/.gitconfig` files upon selection.
