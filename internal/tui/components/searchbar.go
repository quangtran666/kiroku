package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/tranducquang/kiroku/internal/tui/styles"
)

// SearchBar represents the search bar component
type SearchBar struct {
	input  textinput.Model
	active bool
	height int
	width  int
}

// NewSearchBar creates a new search bar component
func NewSearchBar() *SearchBar {
	ti := textinput.New()
	ti.Placeholder = "Search notes..."
	ti.CharLimit = 100
	ti.Width = 30

	return &SearchBar{
		input: ti,
	}
}

// SetSize sets the search bar dimensions
func (s *SearchBar) SetSize(width, height int) {
	s.width = width
	s.height = height
	s.input.Width = width - 10
}

// Focus focuses the search bar
func (s *SearchBar) Focus() {
	s.active = true
	s.input.Focus()
}

// Blur removes focus from the search bar
func (s *SearchBar) Blur() {
	s.active = false
	s.input.Blur()
}

// IsActive returns whether the search bar is active
func (s *SearchBar) IsActive() bool {
	return s.active
}

// Value returns the current search value
func (s *SearchBar) Value() string {
	return s.input.Value()
}

// SetValue sets the search value
func (s *SearchBar) SetValue(value string) {
	s.input.SetValue(value)
}

// Clear clears the search bar
func (s *SearchBar) Clear() {
	s.input.SetValue("")
}

// Update handles input
func (s *SearchBar) Update(msg tea.Msg) (*SearchBar, tea.Cmd) {
	if !s.active {
		return s, nil
	}

	var cmd tea.Cmd
	s.input, cmd = s.input.Update(msg)
	return s, cmd
}

// View renders the search bar
func (s *SearchBar) View() string {
	icon := styles.SearchIconStyle.Render("🔍 ")

	style := styles.SearchBarStyle.Width(s.width - 4)
	if s.active {
		style = style.BorderForeground(styles.Primary)
	}

	return style.Render(icon + s.input.View())
}
