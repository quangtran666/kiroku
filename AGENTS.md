# Kiroku Development Guidelines

> **For AI Assistants**: This file contains best practices and guidelines for working on the Kiroku codebase. Follow these strictly.

## Project Overview

Kiroku (記録) is a terminal-based note-taking application built with:

- **Go 1.21+** - Using modern Go features (generics, slog, etc.)
- **BubbleTea** - Elm-architecture TUI framework
- **Cobra** - CLI framework
- **SQLite** - Local database

## Code Organization

```
kiroku/
├── cmd/kiroku/          # Entry point only, minimal code
├── internal/
│   ├── app/             # Application wiring, dependency injection
│   ├── cli/             # Cobra commands
│   ├── config/          # Configuration loading
│   ├── database/        # Database connection, migrations
│   ├── logging/         # Logging setup
│   ├── models/          # Domain models (pure data structures)
│   ├── repository/      # Data access layer
│   ├── service/         # Business logic
│   └── tui/             # TUI layer
│       ├── app.go       # Main TUI model
│       ├── components/  # Reusable UI components
│       ├── keys/        # Keybindings
│       ├── messages/    # Custom tea.Msg types (NEW)
│       └── styles/      # Lipgloss styles
└── docs/                # Documentation
```

---

## Go Best Practices

### 1. Error Handling

```go
// ✅ GOOD: Wrap errors with context
if err := db.Query(ctx, query); err != nil {
    return fmt.Errorf("fetch notes: %w", err)
}

// ❌ BAD: Lose error context
if err := db.Query(ctx, query); err != nil {
    return err
}

// ❌ BAD: String formatting errors
if err := db.Query(ctx, query); err != nil {
    return fmt.Errorf("error: %s", err.Error())
}
```

### 2. Context Usage

```go
// ✅ GOOD: Accept context as first parameter
func (s *NoteService) GetByID(ctx context.Context, id int64) (*Note, error)

// ✅ GOOD: Create context with timeout for external operations
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()

// ❌ BAD: Use context.Background() in business logic
func (s *NoteService) GetByID(id int64) (*Note, error) {
    ctx := context.Background() // Don't do this
}
```

### 3. Dependency Injection

```go
// ✅ GOOD: Accept interfaces, return structs
type NoteRepository interface {
    GetByID(ctx context.Context, id int64) (*models.Note, error)
}

type NoteService struct {
    repo   NoteRepository  // Interface
    logger *zerolog.Logger
}

func NewNoteService(repo NoteRepository, logger *zerolog.Logger) *NoteService {
    return &NoteService{repo: repo, logger: logger}
}

// ❌ BAD: Concrete dependencies
type NoteService struct {
    repo *repository.NoteRepository  // Concrete type
}
```

### 4. Logging

```go
// ✅ GOOD: Structured logging with context
logging.Info().
    Str("note_id", noteID).
    Str("action", "create").
    Msg("Note created successfully")

// ✅ GOOD: Log errors with stack context
logging.Error().
    Err(err).
    Str("note_id", noteID).
    Msg("Failed to create note")

// ❌ BAD: Unstructured logging
log.Printf("Created note %s", noteID)

// ❌ BAD: Logging sensitive data
logging.Info().Str("password", password).Msg("User login")
```

---

## BubbleTea Best Practices

### 1. Message Types - Use Dedicated Package

```go
// ✅ GOOD: Separate messages package
// internal/tui/messages/messages.go
package messages

// DataLoadedMsg indicates data has been loaded
type DataLoadedMsg struct {
    Folders   []*models.Folder
    Notes     []*models.Note
    Templates []models.Template
}

// ErrorMsg wraps an error as a message
type ErrorMsg struct {
    Err     error
    Context string
}

func (e ErrorMsg) Error() string {
    return fmt.Sprintf("%s: %v", e.Context, e.Err)
}

// EditorFinishedMsg indicates editor process completed
type EditorFinishedMsg struct {
    TempFile string
    NoteID   int64
}
```

