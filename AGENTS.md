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
│   │   └── interfaces.go # Repository interfaces for DI
│   ├── service/         # Business logic
│   │   └── interfaces.go # Service interfaces for DI
│   └── tui/             # TUI layer
│       ├── app.go       # Main TUI model
│       ├── commands/    # tea.Cmd builders (BubbleTea best practice #4)
│       ├── components/  # Reusable UI components
│       ├── constants/   # Magic numbers as named constants
│       ├── keys/        # Keybindings
│       ├── messages/    # Custom tea.Msg types (BubbleTea best practice #1)
│       └── styles/      # Lipgloss styles
└── docs/                # Documentation
```

---

## ⭐ Clean Code Principles (QUAN TRỌNG)

> **Nguyên tắc số 1**: Code phải đọc như văn xuôi. Nếu cần comment để giải thích, code chưa đủ rõ ràng.

### 1. Early Return - Thoát sớm, giảm lồng ghép

```go
// ❌ BAD: Nested nightmare
func (s *NoteService) Create(ctx context.Context, note *Note) error {
    if note != nil {
        if note.Title != "" {
            if s.repo != nil {
                err := s.repo.Insert(ctx, note)
                if err == nil {
                    return nil
                } else {
                    return err
                }
            } else {
                return errors.New("repo is nil")
            }
        } else {
            return errors.New("title is empty")
        }
    } else {
        return errors.New("note is nil")
    }
}

// ✅ GOOD: Early returns, flat structure
func (s *NoteService) Create(ctx context.Context, note *Note) error {
    if note == nil {
        return errors.New("note is nil")
    }
    if note.Title == "" {
        return errors.New("title is empty")
    }
    if s.repo == nil {
        return errors.New("repo is nil")
    }

    return s.repo.Insert(ctx, note)
}
```

### 2. Guard Clauses - Kiểm tra điều kiện đầu tiên

```go
// ❌ BAD: Main logic wrapped in conditions
func (a *App) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    if a.ready {
        if !a.showDialog {
            if a.currentPanel == PanelNoteList {
                // Main logic here, deeply nested
            }
        }
    }
    return a, nil
}

// ✅ GOOD: Guards at top, main logic at bottom
func (a *App) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    if !a.ready {
        return a, nil
    }
    if a.showDialog {
        return a.handleDialogInput(msg)
    }
    if a.currentPanel != PanelNoteList {
        return a.handleOtherPanel(msg)
    }

    // Main logic here, no nesting
    return a.handleNoteListInput(msg)
}
```

### 3. Single Responsibility - Một hàm làm một việc

```go
// ❌ BAD: Function does too many things
func (a *App) editNote(note *Note) (tea.Model, tea.Cmd) {
    // Create temp file
    tmpFile, _ := os.CreateTemp("", "*.md")
    content := fmt.Sprintf("# %s\n\n%s", note.Title, note.Content)
    tmpFile.WriteString(content)
    tmpFile.Close()

    // Run editor
    cmd := exec.Command("vim", tmpFile.Name())
    // ... 20 more lines

    // Parse result
    data, _ := os.ReadFile(tmpFile.Name())
    // ... 15 more lines

    // Update database
    // ... 10 more lines
}

// ✅ GOOD: Each function does ONE thing
func (a *App) editNote(note *Note) (tea.Model, tea.Cmd) {
    return a, commands.OpenEditor(a.editorService, note)
}

// In commands package:
func OpenEditor(svc EditorService, note *Note) tea.Cmd {
    tmpFile, cmd := svc.PrepareEdit(note)
    return tea.ExecProcess(cmd, handleEditorResult(tmpFile, note.ID))
}

// In editor_service.go:
func (s *EditorService) PrepareEdit(note *Note) (string, *exec.Cmd) { ... }
func (s *EditorService) ReadResult(tmpFile string) (string, string) { ... }
```

### 4. Meaningful Names - Tên có nghĩa

```go
// ❌ BAD: Cryptic names
func (a *App) hn(m tea.KeyMsg) (tea.Model, tea.Cmd) {
    n := a.nl.Sel()
    if n == nil {
        return a, nil
    }
    return a.e(n)
}

