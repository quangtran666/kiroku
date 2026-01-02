package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/tranducquang/kiroku/internal/models"
	"github.com/tranducquang/kiroku/internal/tui/keys"
	"github.com/tranducquang/kiroku/internal/tui/styles"
)

// NoteList represents the note list component
type NoteList struct {
	notes      []*models.Note
	cursor     int
	height     int
	width      int
	focused    bool
	showTodos  bool
	folderName string
}

// NewNoteList creates a new note list component
func NewNoteList() *NoteList {
	return &NoteList{
		notes: make([]*models.Note, 0),
	}
}

// SetNotes sets the notes to display
func (n *NoteList) SetNotes(notes []*models.Note) {
	n.notes = notes
	if n.cursor >= len(notes) {
		n.cursor = len(notes) - 1
	}
	if n.cursor < 0 {
		n.cursor = 0
	}
}

// SetSize sets the note list dimensions
func (n *NoteList) SetSize(width, height int) {
	n.width = width
	n.height = height
}

// SetFocused sets the focus state
func (n *NoteList) SetFocused(focused bool) {
	n.focused = focused
}

// SetFolderName sets the folder name for the header
func (n *NoteList) SetFolderName(name string) {
	n.folderName = name
}

// SetShowTodos sets whether to show todo indicators
func (n *NoteList) SetShowTodos(show bool) {
	n.showTodos = show
}

// IsFocused returns whether the note list is focused
func (n *NoteList) IsFocused() bool {
	return n.focused
}

// Cursor returns the current cursor position
func (n *NoteList) Cursor() int {
	return n.cursor
}

// SelectedNote returns the currently selected note
func (n *NoteList) SelectedNote() *models.Note {
	if n.cursor < 0 || n.cursor >= len(n.notes) {
		return nil
	}
	return n.notes[n.cursor]
}

// ResetCursor resets the cursor to the top
func (n *NoteList) ResetCursor() {
	n.cursor = 0
}

// Update handles input
func (n *NoteList) Update(msg tea.Msg) (*NoteList, tea.Cmd) {
	if !n.focused {
		return n, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.DefaultKeyMap.Up):
			if n.cursor > 0 {
				n.cursor--
			}
		case key.Matches(msg, keys.DefaultKeyMap.Down):
			if n.cursor < len(n.notes)-1 {
				n.cursor++
			}
		}
	}

	return n, nil
}

// View renders the note list
func (n *NoteList) View() string {
	var b strings.Builder

	// Ensure minimum dimensions
	width := n.width
	if width < 30 {
		width = 30
	}

	// Account for border (2) and padding
	contentHeight := n.height - 2
	if contentHeight < 5 {
		contentHeight = 5
	}

	// Title
	title := n.folderName
	if title == "" {
		title = "All Notes"
	}
	b.WriteString(styles.NoteListTitleStyle.Render(fmt.Sprintf("📝 %s (%d)", title, len(n.notes))))
	b.WriteString("\n")

	sepWidth := width - 4
	if sepWidth < 10 {
		sepWidth = 10
	}
	b.WriteString(strings.Repeat("─", sepWidth))
	b.WriteString("\n")

	if len(n.notes) == 0 {
		b.WriteString(styles.TextMuted.Render("No notes yet. Press 'n' to create one."))
	} else {
		// Calculate visible range (subtract 3 for title, separator, padding)
		visibleHeight := contentHeight - 3
		if visibleHeight < 1 {
			visibleHeight = 1
		}

		startIdx := 0
		if n.cursor >= visibleHeight {
			startIdx = n.cursor - visibleHeight + 1
		}
		endIdx := startIdx + visibleHeight
		if endIdx > len(n.notes) {
			endIdx = len(n.notes)
		}

		// Render notes
		for i := startIdx; i < endIdx; i++ {
			note := n.notes[i]
			line := n.renderNote(note, i == n.cursor)
			b.WriteString(line)
			if i < endIdx-1 {
				b.WriteString("\n")
			}
		}
	}

	style := styles.NoteListStyle.Width(width - 4).Height(contentHeight)
	if n.focused {
		style = style.BorderForeground(styles.Primary)
	}
	return style.Render(b.String())
}

func (n *NoteList) renderNote(note *models.Note, selected bool) string {
	var parts []string

	// Todo checkbox
	if note.IsTodo {
		parts = append(parts, styles.RenderTodoStatus(note.IsDone))
	}

	// Priority
	if note.Priority > 0 {
		parts = append(parts, styles.RenderPriority(note.Priority))
	}

	// Star
	if note.Starred {
		parts = append(parts, styles.RenderStar(true))
	}

	// Title - calculate available space for title
	title := note.Title
	maxTitleLen := n.width - 25 // Reserve space for icons and date
	if maxTitleLen < 10 {
		maxTitleLen = 10
	}
	if len(title) > maxTitleLen {
		title = title[:maxTitleLen-3] + "..."
	}
	parts = append(parts, title)

	// Date
	date := note.UpdatedAt.Format("Jan 02")
	parts = append(parts, styles.NoteDateStyle.Render(date))

	text := strings.Join(parts, " ")

	// Calculate render width, ensuring it's positive
	renderWidth := n.width - 4
	if renderWidth < 20 {
		renderWidth = 20
	}

	if selected {
		return styles.NoteItemSelectedStyle.Width(renderWidth).Render(text)
	}

	if note.IsTodo && note.IsDone {
		return styles.TodoDoneStyle.Width(renderWidth).Render(text)
	}

	return styles.NoteItemStyle.Width(renderWidth).Render(text)
}

// Width returns the note list width
func (n *NoteList) Width() int {
	return n.width
}

// Height returns the note list height
func (n *NoteList) Height() int {
	return n.height
}