### 2. Model Structure - Single Responsibility

```go
// ✅ GOOD: Small, focused model
type App struct {
    // Services (injected)
    noteService   NoteService
    editorService EditorService

    // UI State
    state      AppState
    dimensions Dimensions

    // Components (composed)
    sidebar  *Sidebar
    noteList *NoteList

    // NO business logic here
}

// ✅ GOOD: Separate state struct
type AppState struct {
    CurrentView   ViewType
    CurrentPanel  Panel
    CurrentNote   *models.Note
    CurrentFolder *models.Folder
    IsEditing     bool
    EditTempFile  string
}

// ✅ GOOD: Separate dimensions
type Dimensions struct {
    Width  int
    Height int
    Ready  bool
}
```

### 3. Update Method - Keep Clean

```go
// ✅ GOOD: Delegate to specific handlers
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        return a.handleWindowResize(msg)
    case tea.KeyMsg:
        return a.handleKeyPress(msg)
    case messages.DataLoadedMsg:
        return a.handleDataLoaded(msg)
    case messages.EditorFinishedMsg:
        return a.handleEditorFinished(msg)
    case messages.ErrorMsg:
        return a.handleError(msg)
    }
    return a, nil
}

// ✅ GOOD: Separate handler methods
func (a *App) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    // Handle global keys first
    if cmd := a.handleGlobalKeys(msg); cmd != nil {
        return a, cmd
    }

    // Delegate to current panel
    return a.currentPanelHandler(msg)
}

// ❌ BAD: Giant switch statement
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if a.showHelp {
            // 50 lines of code...
        }
        if a.showDialog {
            // 50 more lines...
        }
        // 200 more lines of nested conditions...
    }
}
```

### 4. Commands - Use Command Builders

```go
// ✅ GOOD: Command builder pattern
// internal/tui/commands/commands.go
package commands

func LoadData(noteService NoteService) tea.Cmd {
    return func() tea.Msg {
        ctx := context.Background()
        notes, err := noteService.GetAll(ctx)
        if err != nil {
            return messages.ErrorMsg{Err: err, Context: "load notes"}
        }
        return messages.DataLoadedMsg{Notes: notes}
    }
}

func OpenEditor(editorService EditorService, note *models.Note) tea.Cmd {
    tmpFile, cmd, err := editorService.PrepareEdit(note.Title, note.Content)
    if err != nil {
        return func() tea.Msg {
            return messages.ErrorMsg{Err: err, Context: "prepare editor"}
        }
    }

    return tea.ExecProcess(cmd, func(err error) tea.Msg {
        if err != nil {
            return messages.ErrorMsg{Err: err, Context: "run editor"}
        }
        return messages.EditorFinishedMsg{TempFile: tmpFile, NoteID: note.ID}
    })
}

// Usage in app.go
func (a *App) editNote(note *models.Note) (tea.Model, tea.Cmd) {
    a.state.IsEditing = true
    a.state.CurrentNote = note
    return a, commands.OpenEditor(a.editorService, note)
}
```

### 5. External Processes

```go
// ✅ GOOD: Use tea.ExecProcess for external programs
func (a *App) editNote(note *models.Note) (tea.Model, tea.Cmd) {
    cmd := exec.Command("vim", tmpFile)
    return a, tea.ExecProcess(cmd, func(err error) tea.Msg {
        return messages.EditorFinishedMsg{NoteID: note.ID}
    })
}

// ❌ BAD: Direct exec.Command in TUI
func (a *App) editNote(note *models.Note) (tea.Model, tea.Cmd) {
    return a, func() tea.Msg {
        cmd := exec.Command("vim", tmpFile)
        cmd.Stdin = os.Stdin   // WRONG: TUI owns stdin
        cmd.Stdout = os.Stdout // WRONG: TUI owns stdout
        cmd.Run()              // BLOCKS and corrupts terminal
        return nil
    }
}
```

### 6. Components - Self-Contained

