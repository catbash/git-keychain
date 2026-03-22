package models

type GitAccount struct {
	Alias    string `yaml:"alias"`
	Username string `yaml:"username"`
	Email    string `yaml:"email"`
	Host     string `yaml:"host"`
	SSHKey   string `yaml:"sshkey"`
	Note     string `yaml:"note"`
}
