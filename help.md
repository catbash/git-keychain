git-keychain — manage multiple git identities via SSH

USAGE
  git-keychain [flags]

FLAGS
  -h, --help             Show this help text
  -m, --mode <mode>      Launch a mode (details | lite)
  -a, --alias <alias>    Switch immediately to the named account
  -c, --config <path>    Path to config file (default: conf.yaml)

MODES
  lite      Inline search-and-apply picker (default)
  details   Interactive TUI — browse and switch git accounts

EXAMPLES
  git-keychain                   Open the lite picker
  git-keychain --mode lite       Open the lite picker
  git-keychain --mode details    Open the details TUI
  git-keychain --alias work      Switch to the 'work' account non-interactively

ALIAS OUTPUT
  On success:  Include catbash/git-keychain.conf
  On error:    ERROR: Duplicate
               ERROR: SSH key not found