```go
// ✅ GOOD: Component manages own state
type NoteList struct {
    notes    []*models.Note
    cursor   int
    focused  bool
    width    int
    height   int
    delegate NoteListDelegate // For callbacks
}

// Component has its own Update
func (n *NoteList) Update(msg tea.KeyMsg) (*NoteList, tea.Cmd) {
    switch {
    case key.Matches(msg, keys.Up):
        n.cursor = max(0, n.cursor-1)
    case key.Matches(msg, keys.Down):
        n.cursor = min(len(n.notes)-1, n.cursor+1)
    }
    return n, nil
}

// Component has its own View
func (n *NoteList) View() string {
    // Render logic here
}

// ❌ BAD: Parent manages component state
func (a *App) handleNoteListInput(msg tea.KeyMsg) {
    if msg.String() == "j" {
        a.noteListCursor++ // Parent manages cursor
    }
}
```

---

## Logging Guidelines

### Log Levels

| Level   | When to Use                                  |
| ------- | -------------------------------------------- |
| `Debug` | Development info, key presses, state changes |
| `Info`  | Important events: note created, user actions |
| `Warn`  | Recoverable issues, deprecated usage         |
| `Error` | Failures that need attention                 |
| `Fatal` | Unrecoverable, app must exit                 |

### What to Log

```go
// ✅ DO log:
- User actions (create, delete, edit)
- State transitions
- External process start/end
- Errors with context
- Performance metrics for slow operations

// ❌ DON'T log:
- Every key press in production (debug only)
- Sensitive data (passwords, tokens)
- High-frequency events without sampling
- Successful operations that happen constantly
```

---

## File Naming Conventions

```
internal/tui/
├── app.go              # Main model
├── app_handlers.go     # Update handlers (split from app.go)
├── app_actions.go      # Action methods (create, delete, etc.)
├── components/
│   ├── sidebar.go
│   ├── note_list.go    # snake_case for multi-word
│   └── dialog.go
├── messages/
│   └── messages.go     # All tea.Msg types
├── commands/
│   └── commands.go     # All tea.Cmd builders
├── keys/
│   └── keys.go
└── styles/
    └── styles.go
```

---

## Testing Guidelines

```go
// ✅ GOOD: Table-driven tests
func TestNoteService_Create(t *testing.T) {
    tests := []struct {
        name    string
        note    *models.Note
        wantErr bool
    }{
        {
            name:    "valid note",
            note:    &models.Note{Title: "Test"},
            wantErr: false,
        },
        {
            name:    "empty title",
            note:    &models.Note{Title: ""},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test logic
        })
    }
}

// ✅ GOOD: Use interfaces for mocking
type MockNoteRepository struct {
    notes []*models.Note
}

func (m *MockNoteRepository) GetByID(ctx context.Context, id int64) (*models.Note, error) {
    // Mock implementation
}
```

---

## Common Patterns

### 1. Option Pattern for Configuration

```go
type Option func(*App)

func WithLogger(logger *zerolog.Logger) Option {
    return func(a *App) {
        a.logger = logger
    }
}

func NewApp(opts ...Option) *App {
    app := &App{/* defaults */}
    for _, opt := range opts {
        opt(app)
    }
    return app
}
```

### 2. Result Type for Operations

```go
type Result[T any] struct {
    Value T
    Err   error
}

func (r Result[T]) Unwrap() (T, error) {
    return r.Value, r.Err
}
```

---

## Refactoring Checklist

When touching TUI code, ensure:

- [ ] Messages are in `messages/` package
- [ ] Commands are in `commands/` package
- [ ] Update handlers are split into focused methods
- [ ] Components are self-contained
- [ ] No business logic in TUI layer
- [ ] Logging uses structured format
- [ ] External processes use `tea.ExecProcess`
- [ ] Context is passed through service calls
- [ ] Errors are wrapped with context

---

## References

- [BubbleTea Best Practices](https://github.com/charmbracelet/bubbletea/tree/master/tutorials)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
