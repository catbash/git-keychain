package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"catbash/git-keychain/src/models"
)

// helpers

func makeAccounts(aliases ...string) []models.GitAccount {
	var out []models.GitAccount
	for _, a := range aliases {
		out = append(out, models.GitAccount{Alias: a})
	}
	return out
}

func liteWith(accs []models.GitAccount) LiteModel {
	return LiteModel{
		cfg:         &models.Config{Accounts: accs},
		allAccounts: accs,
		duplicates:  map[string]bool{},
		accounts:    accs,
	}
}

// --- counter ---

func TestLiteView_CounterShowsFilteredOfTotal(t *testing.T) {
	accs := makeAccounts("alpha", "beta", "gamma")
	m := liteWith(accs)
	// simulate a query that narrows to 2 results
	m.query = "a"
	m.accounts = makeAccounts("alpha", "gamma") // filtered
	out := m.View()
	if !strings.Contains(out, "Showing 2 of 3 aliases") {
		t.Errorf("want counter 'Showing 2 of 3 aliases', got:\n%s", out)
	}
}

func TestLiteView_CounterEmptyQuery(t *testing.T) {
	accs := makeAccounts("alpha", "beta", "gamma")
	m := liteWith(accs)
	out := m.View()
	if !strings.Contains(out, "Showing 3 of 3 aliases") {
		t.Errorf("want counter 'Showing 3 of 3 aliases' for empty query, got:\n%s", out)
	}
}

func TestLiteView_CounterSingular(t *testing.T) {
	accs := makeAccounts("alpha", "beta", "gamma")
	m := liteWith(accs)
	m.accounts = makeAccounts("alpha") // 1 filtered result
	out := m.View()
	if !strings.Contains(out, "Showing 1 of 3 alias") {
		t.Errorf("want singular 'alias', got:\n%s", out)
	}
	if strings.Contains(out, "Showing 1 of 3 aliases") {
		t.Errorf("should not use plural 'aliases' for count of 1, got:\n%s", out)
	}
}

func TestLiteView_CounterZeroResults(t *testing.T) {
	accs := makeAccounts("alpha", "beta")
	m := liteWith(accs)
	m.query = "zzz"
	m.accounts = nil
	out := m.View()
	if !strings.Contains(out, "Showing 0 of 2 aliases") {
		t.Errorf("want 'Showing 0 of 2 aliases', got:\n%s", out)
	}
}

// --- scroll indicators ---

func TestLiteView_NoIndicators_FourOrFewer(t *testing.T) {
	accs := makeAccounts("alpha", "beta", "gamma", "delta")
	m := liteWith(accs)
	out := m.View()
	if strings.Contains(out, "▲") {
		t.Errorf("no ▲ expected with exactly 4 items at top, got:\n%s", out)
	}
	if strings.Contains(out, "▼") {
		t.Errorf("no ▼ expected with exactly 4 items at top, got:\n%s", out)
	}
}

func TestLiteView_DownIndicator_WhenMoreItemsBelow(t *testing.T) {
	accs := makeAccounts("a", "b", "c", "d", "e")
	m := liteWith(accs) // visibleStart=0, window shows a-d, e is below
	out := m.View()
	if !strings.Contains(out, "▼") {
		t.Errorf("expected ▼ when items below window, got:\n%s", out)
	}
	if strings.Contains(out, "▲") {
		t.Errorf("no ▲ expected at top, got:\n%s", out)
	}
}

func TestLiteView_UpIndicator_WhenItemsAbove(t *testing.T) {
	accs := makeAccounts("a", "b", "c", "d", "e")
	m := liteWith(accs)
	m.visibleStart = 1 // items above
	m.cursor = 1
	out := m.View()
	if !strings.Contains(out, "▲") {
		t.Errorf("expected ▲ when items above window, got:\n%s", out)
	}
}

func TestLiteView_BothIndicators_InMiddle(t *testing.T) {
	accs := makeAccounts("a", "b", "c", "d", "e", "f")
	m := liteWith(accs)
	m.visibleStart = 1 // items above and below
	m.cursor = 1
	out := m.View()
	if !strings.Contains(out, "▲") {
		t.Errorf("expected ▲ in middle, got:\n%s", out)
	}
	if !strings.Contains(out, "▼") {
		t.Errorf("expected ▼ in middle, got:\n%s", out)
	}
}

func TestLiteView_NoIndicators_ZeroResults(t *testing.T) {
	accs := makeAccounts("alpha")
	m := liteWith(accs)
	m.query = "zzz"
	m.accounts = nil
	out := m.View()
	if strings.Contains(out, "▲") || strings.Contains(out, "▼") {
		t.Errorf("no indicators expected when 0 results, got:\n%s", out)
	}
}

