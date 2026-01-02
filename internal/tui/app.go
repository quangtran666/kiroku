package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/tranducquang/kiroku/internal/config"
	"github.com/tranducquang/kiroku/internal/logging"
	"github.com/tranducquang/kiroku/internal/models"
	"github.com/tranducquang/kiroku/internal/service"
	"github.com/tranducquang/kiroku/internal/tui/components"
	"github.com/tranducquang/kiroku/internal/tui/keys"
	"github.com/tranducquang/kiroku/internal/tui/styles"
)

// ViewType represents the current view
type ViewType int

const (
	ViewMain ViewType = iota
	ViewNote
	ViewTodo
	ViewSearch
	ViewNewNote
	ViewNewTodo
	ViewTemplate
)

// Panel represents the focused panel
type Panel int

const (
	PanelSidebar Panel = iota
	PanelNoteList
	PanelPreview
)

// App represents the main TUI application
type App struct {
	// Services
	noteService     *service.NoteService
	folderService   *service.FolderService
	templateService *service.TemplateService
	searchService   *service.SearchService
	editorService   *service.EditorService
	cfg             *config.Config

	// State
	currentView  ViewType
	currentPanel Panel
	width        int
	height       int
	ready        bool

	// Data
	folders   []*models.Folder
	notes     []*models.Note
	templates []models.Template

	// Current selection
	currentFolder *models.Folder
	currentNote   *models.Note
	currentFilter string // "all", "todos", "starred"

	// Components
	sidebar   *components.Sidebar
	noteList  *components.NoteList
	preview   *components.Preview
	statusBar *components.StatusBar
	searchBar *components.SearchBar
	help      *components.Help
	dialog    *components.Dialog

	// Flags
	showHelp        bool
	showDialog      bool
	dialogType      string
	searchMode      bool
	searchQuery     string
	editingTempFile string // Temp file path when editing a note
}

// NewApp creates a new TUI application
func NewApp(
	noteService *service.NoteService,
	folderService *service.FolderService,
	templateService *service.TemplateService,
	searchService *service.SearchService,
	editorService *service.EditorService,
	cfg *config.Config,
) *App {
	noteList := components.NewNoteList()
	noteList.SetFocused(true)

	return &App{
		noteService:     noteService,
		folderService:   folderService,
		templateService: templateService,
		searchService:   searchService,
		editorService:   editorService,
		cfg:             cfg,
		currentView:     ViewMain,
		currentPanel:    PanelNoteList,
		currentFilter:   "all",
		sidebar:         components.NewSidebar(),
		noteList:        noteList,
		preview:         components.NewPreview(),
		statusBar:       components.NewStatusBar(),
		searchBar:       components.NewSearchBar(),
		help:            components.NewHelp(),
		dialog:          components.NewDialog(),
	}
}

// Init initializes the application
func (a *App) Init() tea.Cmd {
	return tea.Batch(
		a.loadData,
		tea.SetWindowTitle("記録 Kiroku"),
	)
}

// loadData loads initial data
func (a *App) loadData() tea.Msg {
	ctx := context.Background()

	// Load folders
	folders, err := a.folderService.GetTree(ctx)
	if err != nil {
		return errMsg{err}
	}

	// Load notes
	notes, err := a.noteService.GetAllNotes(ctx)
	if err != nil {
		return errMsg{err}
	}

	// Load templates
	templates, err := a.templateService.List(ctx)
	if err != nil {
		return errMsg{err}
	}

	return dataLoadedMsg{
		folders:   folders,
		notes:     notes,
		templates: templates,
	}
}

// Message types
type dataLoadedMsg struct {
	folders   []*models.Folder
	notes     []*models.Note
	templates []models.Template
}

type errMsg struct {
	err error
}

type statusClearMsg struct{}

type editorFinishedMsg struct{}

