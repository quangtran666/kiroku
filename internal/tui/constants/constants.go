// Package constants contains all magic numbers and constants for TUI layout.
// This follows Clean Code principle #10: Constants Over Magic Numbers.
package constants

import "time"

// Layout constants
const (
	// SidebarMinWidth is the minimum width for the sidebar panel.
	SidebarMinWidth = 20
	// SidebarMaxWidth is the maximum width for the sidebar panel.
	SidebarMaxWidth = 40
	// SidebarWidthPercent is the percentage of screen width for sidebar.
	SidebarWidthPercent = 25
	// StatusBarHeight is the height reserved for status bar.
	StatusBarHeight = 3
	// PreviewHeightRatio is the ratio of note list height for preview.
	PreviewHeightRatio = 0.5
)

// Timing constants
const (
	// StatusMessageDuration is how long status messages are displayed.
	StatusMessageDuration = 2 * time.Second
	// ErrorMessageDuration is how long error messages are displayed.
	ErrorMessageDuration = 3 * time.Second
)

// Priority constants for todos
const (
	// PriorityNone represents no priority.
	PriorityNone = 0
	// PriorityLow represents low priority.
	PriorityLow = 1
	// PriorityMedium represents medium priority.
	PriorityMedium = 2
	// PriorityHigh represents high priority.
	PriorityHigh = 3
	// PriorityMax is the maximum priority value (for cycling).
	PriorityMax = 4
)

// Dialog types
const (
	DialogTypeNewNote = "new_note"
	DialogTypeNewTodo = "new_todo"
	DialogTypeDelete  = "delete"
	DialogTypeConfirm = "confirm"
)

// Filter types for sidebar
const (
	FilterAll     = "all"
	FilterTodos   = "todos"
	FilterStarred = "starred"
)
