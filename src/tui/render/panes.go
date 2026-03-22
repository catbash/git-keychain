package render

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"catbash/git-keychain/src/tui/state"
	"catbash/git-keychain/src/tui/styles"
)

// RenderView builds the full TUI string.
func RenderView(m state.AppState, totalWidth, totalHeight, leftWidth, rightWidth int) string {
	visibleRows := totalHeight - 6
	if visibleRows < 1 {
		visibleRows = 1
	}
	colPad := lipgloss.NewStyle().Padding(1)
	left := colPad.Render(RenderLeft(m, leftWidth-2))
	right := colPad.Render(RenderRight(m, rightWidth-2, visibleRows))

	panes := lipgloss.JoinHorizontal(lipgloss.Top, left, right)

	var hint string
	if m.Commanding {
		hint = ":" + m.Command
	} else if m.Searching {
		hint = "  type to filter   backspace delete   esc done"
	} else {
		hint = "  ↑/k up   ↓/j down   s search   :q quit"
	}
	status := styles.StatusStyle.Width(totalWidth).Render(hint)

	return lipgloss.JoinVertical(lipgloss.Left, panes, status)
}

// RenderLeft renders the alias list pane with search bar.
func RenderLeft(m state.AppState, width int) string {
	header := styles.HeaderStyle.Render("@GIT-KEYCHAIN")

	searchText := "/ " + m.Query
	var searchBar string
	if m.Searching {
		searchBar = styles.SearchBarActiveStyle.Width(width - 2).Render(searchText)
	} else {
		searchBar = styles.SearchBarStyle.Width(width - 2).Render(searchText)
	}

	rows := []string{header, searchBar}

	if len(m.Accounts) == 0 {
		rows = append(rows, styles.UnselectedStyle.Render("  no matches"))
		return lipgloss.JoinVertical(lipgloss.Left, rows...)
	}

	for i, a := range m.Accounts {
		label := a.Alias
		if m.Duplicates[a.Alias] {
			label = label + styles.DupSuffix
		}
		if i == m.Cursor {
			selColor := styles.ColorActive
			if m.RightFocused || m.Searching {
				selColor = styles.ColorMuted
			}
			rows = append(rows, styles.SelectedStyle.BorderForeground(selColor).Width(width-2).Render(label))
		} else {
			rows = append(rows, styles.UnselectedStyle.Width(width-3).Render(label))
		}
	}
	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

// RenderRight renders the account detail pane.
func RenderRight(m state.AppState, width, visibleRows int) string {
	header := styles.HeaderStyle.Render("DETAILS")
	if len(m.Accounts) == 0 {
		return lipgloss.JoinVertical(lipgloss.Left, header)
	}

	a := m.Accounts[m.Cursor]

	dash := "—"
	sshkey := a.SSHKey
	if sshkey == "" {
		sshkey = dash
	}
	note := strings.TrimRight(a.Note, "\n")
	if note == "" {
		note = dash
	}

	fields := [][2]string{
		{"alias", a.Alias},
		{"username", a.Username},
		{"email", a.Email},
		{"host", a.Host},
		{"sshkey", sshkey},
		{"note", note},
	}

	var fieldRows []string
	for i, f := range fields {
		fieldRows = append(fieldRows, styles.DetailLabelStyle.Render(f[0]))
		for _, line := range strings.Split(f[1], "\n") {
			fieldRows = append(fieldRows, styles.DetailValueStyle.Render(line))
		}
		if i < len(fields)-1 {
			fieldRows = append(fieldRows, "")
		}
	}

	start := m.ScrollOffset
	if start > len(fieldRows) {
		start = len(fieldRows)
	}
	end := start + visibleRows
	if end > len(fieldRows) {
		end = len(fieldRows)
	}
	visible := fieldRows[start:end]

	borderColor := styles.ColorMuted
	if m.RightFocused {
		borderColor = styles.ColorActive
	}
	contentBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1).
		Width(width - 4).
		Height(visibleRows).
		Render(lipgloss.JoinVertical(lipgloss.Left, visible...))

	return lipgloss.JoinVertical(lipgloss.Left, header, contentBox)
}
