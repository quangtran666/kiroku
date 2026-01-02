# 記録 Kiroku - Implementation Plan

> A lightweight, terminal-based note-taking app with TUI built in Go + SQLite

## 📋 Table of Contents

1. [Overview](#overview)
2. [Tech Stack](#tech-stack)
3. [Project Structure](#project-structure)
4. [Database Schema](#database-schema)
5. [Core Features](#core-features)
6. [Templates System](#templates-system)
7. [TUI Components](#tui-components)
8. [Implementation Phases](#implementation-phases)
9. [Testing Strategy](#testing-strategy)
10. [CLI Commands](#cli-commands)

---

## Overview

### Goals

- ⚡ Fast and lightweight
- 🔍 Quick full-text search
- 📁 Folder-based organization
- ☐ Todo management with priorities
- 📝 Templates for quick note creation
- 🖥️ Beautiful TUI interface
- ✏️ Vim/Neovim integration for editing

### Non-Goals (v1)

- Cloud sync
- Collaboration
- Mobile app
- Encryption

---

## Tech Stack

| Component      | Library                     | Version         |
| -------------- | --------------------------- | --------------- |
| Language       | Go                          | 1.21+           |
| Database       | SQLite (modernc.org/sqlite) | Pure Go, no CGO |
| TUI Framework  | Bubble Tea                  | Latest          |
| TUI Styling    | Lipgloss                    | Latest          |
| TUI Components | Bubbles                     | Latest          |
| CLI            | Cobra                       | v1.8+           |
| Config         | Viper                       | Latest          |
| Testing        | testify                     | Latest          |
| Mocking        | gomock                      | Latest          |

---

## Project Structure

```
kiroku/
├── cmd/
│   └── kiroku/
│       └── main.go                 # Entry point
│
├── internal/
│   ├── app/
│   │   └── app.go                  # Application orchestrator
│   │
│   ├── config/
│   │   ├── config.go               # Configuration struct
│   │   └── config_test.go
│   │
│   ├── database/
│   │   ├── database.go             # DB connection & migrations
│   │   ├── database_test.go
│   │   └── migrations/
│   │       └── 001_initial.sql
│   │
│   ├── models/
│   │   ├── note.go                 # Note model
│   │   ├── folder.go               # Folder model
│   │   ├── template.go             # Template model
│   │   └── todo.go                 # Todo model (extends Note)
│   │
│   ├── repository/
│   │   ├── note_repo.go            # Note CRUD operations
│   │   ├── note_repo_test.go
│   │   ├── folder_repo.go          # Folder CRUD
│   │   ├── folder_repo_test.go
│   │   ├── template_repo.go        # Template CRUD
│   │   ├── template_repo_test.go
│   │   └── search_repo.go          # Full-text search
│   │
│   ├── service/
│   │   ├── note_service.go         # Business logic for notes
│   │   ├── note_service_test.go
│   │   ├── folder_service.go
│   │   ├── template_service.go
│   │   └── editor_service.go       # Vim/Neovim integration
│   │
│   ├── tui/
│   │   ├── app.go                  # Main TUI application
│   │   ├── styles.go               # Lipgloss styles
│   │   ├── keys.go                 # Keybindings
│   │   ├── components/
│   │   │   ├── sidebar.go          # Folder sidebar
│   │   │   ├── notelist.go         # Note list panel
│   │   │   ├── preview.go          # Note preview panel
│   │   │   ├── searchbar.go        # Search input
│   │   │   ├── dialog.go           # Modal dialogs
│   │   │   ├── statusbar.go        # Bottom status bar
│   │   │   └── help.go             # Help overlay
│   │   └── views/
│   │       ├── main_view.go        # Main dashboard
│   │       ├── note_view.go        # Note detail view
│   │       ├── todo_view.go        # Todo list view
│   │       ├── search_view.go      # Search results
│   │       ├── template_view.go    # Template picker
│   │       └── new_note_view.go    # New note dialog
│   │
│   └── cli/
│       ├── root.go                 # Root command
│       ├── add.go                  # Quick add command
│       ├── list.go                 # List notes command
│       └── search.go               # Search command
│
├── testdata/
│   └── fixtures/                   # Test fixtures
│
├── docs/
│   ├── IMPLEMENTATION_PLAN.md      # This file
│   └── ARCHITECTURE.md             # Architecture details
│
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

---

## Database Schema

### ERD Diagram

```
┌─────────────┐       ┌─────────────┐       ┌─────────────┐
│   folders   │       │    notes    │       │  templates  │
├─────────────┤       ├─────────────┤       ├─────────────┤
│ id (PK)     │◄──┐   │ id (PK)     │       │ id (PK)     │
│ name        │   │   │ title       │       │ name        │
│ parent_id   │───┘   │ content     │       │ content     │
│ icon        │       │ folder_id   │───────│ type        │
│ created_at  │       │ is_todo     │       │ icon        │
│ updated_at  │       │ is_done     │       │ variables   │
└─────────────┘       │ priority    │       │ created_at  │
                      │ due_date    │       │ updated_at  │
                      │ tags        │       └─────────────┘
                      │ starred     │
                      │ template_id │───────────────┘
                      │ created_at  │
                      │ updated_at  │
                      └─────────────┘
                             │
                             ▼
                      ┌─────────────┐
                      │  notes_fts  │ (Virtual Table - FTS5)
                      ├─────────────┤
                      │ title       │
                      │ content     │
                      │ tags        │
                      └─────────────┘
```

### SQL Schema

```sql
-- migrations/001_initial.sql

-- Enable foreign keys
PRAGMA foreign_keys = ON;

-- Folders table
CREATE TABLE IF NOT EXISTS folders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    parent_id INTEGER REFERENCES folders(id) ON DELETE CASCADE,
    icon TEXT DEFAULT '📁',
    position INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create index for parent lookup
CREATE INDEX idx_folders_parent ON folders(parent_id);

-- Templates table
CREATE TABLE IF NOT EXISTS templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    content TEXT NOT NULL,
    type TEXT NOT NULL CHECK(type IN ('note', 'todo')),
    icon TEXT DEFAULT '📄',
    variables TEXT, -- JSON: [{"name": "project", "default": ""}]
    is_default BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Notes table
CREATE TABLE IF NOT EXISTS notes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    content TEXT DEFAULT '',
    folder_id INTEGER REFERENCES folders(id) ON DELETE SET NULL,
    template_id INTEGER REFERENCES templates(id) ON DELETE SET NULL,

    -- Todo specific fields
    is_todo BOOLEAN DEFAULT FALSE,
    is_done BOOLEAN DEFAULT FALSE,
    priority INTEGER DEFAULT 0 CHECK(priority >= 0 AND priority <= 3),
    due_date DATETIME,

    -- Metadata
    tags TEXT DEFAULT '', -- Comma-separated or JSON array
    starred BOOLEAN DEFAULT FALSE,

    -- Timestamps
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_notes_folder ON notes(folder_id);
CREATE INDEX idx_notes_is_todo ON notes(is_todo);
CREATE INDEX idx_notes_is_done ON notes(is_done);
CREATE INDEX idx_notes_priority ON notes(priority);
CREATE INDEX idx_notes_starred ON notes(starred);
CREATE INDEX idx_notes_created ON notes(created_at);

-- Full-text search virtual table
CREATE VIRTUAL TABLE IF NOT EXISTS notes_fts USING fts5(
    title,
    content,
    tags,
    content='notes',
    content_rowid='id'
);

-- Triggers to keep FTS in sync
CREATE TRIGGER notes_ai AFTER INSERT ON notes BEGIN
    INSERT INTO notes_fts(rowid, title, content, tags)
    VALUES (new.id, new.title, new.content, new.tags);
END;

CREATE TRIGGER notes_ad AFTER DELETE ON notes BEGIN
    INSERT INTO notes_fts(notes_fts, rowid, title, content, tags)
    VALUES ('delete', old.id, old.title, old.content, old.tags);
END;

CREATE TRIGGER notes_au AFTER UPDATE ON notes BEGIN
    INSERT INTO notes_fts(notes_fts, rowid, title, content, tags)
    VALUES ('delete', old.id, old.title, old.content, old.tags);
    INSERT INTO notes_fts(rowid, title, content, tags)
    VALUES (new.id, new.title, new.content, new.tags);
END;

-- Insert default folders
INSERT INTO folders (name, icon) VALUES
    ('Work', '💼'),
    ('Personal', '🏠'),
    ('Ideas', '💡');

-- Insert default templates
INSERT INTO templates (name, description, content, type, icon, is_default) VALUES
    ('Blank Note', 'Empty note', '', 'note', '📝', TRUE),
    ('Blank Todo', 'Empty todo item', '', 'todo', '☐', TRUE);
```

---

## Core Features

### 1. Notes Management

| Feature     | Description                                 | Priority |
| ----------- | ------------------------------------------- | -------- |
| Create note | Create new note with title, content, folder | P0       |
| Edit note   | Open in vim/neovim                          | P0       |
| Delete note | Soft delete or hard delete                  | P0       |
| Move note   | Change folder                               | P1       |
| Star note   | Mark as favorite                            | P1       |
| Tags        | Add/remove tags                             | P1       |

### 2. Folder Management

| Feature         | Description                  | Priority |
| --------------- | ---------------------------- | -------- |
| Create folder   | With custom icon             | P0       |
| Rename folder   | Change name                  | P0       |
| Delete folder   | Move notes to root or delete | P0       |
| Nested folders  | Parent-child relationship    | P1       |
| Reorder folders | Drag and drop position       | P2       |

### 3. Todo Management

| Feature      | Description                             | Priority |
| ------------ | --------------------------------------- | -------- |
| Create todo  | Quick todo creation                     | P0       |
| Toggle done  | Mark complete/incomplete                | P0       |
| Set priority | High (3), Medium (2), Low (1), None (0) | P0       |
| Due date     | Optional deadline                       | P1       |
| Filter todos | By status, priority, folder             | P1       |

### 4. Search

| Feature             | Description                    | Priority |
| ------------------- | ------------------------------ | -------- |
| Full-text search    | Search in title, content, tags | P0       |
| Filter by folder    | Scope search to folder         | P1       |
| Filter by type      | Notes only, todos only         | P1       |
| Search highlighting | Highlight matches              | P2       |

### 5. Templates

| Feature            | Description                             | Priority |
| ------------------ | --------------------------------------- | -------- |
| Use template       | Create note from template               | P0       |
| Create template    | Save note as template                   | P1       |
| Edit template      | Modify existing template                | P1       |
| Template variables | Placeholders like {{date}}, {{project}} | P2       |

---

## Templates System

### Default Templates

```yaml
# Template 1: Meeting Notes
name: "Meeting Notes"
type: "note"
icon: "🤝"
content: |
  # {{title}}

  **Date:** {{date}}
  **Attendees:** 

  ## Agenda
  - 

  ## Discussion Points


  ## Action Items
  - [ ] 

  ## Next Steps

---
# Template 2: Daily Standup
name: "Daily Standup"
type: "note"
icon: "☀️"
content: |
  # Standup - {{date}}

  ## Yesterday
  - 

  ## Today
  - 

  ## Blockers
  - None

---
# Template 3: Bug Report
name: "Bug Report"
type: "todo"
icon: "🐛"
content: |
  ## Description


  ## Steps to Reproduce
  1. 
  2. 
  3. 

  ## Expected Behavior


  ## Actual Behavior


  ## Environment
  - OS: 
  - Version:

---
# Template 4: Feature Request
name: "Feature Request"
type: "todo"
icon: "✨"
content: |
  ## Summary


  ## User Story
  As a [user type], I want [goal] so that [benefit].

  ## Acceptance Criteria
  - [ ] 
  - [ ] 

  ## Technical Notes

---
# Template 5: Weekly Review
name: "Weekly Review"
type: "note"
icon: "📅"
content: |
  # Week {{week_number}} Review

  ## Accomplishments
  - 

  ## Challenges
  - 

  ## Learnings
  - 

  ## Next Week Goals
  - [ ] 
  - [ ]

---
# Template 6: Quick Todo
name: "Quick Todo"
type: "todo"
icon: "⚡"
content: ""

---
# Template 7: Project Note
name: "Project Note"
type: "note"
icon: "📊"
content: |
  # Project: {{title}}

  ## Overview


  ## Goals
  - 

  ## Timeline
  | Phase | Start | End | Status |
  |-------|-------|-----|--------|
  |       |       |     |        |

  ## Resources
  - 

  ## Notes
```

### Template Variables

| Variable          | Description          | Example Output     |
| ----------------- | -------------------- | ------------------ |
| `{{title}}`       | Note title           | "Sprint Planning"  |
| `{{date}}`        | Current date         | "2026-01-02"       |
| `{{datetime}}`    | Current datetime     | "2026-01-02 16:30" |
| `{{week_number}}` | Current week number  | "1"                |
| `{{folder}}`      | Target folder name   | "Work"             |
| `{{author}}`      | Username from config | "quang"            |

---

## TUI Components

### Component Hierarchy

```
App (BubbleTea Model)
│
├── Header
│   └── Title + Breadcrumb + Time
│
├── Content (depends on current view)
│   │
│   ├── MainView
│   │   ├── Sidebar (folders)
│   │   └── NoteList
│   │
│   ├── NoteView
│   │   ├── NoteList (narrow)
│   │   └── Preview (wide)
│   │
│   ├── TodoView
│   │   └── TodoList (grouped by priority)
│   │
│   ├── SearchView
│   │   ├── SearchBar
│   │   └── ResultsList
│   │
│   └── TemplatePicker (modal)
│       └── TemplateList
│
├── StatusBar
│   └── Keybindings + Stats
│
└── Dialog (modal, optional)
    └── Form fields
```

### Keybindings

```go
// Global
"q", "ctrl+c" → Quit
"?"           → Show help
"/"           → Search
"ctrl+n"      → New note
"ctrl+t"      → New todo

// Navigation
"j", "down"   → Move down
"k", "up"     → Move up
"h", "left"   → Collapse/Back
"l", "right"  → Expand/Enter
"gg"          → Go to top
"G"           → Go to bottom
"tab"         → Switch panel

// Actions
"enter"       → Open/Select
"e"           → Edit in vim
"d"           → Delete (with confirm)
"m"           → Move to folder
"s"           → Toggle star
"t"           → Toggle todo status
"p"           → Change priority (for todos)
"T"           → Pick template
"y"           → Copy to clipboard

// Folders
"a"           → Add folder
"r"           → Rename folder
```

### Color Palette

```go
// styles.go
var (
    // Base colors
    Primary    = lipgloss.Color("#7C3AED") // Purple
    Secondary  = lipgloss.Color("#06B6D4") // Cyan
    Success    = lipgloss.Color("#10B981") // Green
    Warning    = lipgloss.Color("#F59E0B") // Amber
    Danger     = lipgloss.Color("#EF4444") // Red

    // Priorities
    PriorityHigh   = lipgloss.Color("#EF4444") // Red
    PriorityMedium = lipgloss.Color("#F59E0B") // Amber
    PriorityLow    = lipgloss.Color("#10B981") // Green

    // UI elements
    Background     = lipgloss.Color("#1E1E2E") // Dark
    Surface        = lipgloss.Color("#313244") // Slightly lighter
    Border         = lipgloss.Color("#45475A") // Border
    TextPrimary    = lipgloss.Color("#CDD6F4") // Light text
    TextSecondary  = lipgloss.Color("#A6ADC8") // Dimmed text
    TextMuted      = lipgloss.Color("#6C7086") // Very dimmed
)
```

---

## Implementation Phases

### Phase 1: Foundation (Week 1)

**Goal:** Project setup, database, basic models

| Task | Description                                | Est. |
| ---- | ------------------------------------------ | ---- |
| 1.1  | Initialize Go module, install dependencies | 1h   |
| 1.2  | Setup project structure                    | 1h   |
| 1.3  | Implement config loading (Viper)           | 2h   |
| 1.4  | Implement database connection & migrations | 3h   |
| 1.5  | Implement models (Note, Folder, Template)  | 2h   |
| 1.6  | Write unit tests for models                | 2h   |
| 1.7  | Implement repositories (CRUD)              | 4h   |
| 1.8  | Write unit tests for repositories          | 3h   |

**Deliverables:**

- [ ] Working database with migrations
- [ ] All models with validation
- [ ] Repository layer with full CRUD
- [ ] 80%+ test coverage for repos

---

### Phase 2: Services & Editor (Week 2)

**Goal:** Business logic, vim integration

| Task | Description                                | Est. |
| ---- | ------------------------------------------ | ---- |
| 2.1  | Implement NoteService                      | 3h   |
| 2.2  | Implement FolderService                    | 2h   |
| 2.3  | Implement TemplateService                  | 3h   |
| 2.4  | Implement SearchService (FTS5)             | 3h   |
| 2.5  | Implement EditorService (vim/neovim spawn) | 3h   |
| 2.6  | Write unit tests for services              | 4h   |
| 2.7  | Integration tests for editor               | 2h   |

**Deliverables:**

- [ ] All services implemented
- [ ] Vim/Neovim integration working
- [ ] Full-text search working
- [ ] Template variable substitution

---

### Phase 3: TUI Basic (Week 3)

**Goal:** Basic TUI with main view

| Task | Description                             | Est. |
| ---- | --------------------------------------- | ---- |
| 3.1  | Setup Bubble Tea app structure          | 2h   |
| 3.2  | Implement styles (Lipgloss)             | 2h   |
| 3.3  | Implement Sidebar component             | 3h   |
| 3.4  | Implement NoteList component            | 3h   |
| 3.5  | Implement Preview component             | 2h   |
| 3.6  | Implement StatusBar component           | 1h   |
| 3.7  | Implement MainView (combine components) | 3h   |
| 3.8  | Implement keybindings                   | 2h   |
| 3.9  | Navigation between panels               | 2h   |

**Deliverables:**

- [ ] Main dashboard showing folders and notes
- [ ] Navigation with keyboard
- [ ] Note preview

---

### Phase 4: TUI Advanced (Week 4)

**Goal:** All views, dialogs, search

| Task | Description                       | Est. |
| ---- | --------------------------------- | ---- |
| 4.1  | Implement SearchView              | 3h   |
| 4.2  | Implement TodoView                | 3h   |
| 4.3  | Implement NewNote dialog          | 3h   |
| 4.4  | Implement TemplatePicker modal    | 3h   |
| 4.5  | Implement Delete confirmation     | 1h   |
| 4.6  | Implement Help overlay            | 2h   |
| 4.7  | Implement folder CRUD in TUI      | 2h   |
| 4.8  | Polish animations and transitions | 2h   |

**Deliverables:**

- [ ] All views working
- [ ] Search with results
- [ ] Template picker
- [ ] Full CRUD from TUI

---

### Phase 5: CLI & Polish (Week 5)

**Goal:** CLI commands, final polish

| Task | Description                        | Est. |
| ---- | ---------------------------------- | ---- |
| 5.1  | Implement CLI root command         | 1h   |
| 5.2  | Implement `kiroku add`             | 2h   |
| 5.3  | Implement `kiroku list`            | 2h   |
| 5.4  | Implement `kiroku search`          | 2h   |
| 5.5  | Implement `kiroku todo`            | 2h   |
| 5.6  | Add default templates on first run | 2h   |
| 5.7  | Error handling & edge cases        | 3h   |
| 5.8  | Performance optimization           | 2h   |
| 5.9  | Documentation (README)             | 2h   |

**Deliverables:**

- [ ] CLI working alongside TUI
- [ ] Default templates seeded
- [ ] README with usage instructions
- [ ] Release-ready binary

---

## Testing Strategy

### Test Structure

```
kiroku/
├── internal/
│   ├── repository/
│   │   ├── note_repo.go
│   │   └── note_repo_test.go      # Unit tests
│   ├── service/
│   │   ├── note_service.go
│   │   └── note_service_test.go   # Unit tests with mocks
│   └── tui/
│       └── components/
│           ├── notelist.go
│           └── notelist_test.go   # Component tests
│
├── tests/
│   ├── integration/
│   │   ├── database_test.go       # DB integration
│   │   └── editor_test.go         # Editor integration
│   └── e2e/
│       └── tui_test.go            # End-to-end TUI tests
│
└── testdata/
    └── fixtures/
        ├── notes.json
        └── templates.json
```

### Test Cases

#### Repository Tests

```go
// note_repo_test.go

func TestNoteRepository_Create(t *testing.T) {
    // Setup
    db := setupTestDB(t)
    repo := NewNoteRepository(db)

    tests := []struct {
        name    string
        note    *models.Note
        wantErr bool
    }{
        {
            name: "valid note",
            note: &models.Note{
                Title:   "Test Note",
                Content: "Test content",
            },
            wantErr: false,
        },
        {
            name: "empty title",
            note: &models.Note{
                Title:   "",
                Content: "Content",
            },
            wantErr: true,
        },
        {
            name: "with folder",
            note: &models.Note{
                Title:    "Note in folder",
                FolderID: ptr(1),
            },
            wantErr: false,
        },
        {
            name: "with invalid folder",
            note: &models.Note{
                Title:    "Note",
                FolderID: ptr(9999),
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := repo.Create(tt.note)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.NotZero(t, tt.note.ID)
            }
        })
    }
}

func TestNoteRepository_GetByID(t *testing.T) {
    db := setupTestDB(t)
    repo := NewNoteRepository(db)

    // Create test note
    note := &models.Note{Title: "Test", Content: "Content"}
    repo.Create(note)

    tests := []struct {
        name    string
        id      int64
        want    *models.Note
        wantErr bool
    }{
        {
            name:    "existing note",
            id:      note.ID,
            want:    note,
            wantErr: false,
        },
        {
            name:    "non-existing note",
            id:      9999,
            want:    nil,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := repo.GetByID(tt.id)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.want.Title, got.Title)
            }
        })
    }
}

func TestNoteRepository_List(t *testing.T) {
    db := setupTestDB(t)
    repo := NewNoteRepository(db)

    // Create test data
    for i := 0; i < 5; i++ {
        repo.Create(&models.Note{Title: fmt.Sprintf("Note %d", i)})
    }

    tests := []struct {
        name     string
        opts     ListOptions
        wantLen  int
    }{
        {
            name:    "all notes",
            opts:    ListOptions{},
            wantLen: 5,
        },
        {
            name:    "with limit",
            opts:    ListOptions{Limit: 3},
            wantLen: 3,
        },
        {
            name:    "with offset",
            opts:    ListOptions{Offset: 2, Limit: 10},
            wantLen: 3,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            notes, err := repo.List(tt.opts)
            assert.NoError(t, err)
            assert.Len(t, notes, tt.wantLen)
        })
    }
}

func TestNoteRepository_Update(t *testing.T) {
    db := setupTestDB(t)
    repo := NewNoteRepository(db)

    note := &models.Note{Title: "Original", Content: "Original content"}
    repo.Create(note)

    // Update
    note.Title = "Updated"
    note.Content = "Updated content"
    err := repo.Update(note)

    assert.NoError(t, err)

    // Verify
    updated, _ := repo.GetByID(note.ID)
    assert.Equal(t, "Updated", updated.Title)
    assert.Equal(t, "Updated content", updated.Content)
}

func TestNoteRepository_Delete(t *testing.T) {
    db := setupTestDB(t)
    repo := NewNoteRepository(db)

    note := &models.Note{Title: "To Delete"}
    repo.Create(note)

    err := repo.Delete(note.ID)
    assert.NoError(t, err)

    // Verify deleted
    _, err = repo.GetByID(note.ID)
    assert.Error(t, err)
}

func TestNoteRepository_Search(t *testing.T) {
    db := setupTestDB(t)
    repo := NewNoteRepository(db)

    // Create test notes
    repo.Create(&models.Note{Title: "Meeting notes", Content: "Discussed project timeline"})
    repo.Create(&models.Note{Title: "Shopping list", Content: "Buy groceries"})
    repo.Create(&models.Note{Title: "Project ideas", Content: "New meeting app concept"})

    tests := []struct {
        name    string
        query   string
        wantLen int
    }{
        {
            name:    "search by title",
            query:   "meeting",
            wantLen: 2, // "Meeting notes" and "New meeting app"
        },
        {
            name:    "search by content",
            query:   "groceries",
            wantLen: 1,
        },
        {
            name:    "no results",
            query:   "nonexistent",
            wantLen: 0,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            results, err := repo.Search(tt.query)
            assert.NoError(t, err)
            assert.Len(t, results, tt.wantLen)
        })
    }
}
```

#### Service Tests

```go
// note_service_test.go

func TestNoteService_CreateFromTemplate(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockNoteRepo := mocks.NewMockNoteRepository(ctrl)
    mockTemplateRepo := mocks.NewMockTemplateRepository(ctrl)

    service := NewNoteService(mockNoteRepo, mockTemplateRepo)

    template := &models.Template{
        ID:      1,
        Name:    "Meeting Notes",
        Content: "# Meeting: {{title}}\nDate: {{date}}",
        Type:    "note",
    }

    mockTemplateRepo.EXPECT().
        GetByID(int64(1)).
        Return(template, nil)

    mockNoteRepo.EXPECT().
        Create(gomock.Any()).
        DoAndReturn(func(note *models.Note) error {
            note.ID = 1
            return nil
        })

    note, err := service.CreateFromTemplate(1, "Sprint Planning", nil)

    assert.NoError(t, err)
    assert.Contains(t, note.Content, "Sprint Planning")
    assert.Contains(t, note.Content, time.Now().Format("2006-01-02"))
}

func TestNoteService_ToggleTodo(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockNoteRepository(ctrl)
    service := NewNoteService(mockRepo, nil)

    todo := &models.Note{
        ID:     1,
        Title:  "Test todo",
        IsTodo: true,
        IsDone: false,
    }

    mockRepo.EXPECT().GetByID(int64(1)).Return(todo, nil)
    mockRepo.EXPECT().Update(gomock.Any()).Return(nil)

    err := service.ToggleTodo(1)

    assert.NoError(t, err)
    assert.True(t, todo.IsDone)
}
```

#### Editor Service Tests

```go
// editor_service_test.go

func TestEditorService_GetEditor(t *testing.T) {
    tests := []struct {
        name     string
        envVar   string
        expected string
    }{
        {
            name:     "EDITOR set to nvim",
            envVar:   "nvim",
            expected: "nvim",
        },
        {
            name:     "EDITOR set to vim",
            envVar:   "vim",
            expected: "vim",
        },
        {
            name:     "EDITOR not set",
            envVar:   "",
            expected: "vim", // default
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if tt.envVar != "" {
                os.Setenv("EDITOR", tt.envVar)
            } else {
                os.Unsetenv("EDITOR")
            }

            service := NewEditorService()
            assert.Equal(t, tt.expected, service.GetEditor())
        })
    }
}

func TestEditorService_CreateTempFile(t *testing.T) {
    service := NewEditorService()

    content := "# Test Note\n\nSome content"

    path, cleanup, err := service.CreateTempFile(content)
    defer cleanup()

    assert.NoError(t, err)
    assert.FileExists(t, path)

    data, _ := os.ReadFile(path)
    assert.Equal(t, content, string(data))
}
```

#### TUI Component Tests

```go
// notelist_test.go

func TestNoteList_Render(t *testing.T) {
    notes := []models.Note{
        {ID: 1, Title: "Note 1", IsTodo: false},
        {ID: 2, Title: "Todo 1", IsTodo: true, IsDone: false},
        {ID: 3, Title: "Todo 2", IsTodo: true, IsDone: true},
    }

    component := NewNoteList(notes)

    output := component.View()

    assert.Contains(t, output, "Note 1")
    assert.Contains(t, output, "☐")  // Pending todo
    assert.Contains(t, output, "☑")  // Done todo
}

func TestNoteList_Navigation(t *testing.T) {
    notes := []models.Note{
        {ID: 1, Title: "Note 1"},
        {ID: 2, Title: "Note 2"},
        {ID: 3, Title: "Note 3"},
    }

    component := NewNoteList(notes)

    // Initial state
    assert.Equal(t, 0, component.Cursor())

    // Move down
    component, _ = component.Update(tea.KeyMsg{Type: tea.KeyDown})
    assert.Equal(t, 1, component.Cursor())

    // Move up
    component, _ = component.Update(tea.KeyMsg{Type: tea.KeyUp})
    assert.Equal(t, 0, component.Cursor())

    // Boundary check - can't go above 0
    component, _ = component.Update(tea.KeyMsg{Type: tea.KeyUp})
    assert.Equal(t, 0, component.Cursor())
}
```

#### Integration Tests

```go
// tests/integration/database_test.go

func TestDatabase_MigrationAndSeed(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    // Use temp file for test DB
    tmpFile, _ := os.CreateTemp("", "kiroku_test_*.db")
    defer os.Remove(tmpFile.Name())

    db, err := database.New(tmpFile.Name())
    assert.NoError(t, err)

    // Run migrations
    err = db.Migrate()
    assert.NoError(t, err)

    // Check default folders exist
    var count int
    err = db.QueryRow("SELECT COUNT(*) FROM folders").Scan(&count)
    assert.NoError(t, err)
    assert.Equal(t, 3, count) // Work, Personal, Ideas

    // Check default templates exist
    err = db.QueryRow("SELECT COUNT(*) FROM templates").Scan(&count)
    assert.NoError(t, err)
    assert.GreaterOrEqual(t, count, 2) // At least blank note and todo
}

func TestDatabase_FTSSearch(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    db := setupIntegrationDB(t)

    // Insert test data
    db.Exec(`INSERT INTO notes (title, content) VALUES
        ('Meeting notes', 'Discussed quarterly goals'),
        ('Shopping list', 'Buy milk and eggs'),
        ('Project plan', 'Meeting with stakeholders')`)

    // Search
    rows, err := db.Query(`
        SELECT n.id, n.title
        FROM notes_fts
        JOIN notes n ON notes_fts.rowid = n.id
        WHERE notes_fts MATCH 'meeting'
    `)
    assert.NoError(t, err)

    var results []string
    for rows.Next() {
        var id int
        var title string
        rows.Scan(&id, &title)
        results = append(results, title)
    }

    assert.Len(t, results, 2)
    assert.Contains(t, results, "Meeting notes")
    assert.Contains(t, results, "Project plan")
}
```

### Test Commands

```makefile
# Makefile

.PHONY: test test-unit test-integration test-coverage

# Run all tests
test:
	go test ./... -v

# Run unit tests only
test-unit:
	go test ./internal/... -v -short

# Run integration tests
test-integration:
	go test ./tests/integration/... -v

# Run with coverage
test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Run specific package tests
test-repo:
	go test ./internal/repository/... -v

test-service:
	go test ./internal/service/... -v

test-tui:
	go test ./internal/tui/... -v
```

---

## CLI Commands

### Command Reference

```bash
# Launch TUI (default)
kiroku

# Quick add note
kiroku add "Note title"
kiroku add "Note title" -f work                    # with folder
kiroku add "Note title" -t meeting-notes           # with template

# Quick add todo
kiroku todo "Review PR #123"
kiroku todo "Deploy to prod" -p high               # with priority
kiroku todo "Update docs" -d 2026-01-05            # with due date

# List notes
kiroku list                                        # all notes
kiroku list -f work                                # by folder
kiroku list --todos                                # todos only
kiroku list --todos --pending                      # pending todos

# Search
kiroku search "meeting"
kiroku search "project" -f work

# Templates
kiroku templates                                   # list templates
kiroku templates add "My Template" --from note.md  # create from file

# Edit note by ID
kiroku edit 123

# Open specific folder in TUI
kiroku open work
```

---

## Configuration

### Config File Location

```
~/.config/kiroku/config.yaml
```

### Config Schema

```yaml
# ~/.config/kiroku/config.yaml

# Database location
database:
  path: ~/.local/share/kiroku/kiroku.db

# Editor preference
editor:
  command: nvim # or vim, nano, etc.
  args: ["-c", "set filetype=markdown"]

# Default folder for new notes
defaults:
  folder: "" # empty = root
  template: "" # empty = blank note

# UI preferences
ui:
  theme: dark # dark, light
  show_preview: true
  date_format: "Jan 2, 15:04"
  sidebar_width: 25 # percentage

# Todo settings
todos:
  show_completed: true
  sort_by: priority # priority, due_date, created_at
```

---

## Summary

### MVP Features (Phase 1-4)

- [x] SQLite database with FTS5
- [x] Note and Todo CRUD
- [x] Folder organization
- [x] Template system with variables
- [x] TUI with Bubble Tea
- [x] Vim/Neovim integration
- [x] Full-text search

### Future Enhancements (v2+)

- [ ] Tags with autocomplete
- [ ] Note linking `[[note-name]]`
- [ ] Export to Markdown/PDF
- [ ] Import from other apps
- [ ] Recurring todos
- [ ] Reminders/notifications
- [ ] Cloud sync (optional)

---

**Ready to start implementation?** 🚀
