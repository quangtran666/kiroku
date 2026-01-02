package components

import (
	"strings"

	"github.com/tranducquang/kiroku/internal/models"
	"github.com/tranducquang/kiroku/internal/tui/styles"
)

// Preview represents the note preview component
type Preview struct {
	note   *models.Note
	height int
	width  int
	scroll int
}

// NewPreview creates a new preview component
func NewPreview() *Preview {
	return &Preview{}
}

// SetNote sets the note to preview
func (p *Preview) SetNote(note *models.Note) {
	p.note = note
	p.scroll = 0
}

// SetSize sets the preview dimensions
func (p *Preview) SetSize(width, height int) {
	p.width = width
	p.height = height
}

// ScrollUp scrolls the preview up
func (p *Preview) ScrollUp() {
	if p.scroll > 0 {
		p.scroll--
	}
}

// ScrollDown scrolls the preview down
func (p *Preview) ScrollDown() {
	p.scroll++
}

// View renders the preview
func (p *Preview) View() string {
	var b strings.Builder

	if p.note == nil {
		b.WriteString(styles.TextMuted.Render("Select a note to preview"))
		return styles.PreviewStyle.Width(p.width).Height(p.height).Render(b.String())
	}

	// Title
	b.WriteString(styles.PreviewTitleStyle.Render(p.note.Title))
	b.WriteString("\n")

	// Meta info
	meta := []string{}
	if p.note.IsTodo {
		if p.note.IsDone {
			meta = append(meta, "✓ Done")
		} else {
			meta = append(meta, "☐ Todo")
		}
	}
	if p.note.Priority > 0 {
		priority := []string{"", "Low", "Medium", "High"}[p.note.Priority]
		meta = append(meta, "Priority: "+priority)
	}
	if p.note.DueDate != nil {
		meta = append(meta, "Due: "+p.note.DueDate.Format("Jan 02, 2006"))
	}
	meta = append(meta, "Updated: "+p.note.UpdatedAt.Format("Jan 02, 2006 15:04"))

	b.WriteString(styles.PreviewMetaStyle.Render(strings.Join(meta, " • ")))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", p.width-4))
	b.WriteString("\n")

	// Content
	content := p.note.Content
	lines := strings.Split(content, "\n")

	// Apply scroll
	if p.scroll < len(lines) {
		lines = lines[p.scroll:]
	}

	// Limit visible lines
	visibleLines := p.height - 6
	if len(lines) > visibleLines {
		lines = lines[:visibleLines]
	}

	b.WriteString(styles.PreviewContentStyle.Render(strings.Join(lines, "\n")))

	return styles.PreviewStyle.Width(p.width).Height(p.height).Render(b.String())
}

// Width returns the preview width
func (p *Preview) Width() int {
	return p.width
}

// Height returns the preview height
func (p *Preview) Height() int {
	return p.height
}
