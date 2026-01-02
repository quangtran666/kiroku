package keys

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines all keybindings for the application
type KeyMap struct {
	// Navigation
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Tab    key.Binding
	Enter  key.Binding
	Escape key.Binding

	// Actions
	NewNote       key.Binding
	NewTodo       key.Binding
	NewFolder     key.Binding
	Edit          key.Binding
	Delete        key.Binding
	Search        key.Binding
	ToggleStar    key.Binding
	ToggleDone    key.Binding
	MoveNote      key.Binding
	CyclePriority key.Binding

	// Views
	Help    key.Binding
	Preview key.Binding
	Quit    key.Binding
	Refresh key.Binding
}

// DefaultKeyMap returns the default key bindings
var DefaultKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "right"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch panel"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select/confirm"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back/cancel"),
	),

	// Actions
	NewNote: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new note"),
	),
	NewTodo: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "new todo"),
	),
	NewFolder: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "new folder"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d", "delete"),
		key.WithHelp("d", "delete"),
	),
	Search: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "search"),
	),
	ToggleStar: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "toggle star"),
	),
	ToggleDone: key.NewBinding(
		key.WithKeys("x", " "),
		key.WithHelp("x/space", "toggle done"),
	),
	MoveNote: key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m", "move to folder"),
	),
	CyclePriority: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "cycle priority"),
	),

	// Views
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Preview: key.NewBinding(
		key.WithKeys("v"),
		key.WithHelp("v", "toggle preview"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r", "ctrl+r"),
		key.WithHelp("r", "refresh"),
	),
}

// ShortHelp returns a short help string
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns the full help keybindings
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Tab, k.Enter, k.Escape},
		{k.NewNote, k.NewTodo, k.NewFolder},
		{k.Edit, k.Delete, k.Search},
		{k.ToggleStar, k.ToggleDone, k.CyclePriority},
		{k.Help, k.Preview, k.Quit, k.Refresh},
	}
}
