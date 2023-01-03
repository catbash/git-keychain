Generates a script to toggle between multiple git accounts. Also rewrites `/home/$USER/.ssh/config` and generates a local `config` file for reference.

Usage:
- Configure your accounts in the `accounts.json` file
- Optionally backup your `/home/$USER/.ssh/config` file
- Run `setup.py`
- Add an alias for the generated script

Accounts can be added to `accounts.json` with the cli using the `-a` flag, or manually as follows:

``` JSON
{
  "name": "test",
  "email": "test@test.com",
  "platform": "github.com"
}
```

Note that the platform field _must_ include the full hostname for the platform. For example:
- `github.com` -> OK
- `github` -> NO
- `codeberg.org` -> OK
- `codeberg` -> NO

An optional `filename` field may be included with each account if the name of your key file does not match the name of your account. Otherwise, the script will default to using the account name. The script assumes all SSH keys are in your `/home/$USER/.ssh/` directory.

----------------------------

If you have two accounts belonging to the same platform, the generated script will use the first one with the normal config Host and all subsequent accounts will be registered with the account name as the config Host. For example, if you have the following accounts both at `github.com` ...

``` JSON
[
  {
    "name": "test",
    "email": "test@users.noreply.github.com",
    "platform": "github.com"
  },
  {
    "name": "example",
    "email": "example@users.noreply.github.com",
    "platform": "github.com"
  }
]
```

... you would be able to add an origin with the first account active with the normal `git remote add origin git@github.com:test/repo.git`, but you would have to use `git remote add origin git@example:example/repo.git` for the second.
