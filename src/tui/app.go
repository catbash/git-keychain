package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"catbash/git-keychain/src/accounts"
	"catbash/git-keychain/src/keychain"
	"catbash/git-keychain/src/tui/constants"
	"catbash/git-keychain/src/tui/render"
	"catbash/git-keychain/src/tui/state"
)

type model struct {
	state.AppState
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.ShowHelp {
			switch msg.String() {
			case "esc", "q", "?":
				m.ShowHelp = false
			case "ctrl+c":
				return m, tea.Quit
			}
			return m, nil
		}
		if m.ShowSuccessAlert {
			switch msg.String() {
			case "ctrl+c", "esc", "o", "enter", "q":
				m.ShowSuccessAlert = false
			}
			return m, nil
		}
		if m.ShowErrorAlert {
			switch msg.String() {
			case "ctrl+c", "esc", "o", "enter", "q":
				m.ShowErrorAlert = false
			}
			return m, nil
		}
		if m.ShowDupAlert {
			switch msg.String() {
			case "ctrl+c", "esc", "o", "enter", "q":
				m.ShowDupAlert = false
			}
			return m, nil
		} else if m.ShowAlert {
			switch msg.String() {
			case "ctrl+c", "esc", "n", "q":
				m.ShowAlert = false
			case "y":
				m.ShowAlert = false
				if errMsg := keychain.ApplyAccount(m.Accounts[m.Cursor]); errMsg != "" {
					m.ShowErrorAlert = true
					m.ErrorMsg = errMsg
				} else {
					m.ShowSuccessAlert = true
				}
			case "enter":
				m.ShowAlert = false
				if m.AlertYes {
					if errMsg := keychain.ApplyAccount(m.Accounts[m.Cursor]); errMsg != "" {
						m.ShowErrorAlert = true
						m.ErrorMsg = errMsg
					} else {
						m.ShowSuccessAlert = true
					}
				}
			case "h", "left":
				m.AlertYes = false
			case "l", "right":
				m.AlertYes = true
			case "tab":
				m.AlertYes = !m.AlertYes
			}
			return m, nil
		} else if m.Commanding {
			switch msg.Type {
			case tea.KeyEsc:
				m.Commanding = false
				m.Command = ""
			case tea.KeyCtrlC:
				return m, tea.Quit
			case tea.KeyEnter:
				if m.Command == "q" {
					return m, tea.Quit
				}
				m.Commanding = false
				m.Command = ""
			case tea.KeyBackspace, tea.KeyDelete:
				if len(m.Command) > 0 {
					runes := []rune(m.Command)
					m.Command = string(runes[:len(runes)-1])
				} else {
					m.Commanding = false
				}
			case tea.KeyRunes:
				m.Command += string(msg.Runes)
			}
		} else if m.Searching {
			switch msg.Type {
			case tea.KeyEsc:
				m.Searching = false
			case tea.KeyCtrlC:
				return m, tea.Quit
			case tea.KeyBackspace, tea.KeyDelete:
				if len(m.Query) > 0 {
					runes := []rune(m.Query)
					m.Query = string(runes[:len(runes)-1])
					m.Accounts = accounts.FilterAccounts(m.AllAccounts, m.Query)
					m.Cursor = 0
				}
			case tea.KeyRunes:
				m.Query += string(msg.Runes)
				m.Accounts = accounts.FilterAccounts(m.AllAccounts, m.Query)
				m.Cursor = 0
			case tea.KeySpace:
				m.Query += " "
				m.Accounts = accounts.FilterAccounts(m.AllAccounts, m.Query)
				m.Cursor = 0
			}
		} else {
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "tab":
				m.RightFocused = !m.RightFocused
			case "h", "left":
				m.RightFocused = false
			case "l", "right":
				m.RightFocused = true
			case ":":
				m.Commanding = true
				m.Command = ""
			case "enter":
				if !m.RightFocused && len(m.Accounts) > 0 {
					if m.Duplicates[m.Accounts[m.Cursor].Alias] {
						m.ShowDupAlert = true
					} else {
						m.ShowAlert = true
						m.AlertYes = false
					}
				}
			case "s":
				m.Searching = true
				m.RightFocused = false
			case "?":
				m.ShowHelp = true
			case "j", "down":
				if m.RightFocused {
					if len(m.Accounts) > 0 {
						maxScroll := render.FieldRowCount(m.Accounts[m.Cursor]) - m.VisibleRows
						if maxScroll < 0 {
							maxScroll = 0
						}
						if m.ScrollOffset < maxScroll {
							m.ScrollOffset++
						}
					}
				} else if len(m.Accounts) > 0 {
					m.Cursor = (m.Cursor + 1) % len(m.Accounts)
					m.ScrollOffset = 0
				}
			case "k", "up":
				if m.RightFocused {
					if m.ScrollOffset > 0 {
						m.ScrollOffset--
					}
				} else if len(m.Accounts) > 0 {
					m.Cursor = (m.Cursor - 1 + len(m.Accounts)) % len(m.Accounts)
					m.ScrollOffset = 0
				}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	if len(m.Accounts) == 0 && !m.Searching && m.Query == "" {
		return "No accounts found.\n"
	}
	width := 80
	leftWidth := width / 3
	rightWidth := width - leftWidth
	return render.RenderView(m.AppState, width, constants.MaxHeight, leftWidth, rightWidth)
}
