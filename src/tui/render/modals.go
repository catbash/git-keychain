package render

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"catbash/git-keychain/src/tui/alerts"
	"catbash/git-keychain/src/tui/state"
	"catbash/git-keychain/src/tui/styles"
)

// RenderAlert renders the confirmation dialog box (no placement).
func RenderAlert(m state.AppState, tuiW int) string {
	dialogMaxW := tuiW - 8
	if dialogMaxW < 20 {
		dialogMaxW = 20
	}
	innerW := dialogMaxW - 6

	mdTitle, mdBody := alerts.ParseAlertMD(alerts.DestructiveMD)
	title := styles.AlertTitleStyle.Render(mdTitle)
	body := WordWrap(mdBody, innerW)

	noBtn := styles.AlertBtnStyle.Render("(N)o")
	yesBtn := styles.AlertBtnStyle.Render("(Y)es")
	if !m.AlertYes {
		noBtn = styles.AlertBtnActiveStyle.Render("(N)o")
	} else {
		yesBtn = styles.AlertBtnActiveStyle.Render("(Y)es")
	}
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, noBtn, lipgloss.NewStyle().Render("  "), yesBtn)

	content := lipgloss.JoinVertical(lipgloss.Center, title, "", body, "", buttons)
	return styles.AlertStyle.Width(dialogMaxW - 2).Render(content)
}

// RenderDupAlert renders the duplicate alias warning dialog (no placement).
func RenderDupAlert() string {
	mdTitle, mdBody := alerts.ParseAlertMD(alerts.DuplicateMD)
	title := styles.AlertTitleStyle.Render(mdTitle)
	okBtn := styles.AlertBtnActiveStyle.Render("(O)k")
	var content string
	if mdBody != "" {
		content = lipgloss.JoinVertical(lipgloss.Center, title, "", mdBody, "", okBtn)
	} else {
		content = lipgloss.JoinVertical(lipgloss.Center, title, "", okBtn)
	}
	contentW := 0
	for _, line := range strings.Split(content, "\n") {
		if w := lipgloss.Width(line); w > contentW {
			contentW = w
		}
	}
	return styles.AlertStyle.Width(contentW + 4).Render(content)
}

// RenderErrorAlert renders a single-button (O)k error dialog.
func RenderErrorAlert(msg string) string {
	title := styles.AlertTitleStyle.Render("ERROR")
	okBtn := styles.AlertBtnActiveStyle.Render("(O)k")
	body := WordWrap(msg, 40)
	content := lipgloss.JoinVertical(lipgloss.Center, title, "", body, "", okBtn)
	contentW := 0
	for _, line := range []string{title, body, okBtn} {
		if w := lipgloss.Width(line); w > contentW {
			contentW = w
		}
	}
	return styles.AlertStyle.Width(contentW + 4).Render(content)
}

// RenderSuccessAlert renders a single-button (O)k success dialog.
func RenderSuccessAlert() string {
	mdTitle, mdBody := alerts.ParseAlertMD(alerts.SuccessMD)
	title := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorActive).Render(mdTitle)
	okBtn := styles.AlertBtnActiveStyle.Render("(O)k")
	body := WordWrap(mdBody, 48)
	content := lipgloss.JoinVertical(lipgloss.Center, title, "", body, "", okBtn)
	contentW := 0
	for _, line := range strings.Split(content, "\n") {
		if w := lipgloss.Width(line); w > contentW {
			contentW = w
		}
	}
	return styles.AlertStyle.BorderForeground(styles.ColorActive).Width(contentW + 4).Render(content)
}

// RenderHelp renders the help panel showing all keybindings.
func RenderHelp(totalWidth, totalHeight int) string {
	type row struct{ key, desc string }
	sections := []struct {
		title string
		rows  []row
	}{
		{"NAVIGATION", []row{
			{"j / ↓", "Move down — cursor (left) or scroll (right)"},
			{"k / ↑", "Move up — cursor (left) or scroll (right)"},
			{"h / ←", "Focus left column"},
			{"l / →", "Focus right column"},
			{"tab", "Toggle column focus"},
			{"enter", "Select account"},
		}},
		{"MODES", []row{
			{"s", "Search"},
			{":", "Command prompt"},
		}},
		{"GENERAL", []row{
			{"?", "Toggle this help panel"},
			{"esc", "Cancel / close"},
			{"q", "Close alert"},
			{":q", "Quit"},
			{"ctrl+c", "Force quit"},
		}},
	}

	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true)
	descStyle := lipgloss.NewStyle().Foreground(styles.ColorMuted)
	sectionStyle := lipgloss.NewStyle().Foreground(styles.ColorActive).Bold(true).MarginTop(1)

	var lines []string
	lines = append(lines, styles.HeaderStyle.Render("KEYBINDINGS"))
	for _, s := range sections {
		lines = append(lines, sectionStyle.Render(s.title))
		for _, r := range s.rows {
			key := keyStyle.Width(12).Render(r.key)
			desc := descStyle.Render(r.desc)
			lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Top, key, desc))
		}
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	panel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorMuted).
		Padding(1, 2).
		Width(totalWidth - 2).
		Height(totalHeight - 3).
		Render(content)

	status := styles.StatusStyle.Width(totalWidth).Render("  esc / q  close help")
	return lipgloss.JoinVertical(lipgloss.Left, panel, status)
}
