package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"catbash/git-keychain/src/accounts"
	"catbash/git-keychain/src/models"
	"catbash/git-keychain/src/tui/styles"
)

// LiteModel is the inline search-and-apply TUI.
// Selected and Confirmed are read by the caller after the program exits.
type LiteModel struct {
	cfg         *models.Config
	allAccounts []models.GitAccount
	duplicates  map[string]bool
	query        string
	accounts     []models.GitAccount // filtered view
	cursor       int
	visibleStart int

	Selected  string // alias chosen by user
	Confirmed bool   // true if user pressed Enter
}

// NewLite returns a LiteModel ready for tea.NewProgram.
func NewLite(cfg *models.Config, dups map[string]bool) tea.Model {
	return LiteModel{
		cfg:         cfg,
		allAccounts: cfg.Accounts,
		duplicates:  dups,
		accounts:    cfg.Accounts,
	}
}

func (m LiteModel) Init() tea.Cmd { return nil }

func (m LiteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyRunes:
			m.query += string(msg.Runes)
			m.accounts = accounts.FilterAccounts(m.allAccounts, m.query)
			m.cursor = 0
			m.visibleStart = 0
		case tea.KeySpace:
			m.query += " "
			m.accounts = accounts.FilterAccounts(m.allAccounts, m.query)
			m.cursor = 0
			m.visibleStart = 0
		case tea.KeyBackspace, tea.KeyDelete:
			if len(m.query) > 0 {
				runes := []rune(m.query)
				m.query = string(runes[:len(runes)-1])
				m.accounts = accounts.FilterAccounts(m.allAccounts, m.query)
				m.cursor = 0
				m.visibleStart = 0
			}
		case tea.KeyUp:
			if len(m.accounts) > 0 && m.cursor > 0 {
				m.cursor--
				if m.cursor < m.visibleStart {
					m.visibleStart = m.cursor
				}
			}
		case tea.KeyDown:
			if len(m.accounts) > 0 && m.cursor < len(m.accounts)-1 {
				m.cursor++
				if m.cursor >= m.visibleStart+visibleCount {
					m.visibleStart = m.cursor - visibleCount + 1
				}
			}
		case tea.KeyEnter:
			if len(m.accounts) > 0 {
				m.Selected = m.accounts[m.cursor].Alias
				m.Confirmed = true
				return m, tea.Quit
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}
	return m, nil
}

const liteWidth = 40
const visibleCount = 4

func (m LiteModel) View() string {
	// Search input with border
	cursor := ""
	if m.query == "" {
		cursor = "_"
	}
	prompt := "> " + m.query + cursor
	inner := lipgloss.NewStyle().Width(liteWidth - 4).Render(prompt)
	inputBox := styles.SearchBarActiveStyle.Width(liteWidth - 2).Render(inner)

	muted := lipgloss.NewStyle().Foreground(styles.ColorMuted)

	noun := "aliases"
	if len(m.accounts) == 1 {
		noun = "alias"
	}
	counter  := muted.Render(fmt.Sprintf("Showing %d of %d %s", len(m.accounts), len(m.allAccounts), noun))
	scrollUp := muted.Render("▲")
	scrollDn := muted.Render("▼")
	hint     := muted.Render("↑↓ navigate · enter apply · esc cancel")

	// Defensive clamp
	if m.visibleStart >= len(m.accounts) {
		m.visibleStart = max(0, len(m.accounts)-visibleCount)
	}

	// Windowed slice
	end := m.visibleStart + visibleCount
	if end > len(m.accounts) {
		end = len(m.accounts)
	}
	visible := m.accounts[m.visibleStart:end]

	// Build alias rows
	var rows []string
	for i, a := range visible {
		alias := a.Alias
		if m.duplicates[alias] {
			alias += styles.DupSuffix
		}
		absIdx := m.visibleStart + i
		if absIdx == m.cursor {
			rows = append(rows, styles.SelectedStyle.Render(alias))
		} else {
			rows = append(rows, styles.UnselectedStyle.Render(alias))
		}
	}

	parts := []string{inputBox, counter}

	if len(m.accounts) > 0 {
		if m.visibleStart > 0 {
			parts = append(parts, scrollUp)
		}
		parts = append(parts, strings.Join(rows, "\n"))
		if m.visibleStart+visibleCount < len(m.accounts) {
			parts = append(parts, scrollDn)
		}
	} else if m.query != "" {
		parts = append(parts, styles.UnselectedStyle.Render("no matches"))
	}

	parts = append(parts, "", hint)

	return lipgloss.NewStyle().PaddingLeft(2).Render(strings.Join(parts, "\n"))
}
