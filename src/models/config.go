package models

import "fmt"

type Config struct {
	Accounts    []GitAccount `yaml:"accounts"`
	ColorActive string       `yaml:"color_active"`
	ColorMuted  string       `yaml:"color_muted"`
}

// Validate checks that the config matches the required schema.
func (c *Config) Validate() error {
	if len(c.Accounts) == 0 {
		return fmt.Errorf("invalid config: accounts must not be empty")
	}
	for i, a := range c.Accounts {
		if a.Alias == "" {
			return fmt.Errorf("invalid config: account[%d]: alias is required", i)
		}
		if a.Username == "" {
			return fmt.Errorf("invalid config: account[%d]: username is required", i)
		}
		if a.Email == "" {
			return fmt.Errorf("invalid config: account[%d]: email is required", i)
		}
		if a.Host == "" {
			return fmt.Errorf("invalid config: account[%d]: host is required", i)
		}
	}
	return nil
}
