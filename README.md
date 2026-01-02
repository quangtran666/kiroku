# 記録 Kiroku

A lightweight, terminal-based note-taking application with a beautiful TUI interface built in Go.

![Kiroku Demo](docs/demo.gif)

## ✨ Features

- ⚡ **Fast and lightweight** - Built with Go, uses SQLite for storage
- 🔍 **Full-text search** - Powered by SQLite FTS5
- 📁 **Folder organization** - Organize notes in nested folders
- ☐ **Todo management** - Create todos with priorities and due dates
- 📝 **Templates** - Quick note creation with customizable templates
- 🖥️ **Beautiful TUI** - Modern terminal interface with Bubble Tea
- ✏️ **Vim integration** - Edit notes in your favorite editor

## 📦 Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/tranducquang/kiroku.git
cd kiroku

# Build
make build

# Install to GOPATH/bin
make install
```

### Using Go

```bash
go install github.com/tranducquang/kiroku/cmd/kiroku@latest
```

## 🚀 Quick Start

```bash
# Launch TUI
kiroku

# Quick add a note
kiroku add "Meeting notes"

# Add a todo
kiroku todo "Review PR #123" -p high

# List notes
kiroku list

# Search
kiroku search "meeting"
```

## ⌨️ Keyboard Shortcuts

### Navigation

| Key       | Action          |
| --------- | --------------- |
| `j/k`     | Move down/up    |
| `h`       | Collapse folder |
| `l/Enter` | Expand folder   |
| `Tab`     | Switch panel    |

### Actions

| Key       | Action          |
| --------- | --------------- |
| `n`       | New note        |
| `t`       | New todo        |
| `f`       | New folder      |
| `e`       | Edit in vim     |
| `d`       | Delete          |
| `s`       | Toggle star     |
| `x/Space` | Toggle done     |
| `p`       | Change priority |
| `/`       | Search          |
| `?`       | Help            |
| `v`       | Toggle preview  |
| `q`       | Quit            |

## 📋 CLI Commands

```bash
# Launch TUI (default)
kiroku

# Quick add note
kiroku add "Note title"
kiroku add "Note title" -f work              # with folder
kiroku add "Note title" -t meeting-notes     # with template

# Quick add todo
kiroku todo "Todo title"
kiroku todo "Todo" -p high                   # with priority
kiroku todo "Todo" -d 2026-01-05            # with due date

# List notes
kiroku list                                  # all notes
kiroku list -f work                          # by folder
kiroku list --todos                          # todos only
kiroku list --todos --pending                # pending todos

# Search
kiroku search "query"
kiroku search "query" -f work                # search in folder

# Edit by ID
kiroku edit 123

# Templates
kiroku templates                             # list templates
```

## ⚙️ Configuration

Configuration file: `~/.config/kiroku/config.yaml`

```yaml
# Database location
database:
  path: ~/.local/share/kiroku/kiroku.db

# Editor preference
editor:
  command: nvim
  args: ["-c", "set filetype=markdown"]

# Default settings
defaults:
  folder: ""
  template: ""

# UI preferences
ui:
  theme: dark
  show_preview: true
  date_format: "Jan 2, 15:04"
  sidebar_width: 25

# Todo settings
todos:
  show_completed: true
  sort_by: priority
```

## 📝 Templates

Built-in templates:

- 📝 Blank Note
- ☐ Blank Todo
- 🤝 Meeting Notes
- ☀️ Daily Standup
- 🐛 Bug Report
- 📅 Weekly Review

Templates support variables:

- `{{title}}` - Note title
- `{{date}}` - Current date
- `{{datetime}}` - Current date and time
- `{{week_number}}` - ISO week number

## 🗂️ Project Structure

```
kiroku/
├── cmd/kiroku/main.go          # Entry point
├── internal/
│   ├── app/                    # Application orchestrator
│   ├── config/                 # Configuration
│   ├── database/               # Database & migrations
│   ├── models/                 # Data models
│   ├── repository/             # Data access layer
│   ├── service/                # Business logic
│   ├── tui/                    # TUI components
│   └── cli/                    # CLI commands
├── docs/                       # Documentation
├── Makefile
└── README.md
```

## 🔧 Development

```bash
# Install dependencies
make deps

# Run in development
make run

# Build
make build

# Run tests
make test

# Run with coverage
make test-coverage

# Format code
make fmt

# Lint
make lint
```

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [SQLite](https://www.sqlite.org/) - Database

---

Made with ❤️ by the Kiroku team
