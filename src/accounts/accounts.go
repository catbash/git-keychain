package accounts

import (
	"strings"

	"catbash/git-keychain/src/models"
)

// BuildDuplicateSet returns a set of aliases that appear more than once.
func BuildDuplicateSet(accts []models.GitAccount) map[string]bool {
	freq := make(map[string]int, len(accts))
	for _, a := range accts {
		freq[a.Alias]++
	}
	dups := make(map[string]bool)
	for alias, count := range freq {
		if count > 1 {
			dups[alias] = true
		}
	}
	return dups
}

// FilterAccounts returns accounts where any field contains query (case-insensitive).
func FilterAccounts(accts []models.GitAccount, query string) []models.GitAccount {
	if query == "" {
		return accts
	}
	q := strings.ToLower(query)
	var result []models.GitAccount
	for _, a := range accts {
		if strings.Contains(strings.ToLower(a.Alias), q) ||
			strings.Contains(strings.ToLower(a.Username), q) ||
			strings.Contains(strings.ToLower(a.Email), q) ||
			strings.Contains(strings.ToLower(a.Host), q) ||
			strings.Contains(strings.ToLower(a.SSHKey), q) ||
			strings.Contains(strings.ToLower(a.Note), q) {
			result = append(result, a)
		}
	}
	return result
}
