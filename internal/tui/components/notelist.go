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

	// Title
	title := n.folderName
	if title == "" {
		title = "All Notes"
	}
	b.WriteString(styles.NoteListTitleStyle.Render(fmt.Sprintf("📝 %s (%d)", title, len(n.notes))))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", n.width-2))
	b.WriteString("\n")

	if len(n.notes) == 0 {
		b.WriteString(styles.TextMuted.Render("No notes yet. Press 'n' to create one."))
	} else {
		// Calculate visible range
		visibleHeight := n.height - 4
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

		// Pad remaining height
		currentLines := endIdx - startIdx + 2
		for i := currentLines; i < n.height; i++ {
			b.WriteString("\n")
		}
	}

	style := styles.NoteListStyle.Width(n.width).Height(n.height)
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

	// Title
	title := note.Title
	if len(title) > n.width-20 {
		title = title[:n.width-23] + "..."
	}
	parts = append(parts, title)

	// Date
	date := note.UpdatedAt.Format("Jan 02")
	parts = append(parts, styles.NoteDateStyle.Render(date))

	text := strings.Join(parts, " ")

	if selected {
		return styles.NoteItemSelectedStyle.Width(n.width - 4).Render(text)
	}

	if note.IsTodo && note.IsDone {
		return styles.TodoDoneStyle.Width(n.width - 4).Render(text)
	}

	return styles.NoteItemStyle.Width(n.width - 4).Render(text)
}

// Width returns the note list width
func (n *NoteList) Width() int {
	return n.width
}

// Height returns the note list height
func (n *NoteList) Height() int {
	return n.height
}