// ✅ GOOD: Self-documenting names
func (a *App) handleNoteSelection(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    selectedNote := a.noteList.SelectedNote()
    if selectedNote == nil {
        return a, nil
    }
    return a.editNote(selectedNote)
}

// Naming conventions:
// - Handlers: handleXxx()
// - Actions: xxxNote(), xxxFolder()
// - Getters: SelectedNote(), CurrentFolder()
// - Predicates: IsReady(), HasNotes(), CanEdit()
```

### 5. Extract Methods - Tách hàm con

```go
// ❌ BAD: Long method
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if key.Matches(msg, keys.Quit) { return a, tea.Quit }
        if key.Matches(msg, keys.Help) { a.showHelp = true; return a, nil }
        if key.Matches(msg, keys.Search) { a.searchMode = true; return a, nil }
        if key.Matches(msg, keys.NewNote) { return a.showNewNoteDialog() }
        if key.Matches(msg, keys.Tab) { a.switchPanel(1); return a, nil }
        // ... 50 more lines
    }
}

// ✅ GOOD: Extracted methods
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        return a.handleKeyPress(msg)
    }
    return a, nil
}

func (a *App) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    if cmd := a.handleGlobalKeys(msg); cmd != nil {
        return a, cmd
    }
    return a.handlePanelKeys(msg)
}

func (a *App) handleGlobalKeys(msg tea.KeyMsg) tea.Cmd {
    switch {
    case key.Matches(msg, keys.Quit):
        return tea.Quit
    case key.Matches(msg, keys.Help):
        a.showHelp = true
    case key.Matches(msg, keys.Search):
        a.searchMode = true
    }
    return nil
}
```

### 6. Avoid Boolean Parameters - Tránh tham số boolean

```go
// ❌ BAD: What does 'true' mean?
note := createNote("Title", true, false)
service.GetNotes(ctx, true, false, true)

// ✅ GOOD: Use options or separate methods
note := createTodoNote("Title")
note := createRegularNote("Title")

service.GetAllNotes(ctx)
service.GetCompletedTodos(ctx)
service.GetStarredNotes(ctx)

// Or use options struct:
service.GetNotes(ctx, NoteFilter{
    IncludeCompleted: true,
    OnlyStarred:      false,
})
```

### 7. Compose, Don't Inherit - Kết hợp, không kế thừa

```go
// ❌ BAD: Trying to do inheritance in Go
type BasePanel struct { ... }
type NoteListPanel struct { BasePanel }  // Embedding for inheritance

// ✅ GOOD: Composition with interfaces
type Panel interface {
    Update(tea.KeyMsg) (Panel, tea.Cmd)
    View() string
    SetFocused(bool)
}

type App struct {
    panels []Panel  // Compose panels
}

func (a *App) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    panel, cmd := a.currentPanel().Update(msg)
    a.setCurrentPanel(panel)
    return a, cmd
}
```

### 8. Table-Driven Logic - Dùng bảng thay vì switch dài

```go
// ❌ BAD: Long switch
func (a *App) handleGlobalKeys(msg tea.KeyMsg) tea.Cmd {
    switch {
    case key.Matches(msg, keys.Quit):
        return tea.Quit
    case key.Matches(msg, keys.Help):
        a.showHelp = true
        return nil
    case key.Matches(msg, keys.Search):
        a.searchMode = true
        return nil
    case key.Matches(msg, keys.NewNote):
        return a.showNewNoteDialog
    // ... 20 more cases
    }
}

// ✅ GOOD: Table-driven
var globalKeyHandlers = map[key.Binding]func(*App) tea.Cmd{
    keys.Quit:    func(a *App) tea.Cmd { return tea.Quit },
    keys.Help:    func(a *App) tea.Cmd { a.showHelp = true; return nil },
    keys.Search:  func(a *App) tea.Cmd { a.searchMode = true; return nil },
    keys.NewNote: func(a *App) tea.Cmd { return a.showNewNoteDialog },
}

func (a *App) handleGlobalKeys(msg tea.KeyMsg) tea.Cmd {
    for binding, handler := range globalKeyHandlers {
        if key.Matches(msg, binding) {
            return handler(a)
        }
    }
    return nil
}
```

### 9. Immutable Data - Tránh mutate dữ liệu

```go
// ❌ BAD: Mutating shared state
func (a *App) updateNote(note *Note) {
    note.Title = "New Title"  // Mutates the original
    a.noteService.Update(note)
}

