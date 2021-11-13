Generates a script to toggle between multiple git accounts.
- Configure accounts in `accounts.json`.
- Add desired alias for generated `toggle.sh` to `.bashrc` and run `source $HOME/.bashrc`.
- If updating, run `source $HOME/.bashrc` on run of `setup.py`.

`.ssh/config`
```
Host test
  HostName github.com
  User git
  IdentityFile /home/{user}/.ssh/test
  IdentitiesOnly yes
# to use this identity, remote has to be set to git@test:test/test.git

Host example
  HostName github.com
  User git
  IdentityFile /home/{user}/.ssh/example
  IdentitiesOnly yes
# to use this identity, remote has to be set to git@example:example/test.git
```

Config entry template:
```
Host {name}
  HostName github.com
  User git
  IdentityFile /home/{user}/.ssh/{name}
  IdentitiesOnly yes
# to use this identity, remote has to be set to git@{name}:{name}/test.git
```
