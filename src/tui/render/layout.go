package render

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"catbash/git-keychain/src/models"
)

// FieldRowCount returns the number of rendered rows for an account's detail content.
func FieldRowCount(a models.GitAccount) int {
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
	count := 0
	for i, f := range fields {
		count++                                 // label
		count += len(strings.Split(f[1], "\n")) // value lines
		if i < len(fields)-1 {
			count++ // blank separator
		}
	}
	return count
}

// WordWrap wraps s to fit within width columns, breaking on word boundaries.
// Existing newlines in s are preserved.
func WordWrap(s string, width int) string {
	if width <= 0 {
		return s
	}
	var lines []string
	for _, para := range strings.Split(s, "\n") {
		if strings.TrimSpace(para) == "" {
			lines = append(lines, "")
			continue
		}
		var line strings.Builder
		lineLen := 0
		for _, word := range strings.Fields(para) {
			wLen := lipgloss.Width(word)
			if lineLen > 0 && lineLen+1+wLen > width {
				lines = append(lines, line.String())
				line.Reset()
				lineLen = 0
			}
			if lineLen > 0 {
				line.WriteByte(' ')
				lineLen++
			}
			line.WriteString(word)
			lineLen += wLen
		}
		if line.Len() > 0 {
			lines = append(lines, line.String())
		}
	}
	return strings.Join(lines, "\n")
}

// PlaceOverlay stamps fg centered on top of bg (bgW x bgH canvas).
// bg lines are assumed to be exactly bgW visible chars wide (as produced by lipgloss.Place).
func PlaceOverlay(bg, fg string, bgW, bgH int) string {
	fgLines := strings.Split(fg, "\n")
	bgLines := strings.Split(bg, "\n")
	for len(bgLines) < bgH {
		bgLines = append(bgLines, strings.Repeat(" ", bgW))
	}

	fgH := len(fgLines)
	fgW := 0
	for _, l := range fgLines {
		if w := lipgloss.Width(l); w > fgW {
			fgW = w
		}
	}

	startY := (bgH - fgH) / 2
	startX := (bgW - fgW) / 2
	if startX < 0 {
		startX = 0
	}
	if startY < 0 {
		startY = 0
	}
	for i, fgLine := range fgLines {
		bgY := startY + i
		if bgY < 0 || bgY >= len(bgLines) {
			continue
		}
		left := ansi.Truncate(bgLines[bgY], startX, "")
		right := ansi.TruncateLeft(bgLines[bgY], startX+lipgloss.Width(fgLine), "")
		bgLines[bgY] = left + fgLine + right
	}
	return strings.Join(bgLines, "\n")
}
