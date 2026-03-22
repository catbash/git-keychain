package styles

import "github.com/charmbracelet/lipgloss"

var (
	ColorActive = lipgloss.Color("10") // green — active selection, focused borders
	ColorMuted  = lipgloss.Color("8")  // grey — inactive items, unfocused borders
)

// ApplyColors overrides the default colors when the user specifies custom
// values in the config file. It must be called before the TUI is rendered.
func ApplyColors(active, muted string) {
	if active != "" {
		ColorActive = lipgloss.Color(active)
	}
	if muted != "" {
		ColorMuted = lipgloss.Color(muted)
	}
	// Rebuild styles that captured the color values at package init.
	SelectedStyle = lipgloss.NewStyle().
		BorderLeft(true).
		BorderStyle(lipgloss.Border{Left: "█"}).
		BorderForeground(ColorActive).
		Foreground(lipgloss.Color("15")).
		PaddingLeft(1)
	UnselectedStyle = lipgloss.NewStyle().
		Foreground(ColorMuted).
		PaddingLeft(3)
	DetailLabelStyle = lipgloss.NewStyle().
		Foreground(ColorMuted)
	StatusStyle = lipgloss.NewStyle().
		Foreground(ColorMuted).
		Reverse(true)
	SearchBarStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("15")).
		Foreground(ColorMuted)
	SearchBarActiveStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorActive).
		Foreground(lipgloss.Color("15"))
	AlertBtnStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#111111")).
		Background(ColorMuted).
		Padding(0, 1)
}

var (
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("15"))

	SelectedStyle = lipgloss.NewStyle().
			BorderLeft(true).
			BorderStyle(lipgloss.Border{Left: "█"}).
			BorderForeground(ColorActive).
			Foreground(lipgloss.Color("15")).
			PaddingLeft(1)

	UnselectedStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			PaddingLeft(3)

	DupSuffix = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Render(" *")

	DetailLabelStyle = lipgloss.NewStyle().
				Foreground(ColorMuted)

	DetailValueStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("15"))

	StatusStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Reverse(true)

	SearchBarStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("15")).
			Foreground(ColorMuted)

	SearchBarActiveStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorActive).
				Foreground(lipgloss.Color("15"))

	AlertStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("9")).
			AlignHorizontal(lipgloss.Center).
			Padding(1, 2)

	AlertTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("9"))

	AlertBtnStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#111111")).
			Background(ColorMuted).
			Padding(0, 1)

	AlertBtnActiveStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#111111")).
				Background(lipgloss.Color("15")).
				Padding(0, 1)
)
