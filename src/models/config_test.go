package models

import "testing"

func TestConfig_Validate_Valid(t *testing.T) {
	cfg := Config{
		Accounts: []GitAccount{
			{Alias: "work", Username: "alice", Email: "alice@example.com", Host: "github.com"},
		},
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected valid config to pass, got: %v", err)
	}
}

func TestConfig_Validate_NoAccounts(t *testing.T) {
	cfg := Config{}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for empty accounts, got nil")
	}
}

func TestConfig_Validate_MissingAlias(t *testing.T) {
	cfg := Config{
		Accounts: []GitAccount{
			{Alias: "", Username: "alice", Email: "alice@example.com", Host: "github.com"},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for missing alias, got nil")
	}
}

func TestConfig_Validate_MissingUsername(t *testing.T) {
	cfg := Config{
		Accounts: []GitAccount{
			{Alias: "work", Username: "", Email: "alice@example.com", Host: "github.com"},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for missing username, got nil")
	}
}

func TestConfig_Validate_MissingEmail(t *testing.T) {
	cfg := Config{
		Accounts: []GitAccount{
			{Alias: "work", Username: "alice", Email: "", Host: "github.com"},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for missing email, got nil")
	}
}

func TestConfig_Validate_MissingHost(t *testing.T) {
	cfg := Config{
		Accounts: []GitAccount{
			{Alias: "work", Username: "alice", Email: "alice@example.com", Host: ""},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for missing host, got nil")
	}
}