// --- scroll window (4 items max) ---

func TestLiteView_AtMostFourRowsVisible(t *testing.T) {
	// Use aliases that don't appear in the hint bar text
	accs := makeAccounts("foo1", "foo2", "foo3", "foo4", "foo5", "foo6", "foo7")
	m := liteWith(accs)
	out := m.View()
	// Count how many alias names appear
	count := 0
	for _, alias := range []string{"foo1", "foo2", "foo3", "foo4", "foo5", "foo6", "foo7"} {
		if strings.Contains(out, alias) {
			count++
		}
	}
	if count > 4 {
		t.Errorf("expected at most 4 visible rows, got %d visible in:\n%s", count, out)
	}
}

func TestLiteView_BottomWindow_LastItemVisible(t *testing.T) {
	accs := makeAccounts("a", "b", "c", "d", "e")
	m := liteWith(accs)
	// scroll to bottom: visibleStart = 5 - 4 = 1
	m.visibleStart = 1
	m.cursor = 4
	out := m.View()
	if !strings.Contains(out, "e") {
		t.Errorf("last item 'e' should be visible at bottom of window, got:\n%s", out)
	}
	if strings.Contains(out, "▼") {
		t.Errorf("no ▼ at bottom of list, got:\n%s", out)
	}
}

// --- cursor movement adjusts visibleStart ---

func TestLiteUpdate_CursorDownScrollsWindow(t *testing.T) {
	accs := makeAccounts("a", "b", "c", "d", "e")
	m := liteWith(accs)
	m.cursor = 3 // at bottom of window (visibleStart=0, window=0..3)
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	got := updated.(LiteModel)
	if got.cursor != 4 {
		t.Errorf("want cursor 4, got %d", got.cursor)
	}
	// window should have shifted: visibleStart = cursor - visibleCount + 1 = 4 - 4 + 1 = 1
	if got.visibleStart != 1 {
		t.Errorf("want visibleStart 1, got %d", got.visibleStart)
	}
}

func TestLiteUpdate_CursorUpScrollsWindowUp(t *testing.T) {
	accs := makeAccounts("a", "b", "c", "d", "e")
	m := liteWith(accs)
	m.cursor = 1
	m.visibleStart = 1 // window shows items 1-4, cursor at item 1 (top of window)
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	got := updated.(LiteModel)
	if got.cursor != 0 {
		t.Errorf("want cursor 0, got %d", got.cursor)
	}
	// cursor < visibleStart → visibleStart = cursor = 0
	if got.visibleStart != 0 {
		t.Errorf("want visibleStart 0, got %d", got.visibleStart)
	}
}

// --- query mutation resets visibleStart ---

func TestLiteUpdate_QueryMutation_ResetsVisibleStart(t *testing.T) {
	accs := makeAccounts("a", "b", "c", "d", "e")
	m := liteWith(accs)
	m.visibleStart = 2
	m.cursor = 3

	// KeyRunes
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})
	got := updated.(LiteModel)
	if got.visibleStart != 0 {
		t.Errorf("KeyRunes: want visibleStart 0, got %d", got.visibleStart)
	}

	// KeySpace
	m.visibleStart = 2
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeySpace})
	got = updated.(LiteModel)
	if got.visibleStart != 0 {
		t.Errorf("KeySpace: want visibleStart 0, got %d", got.visibleStart)
	}

	// KeyBackspace
	m.visibleStart = 2
	m.query = "abc"
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	got = updated.(LiteModel)
	if got.visibleStart != 0 {
		t.Errorf("KeyBackspace: want visibleStart 0, got %d", got.visibleStart)
	}
}

// --- no matches row ---

func TestLiteView_NoMatchesRow(t *testing.T) {
	accs := makeAccounts("alpha")
	m := liteWith(accs)
	m.query = "zzz"
	m.accounts = nil
	out := m.View()
	if !strings.Contains(out, "no matches") {
		t.Errorf("expected 'no matches' row, got:\n%s", out)
	}
}

// --- hint bar ---

func TestLiteView_HintBarAlwaysVisible(t *testing.T) {
	accs := makeAccounts("alpha")
	m := liteWith(accs)
	out := m.View()
	if !strings.Contains(out, "navigate") {
		t.Errorf("hint bar should always be visible, got:\n%s", out)
	}
}

func TestLiteView_HintBarVisibleWhenEmpty(t *testing.T) {
	accs := makeAccounts("alpha")
	m := liteWith(accs)
	m.query = "zzz"
	m.accounts = nil
	out := m.View()
	if !strings.Contains(out, "navigate") {
		t.Errorf("hint bar should be visible even with 0 results, got:\n%s", out)
	}
}