// ✅ GOOD: Create new copy
func (a *App) updateNote(note *Note) {
    updated := *note  // Copy
    updated.Title = "New Title"
    a.noteService.Update(&updated)
}

// Or use functional approach:
func (n Note) WithTitle(title string) Note {
    n.Title = title
    return n
}
```

### 10. Constants Over Magic Numbers

```go
// ❌ BAD: Magic numbers
if len(title) > 100 { ... }
time.Sleep(3 * time.Second)
width := a.width * 25 / 100

// ✅ GOOD: Named constants
const (
    MaxTitleLength     = 100
    StatusMessageDelay = 3 * time.Second
    SidebarWidthRatio  = 0.25
)

if len(title) > MaxTitleLength { ... }
time.Sleep(StatusMessageDelay)
width := int(float64(a.width) * SidebarWidthRatio)
```

### 11. 🚨 Comments - Chỉ khi THỰC SỰ cần thiết

> **Quy tắc vàng**: Code tốt tự giải thích. Comment là thừa nhận code chưa đủ rõ.

```go
// ❌ BAD: Comment giải thích code đang làm gì (obvious)
// Check if note is nil
if note == nil {
    return nil
}

// Loop through all notes
for _, note := range notes {
    // Print the note title
    fmt.Println(note.Title)
}

// Create a new note service
noteService := NewNoteService(repo)

// ❌ BAD: Comment dư thừa
i++ // Increment i by 1
return nil // Return nil

// ❌ BAD: Comment outdated (LIE!)
// Get all active users  <-- Code actually gets notes!
notes, _ := repo.GetAll(ctx)
```

```go
// ✅ GOOD: Không cần comment - code tự giải thích
if note == nil {
    return nil
}

for _, note := range notes {
    fmt.Println(note.Title)
}

noteService := NewNoteService(repo)
```

**Khi nào ĐƯỢC comment:**

```go
// ✅ GOOD: WHY, not WHAT - Giải thích lý do, không phải hành động
// We use a 100ms delay because the terminal needs time to
// restore after the editor process exits
time.Sleep(100 * time.Millisecond)

// ✅ GOOD: Edge case / Bug workaround
// SQLite doesn't support concurrent writes, so we serialize
// all write operations through a single channel
writeChan <- writeRequest

// ✅ GOOD: API documentation (exported functions)
// NewNoteService creates a new note service with the given repository.
// It panics if repo is nil.
func NewNoteService(repo NoteRepository) *NoteService { ... }

// ✅ GOOD: Complex algorithm explanation
// Use binary search because notes are sorted by date.
// Linear search would be O(n) but this is O(log n).
idx := sort.Search(len(notes), func(i int) bool {
    return notes[i].CreatedAt.After(targetDate)
})

// ✅ GOOD: TODO with context
// TODO(quang): Refactor this when we add folder nesting support
// See: https://github.com/project/issues/123

// ✅ GOOD: Warning about non-obvious behavior
// WARNING: This function modifies the input slice in-place
func sortNotes(notes []*Note) { ... }
```

**Thay comment bằng code tốt hơn:**

```go
// ❌ BAD: Comment explaining complex condition
// Check if note is a todo that is not completed and has high priority
if note.IsTodo && !note.Completed && note.Priority > 2 {
    ...
}

// ✅ GOOD: Extract to meaningful function
if note.IsUrgentTodo() {
    ...
}

func (n *Note) IsUrgentTodo() bool {
    return n.IsTodo && !n.Completed && n.Priority > 2
}
```

```go
// ❌ BAD: Comment explaining magic number
if retryCount > 3 { // Max retries is 3
    return err
}

// ✅ GOOD: Use constant
const MaxRetries = 3
if retryCount > MaxRetries {
    return err
}
```

**Comment checklist:**

- [ ] Có thể rename variable/function để không cần comment?
- [ ] Có thể extract method với tên rõ nghĩa?
- [ ] Comment giải thích WHY hay WHAT? (chỉ WHY mới cần)
- [ ] Comment có thể bị outdated không? (nguy hiểm!)
- [ ] Đây có phải exported API cần document?

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
