package alerts

import (
	_ "embed"
	"strings"
)

//go:embed alerts.md/destructive.md
var DestructiveMD string

//go:embed alerts.md/duplicate.md
var DuplicateMD string

//go:embed alerts.md/success.md
var SuccessMD string

// ParseAlertMD splits a markdown string into (title, body).
// The first line starting with "# " is the title; the rest is the body.
func ParseAlertMD(md string) (title, body string) {
	lines := strings.Split(strings.TrimSpace(md), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "# ") {
			title = strings.TrimPrefix(line, "# ")
			body = strings.TrimSpace(strings.Join(lines[i+1:], "\n"))
			return
		}
	}
	body = strings.TrimSpace(md)
	return
}
