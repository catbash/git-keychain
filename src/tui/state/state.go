package state

import "catbash/git-keychain/src/models"

// AppState holds all TUI application state. It is defined here so that
// src/tui/render can accept it as a parameter without importing src/tui.
type AppState struct {
	AllAccounts      []models.GitAccount
	Duplicates       map[string]bool
	Searching        bool
	Query            string
	Commanding       bool
	Command          string
	RightFocused     bool
	ScrollOffset     int
	VisibleRows      int
	ShowAlert        bool
	AlertYes         bool
	ShowDupAlert     bool
	ShowErrorAlert   bool
	ErrorMsg         string
	ShowSuccessAlert bool
	ShowHelp         bool
	Cfg              *models.Config
	Accounts         []models.GitAccount
	Cursor           int
}