// Update handles messages
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		logging.Debug().Int("width", msg.Width).Int("height", msg.Height).Msg("Window resize")
		a.width = msg.Width
		a.height = msg.Height
		a.ready = true
		a.updateLayout()
		return a, nil

	case dataLoadedMsg:
		logging.Debug().
			Int("folders", len(msg.folders)).
			Int("notes", len(msg.notes)).
			Int("templates", len(msg.templates)).
			Msg("Data loaded")
		a.folders = msg.folders
		a.notes = msg.notes
		a.templates = msg.templates
		a.sidebar.SetFolders(a.folders)
		a.noteList.SetNotes(a.notes)
		a.updatePreview()
		return a, nil

	case errMsg:
		logging.Error().Err(msg.err).Msg("TUI error")
		a.statusBar.SetMessage(fmt.Sprintf("Error: %v", msg.err))
		return a, a.clearStatusAfter(3 * time.Second)

	case statusClearMsg:
		a.statusBar.ClearMessage()
		return a, nil

	case editorFinishedMsg:
		logging.Debug().Msg("Processing editor finished message")
		if a.currentNote == nil || a.editingTempFile == "" {
			logging.Warn().Msg("No note or temp file to process")
			return a, nil
		}

		// Read edited content
		newTitle, newContent, err := a.editorService.ReadEditedContent(a.editingTempFile, a.currentNote.Title)
		if err != nil {
			logging.Error().Err(err).Msg("Failed to read edited content")
			a.statusBar.SetMessage(fmt.Sprintf("Error: %v", err))
			a.editingTempFile = ""
			return a, a.clearStatusAfter(3 * time.Second)
		}

		// Update note
		ctx := context.Background()
		a.currentNote.Title = newTitle
		a.currentNote.Content = newContent
		if err := a.noteService.Update(ctx, a.currentNote); err != nil {
			logging.Error().Err(err).Msg("Failed to update note")
			a.statusBar.SetMessage(fmt.Sprintf("Error: %v", err))
			a.editingTempFile = ""
			return a, a.clearStatusAfter(3 * time.Second)
		}

		logging.Info().Int64("note_id", a.currentNote.ID).Msg("Note updated successfully")
		a.statusBar.SetMessage("Note saved")
		a.editingTempFile = ""
		return a, tea.Batch(a.reloadNotes, a.clearStatusAfter(2*time.Second))

	case tea.KeyMsg:
		logging.Debug().
			Str("key", msg.String()).
			Str("panel", panelName(a.currentPanel)).
			Bool("dialog", a.showDialog).
			Bool("search", a.searchMode).
			Bool("help", a.showHelp).
			Msg("Key pressed")

		// Global keys
		if a.showHelp {
			a.help, _ = a.help.Update(msg)
			if !a.help.IsVisible() {
				a.showHelp = false
			}
			return a, nil
		}

		if a.showDialog {
			logging.Debug().Str("dialog_type", a.dialogType).Msg("Handling dialog input")
			return a.handleDialogInput(msg)
		}

		if a.searchMode {
			logging.Debug().Msg("Handling search input")
			return a.handleSearchInput(msg)
		}

		// Handle quit
		if key.Matches(msg, keys.DefaultKeyMap.Quit) {
			logging.Info().Msg("User quit application")
			return a, tea.Quit
		}

		// Handle help
		if key.Matches(msg, keys.DefaultKeyMap.Help) {
			logging.Debug().Msg("Showing help")
			a.showHelp = true
			a.help.Show()
			return a, nil
		}

		// Handle search
		if key.Matches(msg, keys.DefaultKeyMap.Search) {
			logging.Debug().Msg("Entering search mode")
			a.searchMode = true
			a.searchBar.Focus()
			return a, nil
		}

		// Handle new note/todo
		if key.Matches(msg, keys.DefaultKeyMap.NewNote) {
			logging.Debug().Msg("Showing new note dialog")
			return a.showNewNoteDialog()
		}

		if key.Matches(msg, keys.DefaultKeyMap.NewTodo) {
			logging.Debug().Msg("Showing new todo dialog")
			return a.showNewTodoDialog()
		}

		// Handle panel switching with Tab, Left, Right
		if key.Matches(msg, keys.DefaultKeyMap.Tab) || key.Matches(msg, keys.DefaultKeyMap.Right) {
			logging.Debug().Msg("Switching panel right")
			a.switchPanel(1)
			return a, nil
		}

		if key.Matches(msg, keys.DefaultKeyMap.Left) {
			logging.Debug().Msg("Switching panel left")
			a.switchPanel(-1)
			return a, nil
		}

		// Handle panel-specific input
		logging.Debug().Str("panel", panelName(a.currentPanel)).Msg("Handling panel input")
		return a.handlePanelInput(msg)
	}

	return a, tea.Batch(cmds...)
}

