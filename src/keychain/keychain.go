package keychain

import (
	"fmt"
	"os"
	"path/filepath"

	"catbash/git-keychain/src/models"
)

const SSHConfigPath = "catbash/git-keychain.conf"

// BackupFile copies src to src+".bak". Silently succeeds if src does not exist.
func BackupFile(src string) error {
	data, err := os.ReadFile(src)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return os.WriteFile(src+".bak", data, 0644)
}

// WriteGitConfig backs up and overwrites ~/.gitconfig with a [user] block.
func WriteGitConfig(home string, a models.GitAccount) error {
	path := filepath.Join(home, ".gitconfig")
	if err := BackupFile(path); err != nil {
		return err
	}
	content := fmt.Sprintf("[user]\n\tname = %s\n\temail = %s\n", a.Username, a.Email)
	return os.WriteFile(path, []byte(content), 0644)
}

// WriteSSHConfig backs up and overwrites ~/.ssh/<sshConfigFile> for the given account.
// Parent directories are created if they don't exist.
func WriteSSHConfig(home, sshConfigFile string, a models.GitAccount) error {
	sshKey := a.SSHKey
	if sshKey == "" {
		sshKey = a.Username
	}
	path := filepath.Join(home, ".ssh", sshConfigFile)
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	if err := BackupFile(path); err != nil {
		return err
	}
	content := fmt.Sprintf(
		"Host %s\n\tHostName %s\n\tUser git\n\tIdentityFile ~/.ssh/%s\n\tIdentitiesOnly yes\n",
		a.Host, a.Host, sshKey,
	)
	return os.WriteFile(path, []byte(content), 0644)
}

// SSHKeyExists reports whether ~/.ssh/<keyName> exists on disk.
func SSHKeyExists(home, keyName string) bool {
	path := filepath.Join(home, ".ssh", keyName)
	_, err := os.Stat(path)
	return err == nil
}

// ApplyAccount writes gitconfig and SSH config for the account.
// It returns a non-empty error message if the SSH key file is missing.
func ApplyAccount(a models.GitAccount) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Sprintf("cannot determine home directory: %v", err)
	}
	sshKey := a.SSHKey
	if sshKey == "" {
		sshKey = a.Username
	}
	WriteGitConfig(home, a)                    //nolint:errcheck — I/O errors surface as missing files
	WriteSSHConfig(home, SSHConfigPath, a)
	if !SSHKeyExists(home, sshKey) {
		return fmt.Sprintf("SSH key %q not found in ~/.ssh/\nPlease ensure the key file exists before using this account.", sshKey)
	}
	return ""
}
