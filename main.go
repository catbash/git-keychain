package main

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"sort"

	"gopkg.in/yaml.v3"

	"catbash/git-keychain/src/accounts"
	"catbash/git-keychain/src/args"
	"catbash/git-keychain/src/keychain"
	"catbash/git-keychain/src/models"
	"catbash/git-keychain/src/tui/styles"
	modes "catbash/git-keychain/src"
)

//go:embed help.md
var helpMD string

func loadConfig(path string) (*models.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg models.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// handleAlias switches the active git account to the one named by alias.
// It writes output to w and returns an exit code (0 = success, 1 = error).
func handleAlias(w io.Writer, cfg *models.Config, accts []models.GitAccount, dups map[string]bool, alias string) int {
	// Duplicate check first — no filesystem writes
	if dups[alias] {
		fmt.Fprintln(w, "ERROR: Duplicate")
		return 1
	}
	// Find account
	var found *models.GitAccount
	for i := range accts {
		if accts[i].Alias == alias {
			found = &accts[i]
			break
		}
	}
	if found == nil {
		fmt.Fprintf(w, "ERROR: account %q not found\n", alias)
		return 1
	}
	// SSH key check before writing any config
	sshKey := found.SSHKey
	if sshKey == "" {
		sshKey = found.Username
	}
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(w, "ERROR: cannot determine home directory: %v\n", err)
		return 1
	}
	if !keychain.SSHKeyExists(home, sshKey) {
		fmt.Fprintln(w, "ERROR: SSH key not found")
		return 1
	}
	keychain.ApplyAccount(*found)
	fmt.Fprintln(w, "Include catbash/git-keychain.conf")
	return 0
}

func main() {
	parsed, err := args.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	configPath := "conf.yaml"
	if parsed.ConfigPath != "" {
		configPath = parsed.ConfigPath
	}

	cfg, err := loadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	styles.ApplyColors(cfg.ColorActive, cfg.ColorMuted)

	sort.Slice(cfg.Accounts, func(i, j int) bool {
		return cfg.Accounts[i].Alias < cfg.Accounts[j].Alias
	})
	dups := accounts.BuildDuplicateSet(cfg.Accounts)

	switch {
	case parsed.Help:
		fmt.Print(helpMD)
	case parsed.Alias != "":
		code := handleAlias(os.Stdout, cfg, cfg.Accounts, dups, parsed.Alias)
		os.Exit(code)
	case parsed.Mode == args.ModeDetails:
		modes.Run(cfg, dups)
	default:
		// ModeLite or no args — inline picker
		alias, confirmed := modes.RunLite(cfg, dups)
		if confirmed {
			code := handleAlias(os.Stdout, cfg, cfg.Accounts, dups, alias)
			os.Exit(code)
		}
	}
}
