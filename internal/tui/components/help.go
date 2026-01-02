package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/tranducquang/kiroku/internal/tui/keys"
	"github.com/tranducquang/kiroku/internal/tui/styles"
)

// Help represents the help overlay component
type Help struct {
	visible bool
	width   int
	height  int
}

// NewHelp creates a new help component
func NewHelp() *Help {
	return &Help{}
}

// Show shows the help overlay
func (h *Help) Show() {
	h.visible = true
}

// Hide hides the help overlay
func (h *Help) Hide() {
	h.visible = false
}

// Toggle toggles the help overlay
func (h *Help) Toggle() {
	h.visible = !h.visible
}

// IsVisible returns whether the help is visible
func (h *Help) IsVisible() bool {
	return h.visible
}

// SetSize sets the help dimensions
func (h *Help) SetSize(width, height int) {
	h.width = width
	h.height = height
}

// Update handles input
func (h *Help) Update(msg tea.Msg) (*Help, tea.Cmd) {
	if !h.visible {
		return h, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.DefaultKeyMap.Escape),
			key.Matches(msg, keys.DefaultKeyMap.Help):
			h.Hide()
		}
	}

	return h, nil
}

// View renders the help overlay
func (h *Help) View() string {
	if !h.visible {
		return ""
	}

	var b strings.Builder

	b.WriteString(styles.HelpTitleStyle.Render("⌨️  Keyboard Shortcuts"))
	b.WriteString("\n\n")

	sections := []struct {
		title string
		keys  []struct{ key, desc string }
	}{
		{
			title: "Navigation",
			keys: []struct{ key, desc string }{
				{"↑/k", "Move up"},
				{"↓/j", "Move down"},
				{"←/h", "Collapse/Left"},
				{"→/l", "Expand/Right"},
				{"Tab", "Switch panel"},
				{"Enter", "Select/Confirm"},
				{"Esc", "Back/Cancel"},
			},
		},
		{
			title: "Actions",
			keys: []struct{ key, desc string }{
				{"n", "New note"},
				{"t", "New todo"},
				{"f", "New folder"},
				{"e", "Edit note"},
				{"d", "Delete"},
				{"/", "Search"},
				{"s", "Toggle star"},
				{"x/Space", "Toggle done"},
				{"p", "Cycle priority"},
				{"m", "Move to folder"},
			},
		},
		{
			title: "Views",
			keys: []struct{ key, desc string }{
				{"?", "Toggle help"},
				{"v", "Toggle preview"},
				{"r", "Refresh"},
				{"q", "Quit"},
			},
		},
	}

	for _, section := range sections {
		b.WriteString(styles.TitleStyle.Render(section.title))
		b.WriteString("\n")
		for _, k := range section.keys {
			b.WriteString(styles.HelpKeyStyle.Render(k.key))
			b.WriteString(styles.HelpDescStyle.Render(k.desc))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	b.WriteString(styles.TextMuted.Render("Press ? or Esc to close"))

	// Calculate size
	helpWidth := 45
	helpHeight := 30

	return styles.HelpStyle.Width(helpWidth).Height(helpHeight).Render(b.String())
}
