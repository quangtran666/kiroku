package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/tranducquang/kiroku/internal/tui/keys"
	"github.com/tranducquang/kiroku/internal/tui/styles"
)

// DialogType represents the type of dialog
type DialogType int

const (
	DialogConfirm DialogType = iota
	DialogInput
	DialogSelect
)

// Dialog represents a modal dialog component
type Dialog struct {
	dialogType DialogType
	title      string
	message    string
	input      textinput.Model
	options    []string
	cursor     int
	visible    bool
	confirmed  bool
	width      int
	height     int
}

// NewDialog creates a new dialog component
func NewDialog() *Dialog {
	ti := textinput.New()
	ti.CharLimit = 100
	ti.Width = 40

	return &Dialog{
		input: ti,
	}
}

// ShowConfirm shows a confirmation dialog
func (d *Dialog) ShowConfirm(title, message string) {
	d.dialogType = DialogConfirm
	d.title = title
	d.message = message
	d.options = []string{"Yes", "No"}
	d.cursor = 1 // Default to "No"
	d.visible = true
	d.confirmed = false
}

// ShowInput shows an input dialog
func (d *Dialog) ShowInput(title, placeholder string) {
	d.dialogType = DialogInput
	d.title = title
	d.message = ""
	d.input.Placeholder = placeholder
	d.input.SetValue("")
	d.input.Focus()
	d.visible = true
	d.confirmed = false
}

// ShowSelect shows a selection dialog
func (d *Dialog) ShowSelect(title string, options []string) {
	d.dialogType = DialogSelect
	d.title = title
	d.options = options
	d.cursor = 0
	d.visible = true
	d.confirmed = false
}

// Hide hides the dialog
func (d *Dialog) Hide() {
	d.visible = false
	d.input.Blur()
}

// IsVisible returns whether the dialog is visible
func (d *Dialog) IsVisible() bool {
	return d.visible
}

// IsConfirmed returns whether the dialog was confirmed
func (d *Dialog) IsConfirmed() bool {
	return d.confirmed
}

// InputValue returns the input value
func (d *Dialog) InputValue() string {
	return d.input.Value()
}

// SelectedIndex returns the selected option index
func (d *Dialog) SelectedIndex() int {
	return d.cursor
}

// SelectedOption returns the selected option string
func (d *Dialog) SelectedOption() string {
	if d.cursor >= 0 && d.cursor < len(d.options) {
		return d.options[d.cursor]
	}
	return ""
}

// SetSize sets the dialog dimensions
func (d *Dialog) SetSize(width, height int) {
	d.width = width
	d.height = height
	d.input.Width = width - 10
}

// Update handles input
func (d *Dialog) Update(msg tea.Msg) (*Dialog, tea.Cmd) {
	if !d.visible {
		return d, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.DefaultKeyMap.Escape):
			d.Hide()
			return d, nil

		case key.Matches(msg, keys.DefaultKeyMap.Enter):
			if d.dialogType == DialogConfirm {
				d.confirmed = d.cursor == 0 // "Yes" is at index 0
			} else {
				d.confirmed = true
			}
			d.Hide()
			return d, nil

		case key.Matches(msg, keys.DefaultKeyMap.Left):
			if d.dialogType == DialogConfirm && d.cursor > 0 {
				d.cursor--
			}

		case key.Matches(msg, keys.DefaultKeyMap.Right):
			if d.dialogType == DialogConfirm && d.cursor < len(d.options)-1 {
				d.cursor++
			}

		case key.Matches(msg, keys.DefaultKeyMap.Up):
			if d.dialogType == DialogSelect && d.cursor > 0 {
				d.cursor--
			}

		case key.Matches(msg, keys.DefaultKeyMap.Down):
			if d.dialogType == DialogSelect && d.cursor < len(d.options)-1 {
				d.cursor++
			}
		}
	}

	// Update text input for input dialogs
	if d.dialogType == DialogInput {
		var cmd tea.Cmd
		d.input, cmd = d.input.Update(msg)
		return d, cmd
	}

	return d, nil
}

// View renders the dialog
func (d *Dialog) View() string {
	if !d.visible {
		return ""
	}

	var b strings.Builder

	// Title
	b.WriteString(styles.DialogTitleStyle.Render(d.title))
	b.WriteString("\n\n")

	// Content based on type
	switch d.dialogType {
	case DialogConfirm:
		b.WriteString(d.message)
		b.WriteString("\n\n")
		// Buttons
		for i, opt := range d.options {
			if i == d.cursor {
				b.WriteString(styles.ButtonFocusedStyle.Render(opt))
			} else {
				b.WriteString(styles.ButtonStyle.Render(opt))
			}
		}

	case DialogInput:
		b.WriteString(d.input.View())
		b.WriteString("\n\n")
		b.WriteString(styles.TextMuted.Render("Press Enter to confirm, Esc to cancel"))

	case DialogSelect:
		for i, opt := range d.options {
			if i == d.cursor {
				b.WriteString(styles.NoteItemSelectedStyle.Render("▸ " + opt))
			} else {
				b.WriteString(styles.NoteItemStyle.Render("  " + opt))
			}
			if i < len(d.options)-1 {
				b.WriteString("\n")
			}
		}
	}

	// Calculate position to center the dialog
	dialogWidth := 50
	if d.width > 0 && d.width < dialogWidth {
		dialogWidth = d.width - 4
	}

	return styles.DialogStyle.Width(dialogWidth).Render(b.String())
}