// panelName returns the name of a panel for logging
func panelName(p Panel) string {
	switch p {
	case PanelSidebar:
		return "sidebar"
	case PanelNoteList:
		return "notelist"
	case PanelPreview:
		return "preview"
	default:
		return "unknown"
	}
}

func (a *App) handlePanelInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch a.currentPanel {
	case PanelSidebar:
		a.sidebar, _ = a.sidebar.Update(msg)

		// Check if selection changed
		special := a.sidebar.SelectedSpecial()
		if special != "" && special != a.currentFilter {
			logging.Debug().Str("filter", special).Msg("Sidebar filter changed")
			a.currentFilter = special
			a.currentFolder = nil
			return a, a.reloadNotes
		}

		folder := a.sidebar.SelectedFolder()
		if folder != nil && (a.currentFolder == nil || folder.ID != a.currentFolder.ID) {
			logging.Debug().Int64("folder_id", folder.ID).Str("folder_name", folder.Name).Msg("Folder selected")
			a.currentFolder = folder
			a.currentFilter = ""
			return a, a.reloadNotes
		}

		// Handle enter to switch to note list
		if key.Matches(msg, keys.DefaultKeyMap.Enter) {
			logging.Debug().Msg("Enter pressed on sidebar, switching to note list")
			a.switchPanel(1)
		}

	case PanelNoteList:
		a.noteList, _ = a.noteList.Update(msg)
		a.updatePreview()

		// Handle actions on selected note
		note := a.noteList.SelectedNote()
		if note == nil {
			logging.Debug().Msg("No note selected")
			return a, nil
		}

		logging.Debug().Int64("note_id", note.ID).Str("note_title", note.Title).Msg("Note selected")

		if key.Matches(msg, keys.DefaultKeyMap.Enter) {
			logging.Info().Int64("note_id", note.ID).Msg("Opening note to edit (Enter)")
			return a.editNote(note)
		}

		if key.Matches(msg, keys.DefaultKeyMap.Edit) {
			logging.Info().Int64("note_id", note.ID).Msg("Opening note to edit (e key)")
			return a.editNote(note)
		}

		if key.Matches(msg, keys.DefaultKeyMap.Delete) {
			logging.Debug().Int64("note_id", note.ID).Msg("Showing delete confirmation")
			return a.showDeleteConfirm(note)
		}

		if key.Matches(msg, keys.DefaultKeyMap.ToggleStar) {
			logging.Debug().Int64("note_id", note.ID).Msg("Toggling star")
			return a.toggleStar(note)
		}

		if key.Matches(msg, keys.DefaultKeyMap.ToggleDone) && note.IsTodo {
			logging.Debug().Int64("note_id", note.ID).Msg("Toggling todo done")
			return a.toggleTodo(note)
		}

		if key.Matches(msg, keys.DefaultKeyMap.CyclePriority) && note.IsTodo {
			logging.Debug().Int64("note_id", note.ID).Int("priority", note.Priority).Msg("Cycling priority")
			return a.cyclePriority(note)
		}
	}

	return a, nil
}

func (a *App) handleSearchInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, keys.DefaultKeyMap.Escape) {
		a.searchMode = false
		a.searchBar.Blur()
		a.searchBar.Clear()
		return a, a.reloadNotes
	}

	if key.Matches(msg, keys.DefaultKeyMap.Enter) {
		query := a.searchBar.Value()
		if query != "" {
			return a, a.performSearch(query)
		}
		return a, nil
	}

	a.searchBar, _ = a.searchBar.Update(msg)
	return a, nil
}

func (a *App) handleDialogInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	a.dialog, _ = a.dialog.Update(msg)

	if !a.dialog.IsVisible() {
		if a.dialog.IsConfirmed() {
			switch a.dialogType {
			case "new_note":
				return a.createNote(a.dialog.InputValue(), false)
			case "new_todo":
				return a.createNote(a.dialog.InputValue(), true)
			case "delete":
				return a.deleteNote()
			}
		}
		a.showDialog = false
	}

	return a, nil
}

func (a *App) switchPanel(delta int) {
	panels := []Panel{PanelSidebar, PanelNoteList}
	currentIdx := 0
	for i, p := range panels {
		if p == a.currentPanel {
			currentIdx = i
			break
		}
	}

	newIdx := (currentIdx + delta + len(panels)) % len(panels)
	a.currentPanel = panels[newIdx]

	a.sidebar.SetFocused(a.currentPanel == PanelSidebar)
	a.noteList.SetFocused(a.currentPanel == PanelNoteList)
}

