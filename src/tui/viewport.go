package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"catbash/git-keychain/src/models"
	"catbash/git-keychain/src/tui/constants"
	"catbash/git-keychain/src/tui/render"
	"catbash/git-keychain/src/tui/state"
)

type viewportModel struct {
	appModel model
	width    int
	height   int
}

// New returns a tea.Model ready to run.
func New(cfg *models.Config, dups map[string]bool) tea.Model {
	return viewportModel{
		appModel: model{
			AppState: state.AppState{
				Cfg:         cfg,
				AllAccounts: cfg.Accounts,
				Duplicates:  dups,
				Accounts:    cfg.Accounts,
				Cursor:      0,
				VisibleRows: constants.MaxHeight - 6,
			},
		},
	}
}

func (m viewportModel) Init() tea.Cmd { return nil }

func (m viewportModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		h := m.height
		if h == 0 || h > constants.MaxHeight {
			h = constants.MaxHeight
		}
		m.appModel.VisibleRows = h - 6
		if m.appModel.VisibleRows < 1 {
			m.appModel.VisibleRows = 1
		}
		return m, nil
	default:
		updated, cmd := m.appModel.Update(msg)
		m.appModel = updated.(model)
		return m, cmd
	}
}

func (m viewportModel) View() string {
	termW := m.width
	termH := m.height
	if termW == 0 {
		termW = 80
	}
	if termH == 0 {
		termH = 24
	}

	if termW < constants.MinWidth || termH < constants.MinHeight {
		if termH < 1 || termW < lipgloss.Width(constants.SmallMsg) {
			return ""
		}
		return lipgloss.Place(termW, termH, lipgloss.Center, lipgloss.Center, constants.SmallMsg)
	}

	w := termW
	if w > constants.MaxWidth {
		w = constants.MaxWidth
	}
	h := termH
	if h > constants.MaxHeight {
		h = constants.MaxHeight
	}

	leftWidth := w / 3
	rightWidth := w - leftWidth

	if m.appModel.ShowHelp {
		content := render.RenderHelp(w, h)
		return lipgloss.Place(termW, termH, lipgloss.Center, lipgloss.Center, content)
	}

	if len(m.appModel.Accounts) == 0 && !m.appModel.Searching && m.appModel.Query == "" && !m.appModel.ShowAlert {
		return lipgloss.Place(termW, termH, lipgloss.Center, lipgloss.Center, "No accounts found.")
	}

	content := render.RenderView(m.appModel.AppState, w, h, leftWidth, rightWidth)
	bg := lipgloss.Place(termW, termH, lipgloss.Center, lipgloss.Center, content)

	if m.appModel.ShowDupAlert {
		dialog := render.RenderDupAlert()
		return render.PlaceOverlay(bg, dialog, termW, termH)
	}

	if m.appModel.ShowSuccessAlert {
		dialog := render.RenderSuccessAlert()
		return render.PlaceOverlay(bg, dialog, termW, termH)
	}

	if m.appModel.ShowErrorAlert {
		dialog := render.RenderErrorAlert(m.appModel.ErrorMsg)
		return render.PlaceOverlay(bg, dialog, termW, termH)
	}

	if m.appModel.ShowAlert {
		dialog := render.RenderAlert(m.appModel.AppState, w)
		return render.PlaceOverlay(bg, dialog, termW, termH)
	}

	return bg
}
