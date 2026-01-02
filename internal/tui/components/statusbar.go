package components

import (
	"strings"

	"github.com/tranducquang/kiroku/internal/tui/styles"
)

// StatusBar represents the status bar component
type StatusBar struct {
	message string
	keys    []KeyHelp
	width   int
}

// KeyHelp represents a key help item
type KeyHelp struct {
	Key  string
	Desc string
}

// NewStatusBar creates a new status bar component
func NewStatusBar() *StatusBar {
	return &StatusBar{
		keys: defaultKeys(),
	}
}

func defaultKeys() []KeyHelp {
	return []KeyHelp{
		{Key: "n", Desc: "new"},
		{Key: "t", Desc: "todo"},
		{Key: "e", Desc: "edit"},
		{Key: "d", Desc: "delete"},
		{Key: "/", Desc: "search"},
		{Key: "?", Desc: "help"},
		{Key: "q", Desc: "quit"},
	}
}

// SetWidth sets the status bar width
func (s *StatusBar) SetWidth(width int) {
	s.width = width
}

// SetMessage sets a temporary message
func (s *StatusBar) SetMessage(message string) {
	s.message = message
}

// ClearMessage clears the message
func (s *StatusBar) ClearMessage() {
	s.message = ""
}

// SetKeys sets the key help items
func (s *StatusBar) SetKeys(keys []KeyHelp) {
	s.keys = keys
}

// ResetKeys resets to default keys
func (s *StatusBar) ResetKeys() {
	s.keys = defaultKeys()
}

// View renders the status bar
func (s *StatusBar) View() string {
	var b strings.Builder

	if s.message != "" {
		b.WriteString(s.message)
	} else {
		// Render key help
		for i, k := range s.keys {
			b.WriteString(styles.StatusKeyStyle.Render(k.Key))
			b.WriteString(" ")
			b.WriteString(styles.StatusDescStyle.Render(k.Desc))
			if i < len(s.keys)-1 {
				b.WriteString("  ")
			}
		}
	}

	return styles.StatusBarStyle.Width(s.width - 2).Render(b.String())
}
