package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette
var (
	// Base colors
	Primary   = lipgloss.Color("#7C3AED") // Purple
	Secondary = lipgloss.Color("#06B6D4") // Cyan
	Success   = lipgloss.Color("#10B981") // Green
	Warning   = lipgloss.Color("#F59E0B") // Amber
	Danger    = lipgloss.Color("#EF4444") // Red

	// Priority colors
	PriorityHigh   = lipgloss.Color("#EF4444") // Red
	PriorityMedium = lipgloss.Color("#F59E0B") // Amber
	PriorityLow    = lipgloss.Color("#10B981") // Green

	// UI element colors
	Background    = lipgloss.Color("#1E1E2E") // Dark
	Surface       = lipgloss.Color("#313244") // Slightly lighter
	Border        = lipgloss.Color("#45475A") // Border
	TextPrimary   = lipgloss.Color("#CDD6F4") // Light text
	TextSecondary = lipgloss.Color("#A6ADC8") // Dimmed text
	TextMutedC    = lipgloss.Color("#6C7086") // Very dimmed
)

// Base styles
var (
	// Text styles
	TextPrimaryStyle = lipgloss.NewStyle().Foreground(TextPrimary)
	TextMuted        = lipgloss.NewStyle().Foreground(TextMutedC)
	SuccessStyle     = lipgloss.NewStyle().Foreground(Success)
	ErrorStyle       = lipgloss.NewStyle().Foreground(Danger)

	// Header
	HeaderStyle = lipgloss.NewStyle().
			Background(Surface).
			Foreground(TextPrimary).
			Padding(0, 1)

	TitleStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true)

	DateStyle = lipgloss.NewStyle().
			Foreground(TextSecondary)

	// Sidebar styles
	SidebarStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Border).
			Padding(0, 1)

	SidebarTitleStyle = lipgloss.NewStyle().
				Foreground(Primary).
				Bold(true)

	FolderStyle = lipgloss.NewStyle().
			Foreground(TextPrimary)

	FolderSelectedStyle = lipgloss.NewStyle().
				Background(Primary).
				Foreground(lipgloss.Color("#FFFFFF")).
				Bold(true)

	FolderCountStyle = lipgloss.NewStyle().
				Foreground(TextMutedC)

	// Note list styles
	NoteListStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Border).
			Padding(0, 1)

	NoteListTitleStyle = lipgloss.NewStyle().
				Foreground(Primary).
				Bold(true)

	NoteItemStyle = lipgloss.NewStyle().
			Foreground(TextPrimary)

	NoteItemSelectedStyle = lipgloss.NewStyle().
				Background(Primary).
				Foreground(lipgloss.Color("#FFFFFF")).
				Bold(true)

	NoteDateStyle = lipgloss.NewStyle().
			Foreground(TextMutedC)

	TodoDoneStyle = lipgloss.NewStyle().
			Foreground(TextMutedC) // Removed Strikethrough - renders raw ANSI in some terminals

	// Preview styles
	PreviewStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Border).
			Padding(1)

	PreviewTitleStyle = lipgloss.NewStyle().
				Foreground(Primary).
				Bold(true)

	PreviewMetaStyle = lipgloss.NewStyle().
				Foreground(TextSecondary)

	PreviewContentStyle = lipgloss.NewStyle().
				Foreground(TextPrimary)

	// Status bar styles
	StatusBarStyle = lipgloss.NewStyle().
			Background(Surface).
			Foreground(TextSecondary).
			Padding(0, 1)

	StatusKeyStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true)

	StatusDescStyle = lipgloss.NewStyle().
			Foreground(TextSecondary)

	// Search bar styles
	SearchBarStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(0, 1)

	SearchIconStyle = lipgloss.NewStyle().
			Foreground(Primary)

	// Dialog styles
	DialogStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(Primary).
			Background(Surface).
			Padding(1, 2)

	DialogTitleStyle = lipgloss.NewStyle().
				Foreground(Primary).
				Bold(true)

	// Input styles
	InputStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(Border).
			Padding(0, 1)

	InputFocusedStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(Primary).
				Padding(0, 1)

	// Button styles
	ButtonStyle = lipgloss.NewStyle().
			Background(Surface).
			Foreground(TextPrimary).
			Padding(0, 2).
			Margin(0, 1)

	ButtonFocusedStyle = lipgloss.NewStyle().
				Background(Primary).
				Foreground(lipgloss.Color("#FFFFFF")).
				Padding(0, 2).
				Margin(0, 1).
				Bold(true)

	// Help styles
	HelpStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Background(Background).
			Padding(1, 2)

	HelpTitleStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true)

	HelpKeyStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Width(12).
			Bold(true)

	HelpDescStyle = lipgloss.NewStyle().
			Foreground(TextSecondary)
)

// Helper functions

// RenderPriority renders a priority indicator
func RenderPriority(priority int) string {
	switch priority {
	case 3:
		return lipgloss.NewStyle().Foreground(PriorityHigh).Render("●")
	case 2:
		return lipgloss.NewStyle().Foreground(PriorityMedium).Render("●")
	case 1:
		return lipgloss.NewStyle().Foreground(PriorityLow).Render("●")
	default:
		return ""
	}
}

// RenderStar renders a star indicator
func RenderStar(starred bool) string {
	if starred {
		return lipgloss.NewStyle().Foreground(Warning).Render("★")
	}
	return ""
}

// RenderTodoStatus renders a todo status indicator
func RenderTodoStatus(done bool) string {
	if done {
		return lipgloss.NewStyle().Foreground(Success).Render("☑")
	}
	return lipgloss.NewStyle().Foreground(TextSecondary).Render("☐")
}
