#! /bin/bash

GITCONFIG="$HOME/.gitconfig"
KEYCHAIN_FOLDER="$HOME/.ssh/catbash/"
KEYCHAIN_FILE="$KEYCHAIN_FOLDER/git-keychain.conf"

dependencies=(
  yq
)

missing=()
for cmd in "${dependencies[@]}"; do
    if ! command -v "$cmd" >/dev/null 2>&1; then
        missing+=("$cmd")
    fi
done

if (( ${#missing[@]} > 0 )); then
    printf 'Missing dependencies:\n' >&2
    printf '  - %s\n' "${missing[@]}" >&2
    exit 1
fi

mkdir -p "$KEYCHAIN_FOLDER"

HELPTEXT="
keychain.sh

--file=/path/to/file       specify conf file to use
--alias=<alias>            set by the specified alias
--cleanup                  remove all backups except most recent
--help                     show this help
"

CONF_FILE="./conf.yaml"

write_git_config() {
    if [[ -f "$GITCONFIG" ]]; then
      cp "$GITCONFIG" "$GITCONFIG.$(date +%Y%m%d.%H%M%S).bak"
    fi
    cat > "$GITCONFIG" <<EOF
[user]
        name = $1
        email = $2
EOF
}

write_keychain_conf() {
    if [[ -f "$KEYCHAIN_FILE" ]]; then
      cp "$KEYCHAIN_FILE" "$KEYCHAIN_FILE.$(date +%Y%m%d.%H%M%S).bak"
    fi
    cat > "$KEYCHAIN_FILE" <<EOF
Host $1
	HostName $1
	User git
	IdentityFile ~/.ssh/$2
	IdentitiesOnly yes
EOF
}

print_accounts() {
  local i=0
  while IFS=$'\t' read -r alias username host; do
      echo "$i: alias=$alias uname=$username host=$host"
      ((i++))
  done < <(yq -r '.accounts[] | [.alias, .username, .host] | @tsv' "$CONF_FILE")
}

get_note() {
  alias="$1" yq '.accounts[] | select(.alias == strenv(alias)) | .note' "$CONF_FILE"
}

cleanup_backups() {
    local target="$1"
    local backups=()
    local f
    while IFS= read -r f; do
        backups+=("$f")
    done < <(find "$(dirname "$target")" -maxdepth 1 -name "$(basename "$target").*.bak" 2>/dev/null | sort)
    local n=${#backups[@]}
    if (( n <= 1 )); then
        return
    fi
    local i
    for (( i=0; i<n-1; i++ )); do
        rm -f -- "${backups[i]}"
        echo "removed ${backups[i]}"
    done
}

load_account_by_alias() {
  local data
  data=$(alias="$1" yq '.accounts[] | select(.alias == strenv(alias)) | [.username, .email, .alias, .host, .sshkey, .note] | @tsv' "$CONF_FILE")
  if [[ -z "$data" ]]; then
      printf 'Alias not found: %s\n' "$1" >&2
      exit 1
  fi
  IFS=$'\t' read -r ACTIVE_USERNAME ACTIVE_EMAIL ACTIVE_ALIAS ACTIVE_HOST ACTIVE_KEY ACTIVE_NOTE <<< "$data"
}

load_account_by_index() {
  local idx="$1"
  if ! [[ "$idx" =~ ^[0-9]+$ ]]; then
      printf 'Index must be a number: %s\n' "$idx" >&2
      exit 1
  fi
  local count
  count=$(yq '.accounts | length' "$CONF_FILE")
  if (( idx < 0 || idx >= count )); then
      printf 'Invalid index: %s (valid: 0-%d)\n' "$idx" "$((count-1))" >&2
      exit 1
  fi
  local data
  data=$(idx="$idx" yq '.accounts[env(idx) | tonumber] | [.username, .email, .alias, .host, .sshkey, .note] | @tsv' "$CONF_FILE")
  IFS=$'\t' read -r ACTIVE_USERNAME ACTIVE_EMAIL ACTIVE_ALIAS ACTIVE_HOST ACTIVE_KEY ACTIVE_NOTE <<< "$data"
}

for arg in "$@"; do
    case "$arg" in
        --file=*)
            CONF_FILE="${arg#--file=}"
            ;;
    esac
done

case "$CONF_FILE" in
    "~")    CONF_FILE="$HOME" ;;
    "~/"*)  CONF_FILE="$HOME/${CONF_FILE#"~/"}" ;;
esac

if [[ ! -f "$CONF_FILE" ]]; then
    printf 'Conf file not found: %s\n' "$CONF_FILE" >&2
    exit 1
fi

for arg in "$@"; do
    case "$arg" in
        --file=*)
            ;;
        --alias=*)
            load_account_by_alias "${arg#--alias=}"
            write_git_config "$ACTIVE_USERNAME" "$ACTIVE_EMAIL"
            write_keychain_conf "$ACTIVE_HOST" "$ACTIVE_KEY"
            echo "git keychain switched to $ACTIVE_USERNAME at $ACTIVE_HOST"
            if ! grep -qF "Include catbash/git-keychain.conf" "$HOME/.ssh/config" 2>/dev/null; then
                echo "Make sure to add this line to your ~/.ssh/config file"
                echo "Include catbash/git-keychain.conf"
            fi
            exit 0
            ;;
        --cleanup)
            cleanup_backups "$GITCONFIG"
            cleanup_backups "$KEYCHAIN_FILE"
            exit 0
            ;;
        --help|-h)
            printf '%s\n' "$HELPTEXT"
            exit 0
            ;;
        --*)
            printf 'Unknown argument: %s\n' "$arg" >&2
            printf '%s\n' "$HELPTEXT" >&2
            exit 1
            ;;
        *)
            load_account_by_alias "$arg"
            write_git_config "$ACTIVE_USERNAME" "$ACTIVE_EMAIL"
            write_keychain_conf "$ACTIVE_HOST" "$ACTIVE_KEY"
            echo "git keychain switched to $ACTIVE_USERNAME at $ACTIVE_HOST"
            if ! grep -qF "Include catbash/git-keychain.conf" "$HOME/.ssh/config" 2>/dev/null; then
                echo "Make sure to add this line to your ~/.ssh/config file"
                echo "Include catbash/git-keychain.conf"
            fi
            exit 0
            ;;
    esac
done

echo "git keychain"
print_accounts