func (a *App) updateLayout() {
	sidebarWidth := a.width * 25 / 100
	if sidebarWidth < 20 {
		sidebarWidth = 20
	}
	if sidebarWidth > 40 {
		sidebarWidth = 40
	}

	noteListWidth := a.width - sidebarWidth
	contentHeight := a.height - 3 // Status bar

	a.sidebar.SetSize(sidebarWidth, contentHeight)
	a.noteList.SetSize(noteListWidth, contentHeight)
	a.preview.SetSize(noteListWidth, contentHeight/2)
	a.statusBar.SetWidth(a.width)
	a.searchBar.SetSize(a.width, 3)
	a.help.SetSize(a.width, a.height)
	a.dialog.SetSize(a.width, a.height)
}

func (a *App) updatePreview() {
	note := a.noteList.SelectedNote()
	a.preview.SetNote(note)
	a.currentNote = note
}

// View renders the UI
func (a *App) View() string {
	if !a.ready {
		return "Loading..."
	}

	// Help overlay
	if a.showHelp {
		return a.renderWithOverlay(a.help.View())
	}

	// Dialog overlay
	if a.showDialog {
		return a.renderWithOverlay(a.dialog.View())
	}

	// Main layout
	var content string

	// Search bar if active
	if a.searchMode {
		content = a.searchBar.View() + "\n"
	}

	// Header
	header := a.renderHeader()

	// Main content
	sidebar := a.sidebar.View()
	noteList := a.noteList.View()

	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, noteList)

	// Status bar
	statusBar := a.statusBar.View()

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		content+mainContent,
		statusBar,
	)
}

func (a *App) renderHeader() string {
	title := styles.TitleStyle.Render("記録 Kiroku")
	date := styles.DateStyle.Render(time.Now().Format("Mon, Jan 2 15:04"))

	spacing := a.width - lipgloss.Width(title) - lipgloss.Width(date) - 2
	if spacing < 0 {
		spacing = 0
	}

	return styles.HeaderStyle.Width(a.width).Render(
		title + strings.Repeat(" ", spacing) + date,
	)
}

func (a *App) renderWithOverlay(overlay string) string {
	return lipgloss.Place(a.width, a.height, lipgloss.Center, lipgloss.Center, overlay)
}

// Action methods

func (a *App) reloadNotes() tea.Msg {
	ctx := context.Background()
	var notes []*models.Note
	var err error

	switch a.currentFilter {
	case "all":
		notes, err = a.noteService.GetAllNotes(ctx)
		a.noteList.SetFolderName("All Notes")
	case "todos":
		notes, err = a.noteService.GetTodos(ctx, a.cfg.Todos.ShowCompleted)
		a.noteList.SetFolderName("Todos")
	case "starred":
		notes, err = a.noteService.GetStarred(ctx)
		a.noteList.SetFolderName("Starred")
	default:
		if a.currentFolder != nil {
			notes, err = a.noteService.GetByFolder(ctx, a.currentFolder.ID)
			a.noteList.SetFolderName(a.currentFolder.Name)
		} else {
			notes, err = a.noteService.GetAllNotes(ctx)
			a.noteList.SetFolderName("Notes")
		}
	}

	if err != nil {
		return errMsg{err}
	}

	return dataLoadedMsg{
		folders:   a.folders,
		notes:     notes,
		templates: a.templates,
	}
}

func (a *App) performSearch(query string) func() tea.Msg {
	return func() tea.Msg {
		ctx := context.Background()
		results, err := a.searchService.Search(ctx, query, models.ListOptions{})
		if err != nil {
			return errMsg{err}
		}

		notes := make([]*models.Note, len(results))
		for i, r := range results {
			note := r.Note
			notes[i] = &note
		}

		a.noteList.SetFolderName(fmt.Sprintf("Search: %s", query))
		return dataLoadedMsg{
			folders:   a.folders,
			notes:     notes,
			templates: a.templates,
		}
	}
}

func (a *App) showNewNoteDialog() (tea.Model, tea.Cmd) {
	a.dialog.ShowInput("New Note", "Enter note title...")
	a.dialogType = "new_note"
	a.showDialog = true
	return a, nil
}

