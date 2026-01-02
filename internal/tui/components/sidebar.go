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

// Sidebar represents the folder sidebar component
type Sidebar struct {
	folders     []*models.Folder
	flatList    []sidebarItem
	cursor      int
	height      int
	width       int
	focused     bool
	showAll     bool
	showTodos   bool
	showStarred bool
}

type sidebarItem struct {
	folder    *models.Folder
	isSpecial bool
	special   string // "all", "todos", "starred"
	level     int
}

// NewSidebar creates a new sidebar component
func NewSidebar() *Sidebar {
	return &Sidebar{
		folders:     make([]*models.Folder, 0),
		flatList:    make([]sidebarItem, 0),
		showAll:     true,
		showTodos:   true,
		showStarred: true,
	}
}

// SetFolders sets the folders to display
func (s *Sidebar) SetFolders(folders []*models.Folder) {
	s.folders = folders
	s.buildFlatList()
}

// SetSize sets the sidebar dimensions
func (s *Sidebar) SetSize(width, height int) {
	s.width = width
	s.height = height
}

// SetFocused sets the focus state
func (s *Sidebar) SetFocused(focused bool) {
	s.focused = focused
}

// IsFocused returns whether the sidebar is focused
func (s *Sidebar) IsFocused() bool {
	return s.focused
}

// Cursor returns the current cursor position
func (s *Sidebar) Cursor() int {
	return s.cursor
}

// SelectedFolder returns the currently selected folder
func (s *Sidebar) SelectedFolder() *models.Folder {
	if s.cursor < 0 || s.cursor >= len(s.flatList) {
		return nil
	}
	item := s.flatList[s.cursor]
	if item.isSpecial {
		return nil
	}
	return item.folder
}

// SelectedSpecial returns the selected special item ("all", "todos", "starred")
func (s *Sidebar) SelectedSpecial() string {
	if s.cursor < 0 || s.cursor >= len(s.flatList) {
		return ""
	}
	item := s.flatList[s.cursor]
	if item.isSpecial {
		return item.special
	}
	return ""
}

// buildFlatList builds a flat list of items for display
func (s *Sidebar) buildFlatList() {
	s.flatList = make([]sidebarItem, 0)

	// Add special items
	if s.showAll {
		s.flatList = append(s.flatList, sidebarItem{isSpecial: true, special: "all"})
	}

	// Add folders
	s.addFoldersToList(s.folders, 0)

	// Add quick access
	if s.showTodos {
		s.flatList = append(s.flatList, sidebarItem{isSpecial: true, special: "todos"})
	}
	if s.showStarred {
		s.flatList = append(s.flatList, sidebarItem{isSpecial: true, special: "starred"})
	}
}

func (s *Sidebar) addFoldersToList(folders []*models.Folder, level int) {
	for _, folder := range folders {
		s.flatList = append(s.flatList, sidebarItem{
			folder: folder,
			level:  level,
		})
		if folder.Expanded && len(folder.Children) > 0 {
			s.addFoldersToList(folder.Children, level+1)
		}
	}
}

// Update handles input
func (s *Sidebar) Update(msg tea.Msg) (*Sidebar, tea.Cmd) {
	if !s.focused {
		return s, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.DefaultKeyMap.Up):
			if s.cursor > 0 {
				s.cursor--
			}
		case key.Matches(msg, keys.DefaultKeyMap.Down):
			if s.cursor < len(s.flatList)-1 {
				s.cursor++
			}
		case key.Matches(msg, keys.DefaultKeyMap.Right), key.Matches(msg, keys.DefaultKeyMap.Enter):
			// Expand folder
			if s.cursor >= 0 && s.cursor < len(s.flatList) {
				item := s.flatList[s.cursor]
				if !item.isSpecial && item.folder != nil && len(item.folder.Children) > 0 {
					item.folder.Expanded = true
					s.buildFlatList()
				}
			}
		case key.Matches(msg, keys.DefaultKeyMap.Left):
			// Collapse folder
			if s.cursor >= 0 && s.cursor < len(s.flatList) {
				item := s.flatList[s.cursor]
				if !item.isSpecial && item.folder != nil {
					if item.folder.Expanded {
						item.folder.Expanded = false
						s.buildFlatList()
					}
				}
			}
		}
	}

	return s, nil
}

// View renders the sidebar
func (s *Sidebar) View() string {
	var b strings.Builder

	// Title
	title := styles.SidebarTitleStyle.Render("📁 FOLDERS")
	b.WriteString(title)
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", s.width-2))
	b.WriteString("\n")

	// Calculate visible range
	visibleHeight := s.height - 4 // Account for title and separators
	startIdx := 0
	if s.cursor >= visibleHeight {
		startIdx = s.cursor - visibleHeight + 1
	}
	endIdx := startIdx + visibleHeight
	if endIdx > len(s.flatList) {
		endIdx = len(s.flatList)
	}

	// Render items
	for i := startIdx; i < endIdx; i++ {
		item := s.flatList[i]
		line := s.renderItem(item, i == s.cursor)
		b.WriteString(line)
		if i < endIdx-1 {
			b.WriteString("\n")
		}
	}

	// Pad remaining height
	currentLines := endIdx - startIdx + 2 // +2 for title and separator
	for i := currentLines; i < s.height; i++ {
		b.WriteString("\n")
	}

	style := styles.SidebarStyle.Width(s.width).Height(s.height)
	if s.focused {
		style = style.BorderForeground(styles.Primary)
	}
	return style.Render(b.String())
}

func (s *Sidebar) renderItem(item sidebarItem, selected bool) string {
	var icon, name, count string
	indent := strings.Repeat("  ", item.level)

	if item.isSpecial {
		switch item.special {
		case "all":
			icon = "📋"
			name = "All Notes"
		case "todos":
			icon = "☐"
			name = "Todos"
		case "starred":
			icon = "⭐"
			name = "Starred"
		}
	} else if item.folder != nil {
		icon = item.folder.Icon
		name = item.folder.Name
		// Show expand/collapse indicator
		if len(item.folder.Children) > 0 {
			if item.folder.Expanded {
				indent += "▾ "
			} else {
				indent += "▸ "
			}
		} else {
			indent += "  "
		}
		if item.folder.NoteCount > 0 {
			count = fmt.Sprintf(" (%d)", item.folder.NoteCount)
		}
	}

	text := fmt.Sprintf("%s%s %s%s", indent, icon, name, styles.FolderCountStyle.Render(count))

	if selected {
		return styles.FolderSelectedStyle.Width(s.width - 4).Render(text)
	}
	return styles.FolderStyle.Width(s.width - 4).Render(text)
}

// Width returns the sidebar width
func (s *Sidebar) Width() int {
	return s.width
}

// Height returns the sidebar height
func (s *Sidebar) Height() int {
	return s.height
}
