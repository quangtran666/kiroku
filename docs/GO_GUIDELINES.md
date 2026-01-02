# Go Coding Guidelines - Kiroku Project

> Best practices và conventions cho project này

---

## 📁 Project Layout

```
kiroku/
├── cmd/kiroku/main.go    # Entry point - chỉ wire dependencies
├── internal/             # Private code, không export ra ngoài
│   ├── domain/           # Models, interfaces (không dependencies)
│   ├── repository/       # Data access layer
│   ├── service/          # Business logic
│   └── tui/              # Presentation layer
└── pkg/                  # Public utilities (nếu cần)
```

---

## 🎯 Naming Conventions

```go
// ✅ Package names: lowercase, singular, ngắn gọn
package repository  // không phải "repositories"
package service     // không phải "services"

// ✅ Interface names: hành động + "er"
type Reader interface { Read() }
type NoteRepository interface { ... }

// ✅ Struct names: danh từ
type Note struct { ... }
type NoteService struct { ... }

// ✅ Function/Method: động từ hoặc câu hỏi
func CreateNote() {}     // hành động
func IsValid() bool {}   // câu hỏi → return bool
func GetByID() {}        // lấy data

// ✅ Variables: camelCase, mô tả rõ ràng
noteCount := 10          // không phải "n" hoặc "nc"
userID := 123            // ID viết hoa
httpClient := ...        // HTTP viết hoa

// ✅ Constants: PascalCase hoặc ALL_CAPS cho env
const MaxRetries = 3
const DefaultTimeout = 30 * time.Second
```

---

## 🏗️ Struct & Interface

```go
// ✅ Interface nhỏ, focused (1-3 methods)
type NoteReader interface {
    GetByID(ctx context.Context, id int64) (*Note, error)
    List(ctx context.Context, opts ListOptions) ([]Note, error)
}

type NoteWriter interface {
    Create(ctx context.Context, note *Note) error
    Update(ctx context.Context, note *Note) error
    Delete(ctx context.Context, id int64) error
}

// Compose interfaces
type NoteRepository interface {
    NoteReader
    NoteWriter
}

// ✅ Struct với constructor
type NoteService struct {
    repo   NoteRepository
    logger *slog.Logger
}

// Constructor function - return interface, not struct
func NewNoteService(repo NoteRepository, logger *slog.Logger) *NoteService {
    return &NoteService{
        repo:   repo,
        logger: logger,
    }
}
```

---

## ⚠️ Error Handling

```go
// ✅ Wrap errors với context
func (s *NoteService) GetByID(ctx context.Context, id int64) (*Note, error) {
    note, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("get note %d: %w", id, err)
    }
    return note, nil
}

// ✅ Define sentinel errors
var (
    ErrNotFound     = errors.New("not found")
    ErrInvalidInput = errors.New("invalid input")
)

// ✅ Check specific errors
if errors.Is(err, ErrNotFound) {
    // handle not found
}

// ✅ Custom error types khi cần thêm context
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}
```

---

## 🔄 Context Usage

```go
// ✅ Context là param đầu tiên
func (r *NoteRepo) GetByID(ctx context.Context, id int64) (*Note, error) {
    // Check context trước khi làm việc nặng
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
    }

    // Dùng context cho DB queries
    row := r.db.QueryRowContext(ctx, "SELECT ... WHERE id = ?", id)
    // ...
}

// ✅ Timeout cho operations
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

note, err := service.GetByID(ctx, 123)
```

---

## 🧪 Testing

