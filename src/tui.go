package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	tuipkg "catbash/git-keychain/src/tui"
	"catbash/git-keychain/src/models"
)

func Run(cfg *models.Config, dups map[string]bool) {
	m := tuipkg.New(cfg, dups)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// RunLite runs the lite inline TUI and returns (selectedAlias, wasConfirmed).
// Returns ("", false) if the user cancelled.
// Unlike Run, no WithAltScreen — renders inline in the terminal.
func RunLite(cfg *models.Config, dups map[string]bool) (string, bool) {
	m := tuipkg.NewLite(cfg, dups)
	p := tea.NewProgram(m) // no WithAltScreen
	fm, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	lm, ok := fm.(tuipkg.LiteModel)
	if !ok {
		fmt.Fprintln(os.Stderr, "error: unexpected model type")
		os.Exit(1)
	}
	return lm.Selected, lm.Confirmed
}