func (a *App) showNewTodoDialog() (tea.Model, tea.Cmd) {
	a.dialog.ShowInput("New Todo", "Enter todo title...")
	a.dialogType = "new_todo"
	a.showDialog = true
	return a, nil
}

func (a *App) showDeleteConfirm(note *models.Note) (tea.Model, tea.Cmd) {
	a.dialog.ShowConfirm("Delete Note", fmt.Sprintf("Delete '%s'?", note.Title))
	a.dialogType = "delete"
	a.showDialog = true
	return a, nil
}

func (a *App) createNote(title string, isTodo bool) (tea.Model, tea.Cmd) {
	ctx := context.Background()

	var folderID *int64
	if a.currentFolder != nil {
		folderID = &a.currentFolder.ID
	}

	note := &models.Note{
		Title:    title,
		FolderID: folderID,
		IsTodo:   isTodo,
	}

	if err := a.noteService.Create(ctx, note); err != nil {
		a.statusBar.SetMessage(fmt.Sprintf("Error: %v", err))
		return a, a.clearStatusAfter(3 * time.Second)
	}

	a.showDialog = false

	a.statusBar.SetMessage(fmt.Sprintf("Created: %s", title))
	return a, tea.Batch(a.reloadNotes, a.clearStatusAfter(2*time.Second))
}

func (a *App) deleteNote() (tea.Model, tea.Cmd) {
	ctx := context.Background()

	if a.currentNote == nil {
		return a, nil
	}

	if err := a.noteService.Delete(ctx, a.currentNote.ID); err != nil {
		a.statusBar.SetMessage(fmt.Sprintf("Error: %v", err))
		return a, a.clearStatusAfter(3 * time.Second)
	}

	a.showDialog = false

	a.statusBar.SetMessage("Note deleted")
	return a, tea.Batch(a.reloadNotes, a.clearStatusAfter(2*time.Second))
}

func (a *App) editNote(note *models.Note) (tea.Model, tea.Cmd) {
	logging.Info().Int64("note_id", note.ID).Str("title", note.Title).Msg("Opening editor for note")

	// Store the note we're editing
	a.currentNote = note

	// Create temp file and get editor command
	tmpFile, editorCmd, err := a.editorService.PrepareEdit(note.Title, note.Content)
	if err != nil {
		logging.Error().Err(err).Msg("Failed to prepare editor")
		a.statusBar.SetMessage(fmt.Sprintf("Error: %v", err))
		return a, a.clearStatusAfter(3 * time.Second)
	}

	// Store temp file path for later
	a.editingTempFile = tmpFile

	// Use tea.ExecProcess to properly release terminal to editor
	return a, tea.ExecProcess(editorCmd, func(err error) tea.Msg {
		if err != nil {
			logging.Error().Err(err).Msg("Editor process failed")
			return errMsg{fmt.Errorf("editor failed: %w", err)}
		}
		logging.Debug().Msg("Editor closed, processing changes")
		return editorFinishedMsg{}
	})
}

func (a *App) toggleStar(note *models.Note) (tea.Model, tea.Cmd) {
	ctx := context.Background()
	if err := a.noteService.ToggleStar(ctx, note.ID); err != nil {
		a.statusBar.SetMessage(fmt.Sprintf("Error: %v", err))
		return a, a.clearStatusAfter(3 * time.Second)
	}

	return a, a.reloadNotes
}

func (a *App) toggleTodo(note *models.Note) (tea.Model, tea.Cmd) {
	ctx := context.Background()
	if err := a.noteService.ToggleTodo(ctx, note.ID); err != nil {
		a.statusBar.SetMessage(fmt.Sprintf("Error: %v", err))
		return a, a.clearStatusAfter(3 * time.Second)
	}

	return a, a.reloadNotes
}

func (a *App) cyclePriority(note *models.Note) (tea.Model, tea.Cmd) {
	ctx := context.Background()
	newPriority := (note.Priority + 1) % 4
	if err := a.noteService.SetPriority(ctx, note.ID, newPriority); err != nil {
		a.statusBar.SetMessage(fmt.Sprintf("Error: %v", err))
		return a, a.clearStatusAfter(3 * time.Second)
	}

	return a, a.reloadNotes
}

func (a *App) clearStatusAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(time.Time) tea.Msg {
		return statusClearMsg{}
	})
}