```go
// ✅ Table-driven tests
func TestNoteService_Create(t *testing.T) {
    tests := []struct {
        name    string
        input   *Note
        wantErr bool
    }{
        {
            name:    "valid note",
            input:   &Note{Title: "Test"},
            wantErr: false,
        },
        {
            name:    "empty title",
            input:   &Note{Title: ""},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            svc := setupTestService(t)

            // Act
            err := svc.Create(context.Background(), tt.input)

            // Assert
            if tt.wantErr {
                require.Error(t, err)
            } else {
                require.NoError(t, err)
            }
        })
    }
}

// ✅ Test helpers
func setupTestService(t *testing.T) *NoteService {
    t.Helper()
    db := setupTestDB(t)
    repo := NewNoteRepository(db)
    return NewNoteService(repo, slog.Default())
}

func setupTestDB(t *testing.T) *sql.DB {
    t.Helper()
    db, err := sql.Open("sqlite", ":memory:")
    require.NoError(t, err)
    t.Cleanup(func() { db.Close() })
    return db
}
```

---

## 📝 Code Style Rules

### DO ✅

```go
// 1. Early return - giảm nesting
func Process(note *Note) error {
    if note == nil {
        return ErrInvalidInput
    }
    if note.Title == "" {
        return &ValidationError{Field: "title", Message: "required"}
    }
    // main logic here
    return nil
}

// 2. Named return values cho documentation
func (r *Repo) Stats() (total int, pending int, err error) {
    // ...
}

// 3. Functional options cho config phức tạp
type Option func(*Config)

func WithTimeout(d time.Duration) Option {
    return func(c *Config) { c.Timeout = d }
}

func NewClient(opts ...Option) *Client {
    cfg := defaultConfig()
    for _, opt := range opts {
        opt(&cfg)
    }
    return &Client{cfg: cfg}
}

// 4. Defer cho cleanup
func (r *Repo) Query() ([]Note, error) {
    rows, err := r.db.Query("SELECT ...")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    // ...
}
```

### DON'T ❌

```go
// 1. Không panic trong library code
func GetNote(id int) *Note {
    panic("not found") // ❌ Dùng error thay vì panic
}

// 2. Không ignore errors
result, _ := doSomething() // ❌ Luôn handle error

// 3. Không dùng init() trừ khi cần thiết
func init() { // ❌ Khó test, side effects
    globalDB = connect()
}

// 4. Không return concrete type khi có thể dùng interface
func NewRepo() *SQLiteRepo { } // ❌
func NewRepo() Repository { }  // ✅

// 5. Không nested quá 3 levels
if x {
    if y {
        if z { // ❌ Quá deep
        }
    }
}
```

---

## 📦 Dependencies

```go
// go.mod - Pin versions
module github.com/user/kiroku

go 1.21

require (
    github.com/charmbracelet/bubbletea v0.25.0
    github.com/charmbracelet/lipgloss v0.9.1
    github.com/spf13/cobra v1.8.0
    modernc.org/sqlite v1.28.0
)
```

### Sử dụng

| Package              | Purpose                         |
| -------------------- | ------------------------------- |
| `modernc.org/sqlite` | SQLite driver (pure Go, no CGO) |
| `bubbletea`          | TUI framework                   |
| `lipgloss`           | TUI styling                     |
| `bubbles`            | TUI components                  |
| `cobra`              | CLI framework                   |
| `testify`            | Testing assertions              |

---

## 🔧 Tools

```bash
# Format code
go fmt ./...

# Lint
golangci-lint run

# Test với coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Build
go build -o kiroku ./cmd/kiroku
```

### golangci-lint config (.golangci.yml)

```yaml
linters:
  enable:
    - errcheck # Check error handling
    - govet # Go vet
    - staticcheck # Static analysis
    - unused # Unused code
    - gosimple # Simplify code
    - ineffassign # Unused assignments
    - gofmt # Format check

linters-settings:
  errcheck:
    check-blank: true
```

---

## 📚 Quick Reference

| Concept     | Pattern                                          |
| ----------- | ------------------------------------------------ |
| Constructor | `func New...(deps) *Type`                        |
| Error       | `return fmt.Errorf("context: %w", err)`          |
| Context     | Param đầu tiên: `func(ctx context.Context, ...)` |
| Interface   | Nhỏ, 1-3 methods                                 |
| Test        | Table-driven với `t.Run()`                       |
| Cleanup     | `defer resource.Close()`                         |
| Check nil   | Early return ở đầu function                      |
