package tui

import (
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
	"github.com/tranducquang/kiroku/internal/tui/commands"
	"github.com/tranducquang/kiroku/internal/tui/components"
	"github.com/tranducquang/kiroku/internal/tui/constants"
	"github.com/tranducquang/kiroku/internal/tui/keys"
	"github.com/tranducquang/kiroku/internal/tui/messages"
	"github.com/tranducquang/kiroku/internal/tui/styles"
)

// ViewType represents the current view.
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

// Panel represents the focused panel.
type Panel int

const (
	PanelSidebar Panel = iota
	PanelNoteList
	PanelPreview
)

// App represents the main TUI application.
type App struct {
	// Services (injected via interfaces)
	noteService     service.NoteServiceInterface
	folderService   service.FolderServiceInterface
	templateService service.TemplateServiceInterface
	searchService   service.SearchServiceInterface
	editorService   service.EditorServiceInterface
	cfg             *config.Config

	// Dimensions
	width  int
	height int
	ready  bool

	// State
	currentView   ViewType
	currentPanel  Panel
	currentFilter string
	currentNote   *models.Note
	currentFolder *models.Folder

	// Data
	folders   []*models.Folder
	notes     []*models.Note
	templates []models.Template

	// Components
	sidebar   *components.Sidebar
	noteList  *components.NoteList
	preview   *components.Preview
	statusBar *components.StatusBar
	searchBar *components.SearchBar
	help      *components.Help
	dialog    *components.Dialog

	// UI State
	showHelp        bool
	showDialog      bool
	showPreview     bool
	dialogType      string
	searchMode      bool
	searchQuery     string
	editingTempFile string
}

// NewApp creates a new TUI application with the given services.
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
		currentFilter:   constants.FilterStarred,
		sidebar:         components.NewSidebar(),
		noteList:        noteList,
		preview:         components.NewPreview(),
		statusBar:       components.NewStatusBar(),
		searchBar:       components.NewSearchBar(),
		help:            components.NewHelp(),
		dialog:          components.NewDialog(),
		showPreview:     true,
	}
}

// Init initializes the application.
func (a *App) Init() tea.Cmd {
	return tea.Batch(
		a.loadData(),
		tea.SetWindowTitle("記録 Kiroku"),
	)
}

// loadData returns a command that loads initial data.
func (a *App) loadData() tea.Cmd {
	return commands.LoadData(commands.LoadDataParams{
		FolderService:   a.folderService,
		NoteService:     a.noteService,
		TemplateService: a.templateService,
	})
}

// Update handles messages.
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return a.handleWindowResize(msg)
	case messages.DataLoadedMsg:
		return a.handleDataLoaded(msg)
	case messages.ErrorMsg:
		return a.handleError(msg)
	case messages.StatusClearMsg:
		return a.handleStatusClear()
	case messages.EditorFinishedMsg:
		return a.handleEditorFinished(msg)
	case messages.NoteCreatedMsg:
		return a.handleNoteCreated(msg)
	case messages.NoteDeletedMsg:
		return a.handleNoteDeleted(msg)
	case messages.NoteUpdatedMsg:
		return a.handleNoteUpdated()
	case messages.SearchResultsMsg:
		return a.handleSearchResults(msg)
	case tea.KeyMsg:
		return a.handleKeyPress(msg)
	}

	return a, nil
}

// handleWindowResize handles window resize events.
func (a *App) handleWindowResize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	logging.Debug().Int("width", msg.Width).Int("height", msg.Height).Msg("Window resize")
	a.width = msg.Width
	a.height = msg.Height
	a.ready = true
	a.updateLayout()
	return a, nil
}

// handleDataLoaded handles data loaded events.
func (a *App) handleDataLoaded(msg messages.DataLoadedMsg) (tea.Model, tea.Cmd) {
	logging.Debug().
		Int("folders", len(msg.Folders)).
		Int("notes", len(msg.Notes)).
		Int("templates", len(msg.Templates)).
		Msg("Data loaded")

	if msg.Folders != nil {
		a.folders = msg.Folders
		a.sidebar.SetFolders(a.folders)
	}

	// Always update notes, treating nil as empty list
	a.notes = msg.Notes

	// If we are in Starred filter but got all notes (e.g. from LoadData), filter them
	if a.currentFilter == constants.FilterStarred {
		var starred []*models.Note
		for _, n := range a.notes {
			if n.Starred {
				starred = append(starred, n)
			}
		}
		a.notes = starred
	}

	a.noteList.SetNotes(a.notes)

	if msg.Templates != nil {
		a.templates = msg.Templates
	}

	// Only show starred folders in the note list if we are in the Starred filter
	if a.currentFilter == constants.FilterStarred {
		// Always update folders in note list, treating nil as empty
		a.noteList.SetFolders(msg.StarredFolders)
	} else {
		a.noteList.SetFolders(nil)
	}

	a.updatePreview()
	return a, nil
}

// handleError handles error events.
func (a *App) handleError(msg messages.ErrorMsg) (tea.Model, tea.Cmd) {
	logging.Error().Err(msg.Err).Str("context", msg.Context).Msg("TUI error")
	a.statusBar.SetMessage(fmt.Sprintf("Error: %v", msg))
	return a, commands.ClearStatusAfter(constants.ErrorMessageDuration)
}

// handleStatusClear handles status clear events.
func (a *App) handleStatusClear() (tea.Model, tea.Cmd) {
	a.statusBar.ClearMessage()
	return a, nil
}

// handleEditorFinished handles editor finished events.
func (a *App) handleEditorFinished(msg messages.EditorFinishedMsg) (tea.Model, tea.Cmd) {
	logging.Debug().Msg("Processing editor finished message")

	// Check if editor process returned an error
	if msg.Err != nil {
		logging.Error().Err(msg.Err).Msg("Editor process failed")
		a.statusBar.SetMessage(fmt.Sprintf("Editor failed: %v", msg.Err))
		a.editingTempFile = ""
		return a, commands.ClearStatusAfter(constants.ErrorMessageDuration)
	}

	if a.currentNote == nil || a.editingTempFile == "" {
		logging.Warn().Msg("No note or temp file to process")
		return a, nil
	}

	newTitle, newContent, err := a.editorService.ReadEditedContent(a.editingTempFile, a.currentNote.Title)
	if err != nil {
		logging.Error().Err(err).Msg("Failed to read edited content")
		a.statusBar.SetMessage(fmt.Sprintf("Error: %v", err))
		a.editingTempFile = ""
		return a, commands.ClearStatusAfter(constants.ErrorMessageDuration)
	}

	a.currentNote.Title = newTitle
	a.currentNote.Content = newContent
	a.editingTempFile = ""

	return a, tea.Batch(
		commands.UpdateNote(a.noteService, a.currentNote),
		a.reloadNotes(),
	)
}

// handleNoteCreated handles note created events.
func (a *App) handleNoteCreated(msg messages.NoteCreatedMsg) (tea.Model, tea.Cmd) {
	a.showDialog = false
	a.statusBar.SetMessage(fmt.Sprintf("Created: %s", msg.Note.Title))
	return a, tea.Batch(
		a.reloadNotes(),
		commands.ClearStatusAfter(constants.StatusMessageDuration),
	)
}

// handleNoteDeleted handles note deleted events.
func (a *App) handleNoteDeleted(msg messages.NoteDeletedMsg) (tea.Model, tea.Cmd) {
	a.showDialog = false
	a.statusBar.SetMessage("Note deleted")

	// Explicitly clear state to ensure immediate visual feedback
	a.currentNote = nil
	a.preview.SetNote(nil)

	// Remove from local list immediately
	var newNotes []*models.Note
	for _, n := range a.notes {
		if n.ID != msg.NoteID {
			newNotes = append(newNotes, n)
		}
	}
	a.notes = newNotes
	a.noteList.SetNotes(a.notes)

	return a, tea.Batch(
		a.reloadNotes(),
		commands.ClearStatusAfter(constants.StatusMessageDuration),
	)
}

// handleNoteUpdated handles note updated events.
func (a *App) handleNoteUpdated() (tea.Model, tea.Cmd) {
	a.statusBar.SetMessage("Note saved")
	return a, tea.Batch(
		a.reloadNotes(),
		commands.ClearStatusAfter(constants.StatusMessageDuration),
	)
}

// handleSearchResults handles search results.
func (a *App) handleSearchResults(msg messages.SearchResultsMsg) (tea.Model, tea.Cmd) {
	a.noteList.SetFolderName(fmt.Sprintf("Search: %s", msg.Query))
	a.notes = msg.Notes
	a.noteList.SetNotes(a.notes)
	a.noteList.ResetCursor()
	a.updatePreview()
	return a, nil
}

// handleKeyPress handles key press events.
func (a *App) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	logging.Debug().
		Str("key", msg.String()).
		Str("panel", panelName(a.currentPanel)).
		Bool("dialog", a.showDialog).
		Bool("search", a.searchMode).
		Bool("help", a.showHelp).
		Msg("Key pressed")

	// Handle overlays first (help, dialog, search)
	if a.showHelp {
		return a.handleHelpInput(msg)
	}
	if a.showDialog {
		return a.handleDialogInput(msg)
	}
	if a.searchMode {
		return a.handleSearchInput(msg)
	}

	// Handle global keys
	if handled, cmd := a.handleGlobalKeys(msg); handled {
		return a, cmd
	}

	// Handle panel-specific input
	return a.handlePanelInput(msg)
}

// handleGlobalKeys handles keys that work regardless of panel.
func (a *App) handleGlobalKeys(msg tea.KeyMsg) (bool, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.DefaultKeyMap.Quit):
		logging.Info().Msg("User quit application")
		return true, tea.Quit

	case key.Matches(msg, keys.DefaultKeyMap.Help):
		logging.Debug().Msg("Showing help")
		a.showHelp = true
		a.help.Show()
		return true, nil

	case key.Matches(msg, keys.DefaultKeyMap.Search):
		logging.Debug().Msg("Entering search mode")
		a.searchMode = true
		a.searchBar.Focus()
		a.updateLayout()
		return true, nil

	case key.Matches(msg, keys.DefaultKeyMap.Preview):
		logging.Debug().Bool("show_preview", !a.showPreview).Msg("Toggling preview")
		a.showPreview = !a.showPreview
		a.updateLayout()
		return true, nil

	case key.Matches(msg, keys.DefaultKeyMap.NewNote):
		logging.Debug().Msg("Showing new note dialog")
		a.showNewNoteDialog()
		return true, nil

	case key.Matches(msg, keys.DefaultKeyMap.NewTodo):
		logging.Debug().Msg("Showing new todo dialog")
		a.showNewTodoDialog()
		return true, nil

	case key.Matches(msg, keys.DefaultKeyMap.NewFolder):
		logging.Debug().Msg("Showing new folder dialog")
		a.showNewFolderDialog()
		return true, nil

	case key.Matches(msg, keys.DefaultKeyMap.Tab), key.Matches(msg, keys.DefaultKeyMap.Right):
		logging.Debug().Msg("Switching panel right")
		a.switchPanel(1)
		return true, nil

	case key.Matches(msg, keys.DefaultKeyMap.Left):
		logging.Debug().Msg("Switching panel left")
		a.switchPanel(-1)
		return true, nil

	default:
		return false, nil
	}
}

// handleHelpInput handles input when help overlay is visible.
func (a *App) handleHelpInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	a.help, _ = a.help.Update(msg)
	if !a.help.IsVisible() {
		a.showHelp = false
	}
	return a, nil
}

// handleDialogInput handles input when dialog is visible.
func (a *App) handleDialogInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	logging.Debug().Str("dialog_type", a.dialogType).Msg("Handling dialog input")
	a.dialog, _ = a.dialog.Update(msg)

	if a.dialog.IsVisible() {
		return a, nil
	}

	if !a.dialog.IsConfirmed() {
		a.showDialog = false
		return a, nil
	}

	a.showDialog = false

	switch a.dialogType {
	case constants.DialogTypeNewNote:
		return a, commands.CreateNote(commands.CreateNoteParams{
			NoteService:   a.noteService,
			Title:         a.dialog.InputValue(),
			IsTodo:        false,
			CurrentFolder: a.currentFolder,
		})

	case constants.DialogTypeNewTodo:
		return a, commands.CreateNote(commands.CreateNoteParams{
			NoteService:   a.noteService,
			Title:         a.dialog.InputValue(),
			IsTodo:        true,
			CurrentFolder: a.currentFolder,
		})

	case constants.DialogTypeDelete:
		if a.currentNote != nil {
			return a, commands.DeleteNote(a.noteService, a.currentNote.ID)
		}

	case constants.DialogTypeNewFolder:
		var parentID *int64
		if a.currentFolder != nil {
			parentID = &a.currentFolder.ID
		}
		return a, commands.CreateFolder(commands.CreateFolderParams{
			FolderService: a.folderService,
			Name:          a.dialog.InputValue(),
			ParentID:      parentID,
		})

	case constants.DialogTypeDeleteFolder:
		if a.currentFolder != nil {
			return a, commands.DeleteFolder(a.folderService, a.currentFolder.ID)
		}
	}

	return a, nil
}

// handleSearchInput handles input when in search mode.
func (a *App) handleSearchInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	logging.Debug().Msg("Handling search input")

	if key.Matches(msg, keys.DefaultKeyMap.Escape) {
		a.searchMode = false
		a.searchBar.Blur()
		a.searchBar.Clear()
		a.updateLayout()
		return a, a.reloadNotes()
	}

	// Submit search on Enter or Down
	if key.Matches(msg, keys.DefaultKeyMap.Enter) || key.Matches(msg, keys.DefaultKeyMap.Down) {
		query := a.searchBar.Value()
		if query == "" {
			return a, nil
		}

		// Exit search mode to view results
		a.searchMode = false
		a.searchBar.Blur()
		a.updateLayout()

		// Focus note list
		a.currentPanel = PanelNoteList
		a.noteList.SetFocused(true)
		a.sidebar.SetFocused(false)

		return a, commands.Search(commands.SearchParams{
			SearchService: a.searchService,
			Query:         query,
		})
	}

	a.searchBar, _ = a.searchBar.Update(msg)
	return a, nil
}

// handlePanelInput handles panel-specific input.
func (a *App) handlePanelInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	logging.Debug().Str("panel", panelName(a.currentPanel)).Msg("Handling panel input")

	switch a.currentPanel {
	case PanelSidebar:
		return a.handleSidebarInput(msg)
	case PanelNoteList:
		return a.handleNoteListInput(msg)
	}

	return a, nil
}

// handleSidebarInput handles sidebar panel input.
func (a *App) handleSidebarInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	a.sidebar, _ = a.sidebar.Update(msg)

	special := a.sidebar.SelectedSpecial()
	if special != "" && special != a.currentFilter {
		logging.Debug().Str("filter", special).Msg("Sidebar filter changed")
		a.currentFilter = special
		a.currentFolder = nil
		return a, a.reloadNotes()
	}

	folder := a.sidebar.SelectedFolder()
	if folder != nil && (a.currentFolder == nil || folder.ID != a.currentFolder.ID) {
		logging.Debug().Int64("folder_id", folder.ID).Str("folder_name", folder.Name).Msg("Folder selected")
		a.currentFolder = folder
		a.currentFilter = ""
		a.notes = nil
		a.noteList.SetNotes(nil)
		return a, a.reloadNotes()
	}

	if key.Matches(msg, keys.DefaultKeyMap.Enter) {
		logging.Debug().Msg("Enter pressed on sidebar")
		// Let Sidebar component handle it (expand/collapse).
		// No-op here means no panel switch.
	}

	if key.Matches(msg, keys.DefaultKeyMap.Delete) {
		folder := a.sidebar.SelectedFolder()
		if folder != nil {
			a.showDeleteFolderConfirm(folder)
		}
	}

	if key.Matches(msg, keys.DefaultKeyMap.ToggleStar) {
		folder := a.sidebar.SelectedFolder()
		if folder != nil {
			return a, commands.ToggleFolderStar(a.folderService, folder.ID)
		}
	}

	return a, nil
}

// handleNoteListInput handles note list panel input.
func (a *App) handleNoteListInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	a.noteList, _ = a.noteList.Update(msg)
	a.updatePreview()

	note := a.noteList.SelectedNote()
	folder := a.noteList.SelectedFolder()

	if note == nil && folder == nil {
		logging.Debug().Msg("No item selected")
		return a, nil
	}

	if folder != nil {
		logging.Debug().Int64("folder_id", folder.ID).Str("folder_name", folder.Name).Msg("Folder selected in list")
		switch {
		case key.Matches(msg, keys.DefaultKeyMap.Enter):
			// Navigate to folder
			a.currentFolder = folder
			a.currentFilter = ""
			// Clean up note list state before reload
			a.notes = nil
			a.noteList.SetNotes(nil)
			a.noteList.SetFolders(nil)

			// We also need to update the sidebar selection to reflect this change
			a.sidebar.SelectFolder(folder.ID)

			return a, a.reloadNotes()

		case key.Matches(msg, keys.DefaultKeyMap.ToggleStar):
			return a, tea.Batch(
				commands.ToggleFolderStar(a.folderService, folder.ID),
				// We need to reload to refresh the list, specifically the Starred list
				a.reloadNotes(),
			)
		}
		return a, nil
	}

	logging.Debug().Int64("note_id", note.ID).Str("note_title", note.Title).Msg("Note selected")

	switch {
	case key.Matches(msg, keys.DefaultKeyMap.Enter), key.Matches(msg, keys.DefaultKeyMap.Edit):
		logging.Info().Int64("note_id", note.ID).Msg("Opening note to edit")
		return a.editNote(note)

	case key.Matches(msg, keys.DefaultKeyMap.Delete):
		logging.Debug().Int64("note_id", note.ID).Msg("Showing delete confirmation")
		a.showDeleteConfirm(note)

	case key.Matches(msg, keys.DefaultKeyMap.ToggleStar):
		logging.Debug().Int64("note_id", note.ID).Msg("Toggling star")
		return a, commands.ToggleStar(a.noteService, note.ID)

	case key.Matches(msg, keys.DefaultKeyMap.ToggleDone) && note.IsTodo:
		logging.Debug().Int64("note_id", note.ID).Msg("Toggling todo done")
		return a, commands.ToggleTodo(a.noteService, note.ID)

	case key.Matches(msg, keys.DefaultKeyMap.CyclePriority) && note.IsTodo:
		logging.Debug().Int64("note_id", note.ID).Int("priority", note.Priority).Msg("Cycling priority")
		return a, commands.CyclePriority(a.noteService, note.ID, note.Priority)
	}

	return a, nil
}

// Panel switching and layout methods

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
	// Calculate sidebar width
	sidebarWidth := a.width * constants.SidebarWidthPercent / 100
	if sidebarWidth < constants.SidebarMinWidth {
		sidebarWidth = constants.SidebarMinWidth
	}
	if sidebarWidth > constants.SidebarMaxWidth {
		sidebarWidth = constants.SidebarMaxWidth
	}

	// Right panel takes remaining width
	rightPanelWidth := a.width - sidebarWidth
	if rightPanelWidth < 40 {
		rightPanelWidth = 40
	}

	// Total available height for main content (subtract header=1, status bar=1)
	availableHeight := a.height - 2

	// If search mode is active, subtract search bar height (3 lines)
	if a.searchMode {
		availableHeight -= 3
	}

	if availableHeight < 15 {
		availableHeight = 15
	}

	// Note list gets 50% of height, preview gets 50%
	// Add 2 to each for borders (top + bottom)
	noteListTotalHeight := (availableHeight * 50) / 100
	previewTotalHeight := availableHeight - noteListTotalHeight

	if !a.showPreview {
		noteListTotalHeight = availableHeight
		previewTotalHeight = 0
	}

	// Ensure minimum heights (including borders)
	if noteListTotalHeight < 8 {
		noteListTotalHeight = 8
	}
	if a.showPreview && previewTotalHeight < 6 {
		previewTotalHeight = 6
	}

	// Sidebar gets full available height
	a.sidebar.SetSize(sidebarWidth, availableHeight)
	a.noteList.SetSize(rightPanelWidth, noteListTotalHeight)
	if a.showPreview {
		a.preview.SetSize(rightPanelWidth, previewTotalHeight)
	} else {
		a.preview.SetSize(0, 0)
	}
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

// View renders the UI.
func (a *App) View() string {
	if !a.ready {
		return "Loading..."
	}

	if a.showHelp {
		return a.renderWithOverlay(a.help.View())
	}

	if a.showDialog {
		return a.renderWithOverlay(a.dialog.View())
	}

	header := a.renderHeader()

	sidebar := a.sidebar.View()

	// Note list and preview combined vertically on the right side
	noteListView := a.noteList.View()
	parts := []string{sidebar}

	if a.showPreview {
		previewView := a.preview.View()
		rightPanel := lipgloss.JoinVertical(lipgloss.Left, noteListView, previewView)
		parts = append(parts, rightPanel)
	} else {
		parts = append(parts, noteListView)
	}

	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, parts...)
	statusBar := a.statusBar.View()

	parts = []string{header}
	if a.searchMode {
		parts = append(parts, a.searchBar.View())
	}
	parts = append(parts, mainContent, statusBar)

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func (a *App) renderHeader() string {
	title := styles.TitleStyle.Render("記録 Kiroku")
	date := styles.DateStyle.Render(time.Now().Format("Mon, Jan 2 15:04"))

	spacing := a.width - lipgloss.Width(title) - lipgloss.Width(date) - 2
	if spacing < 0 {
		spacing = 0
	}

	return styles.HeaderStyle.Width(a.width - 2).Render(
		title + strings.Repeat(" ", spacing) + date,
	)
}

func (a *App) renderWithOverlay(overlay string) string {
	return lipgloss.Place(a.width, a.height, lipgloss.Center, lipgloss.Center, overlay)
}

// Action methods

func (a *App) reloadNotes() tea.Cmd {
	folderName := a.getFolderDisplayName()
	a.noteList.SetFolderName(folderName)

	return commands.ReloadNotes(commands.ReloadNotesParams{
		NoteService:   a.noteService,
		FolderService: a.folderService,
		CurrentFilter: a.currentFilter,
		CurrentFolder: a.currentFolder,
		ShowCompleted: a.cfg.Todos.ShowCompleted,
	})
}

func (a *App) getFolderDisplayName() string {
	switch a.currentFilter {
	case constants.FilterAll:
		return "All Notes"
	case constants.FilterTodos:
		return "Todos"
	case constants.FilterStarred:
		return "Starred"
	default:
		if a.currentFolder != nil {
			return a.currentFolder.Name
		}
		return "Notes"
	}
}

func (a *App) showNewNoteDialog() {
	a.dialog.ShowInput("New Note", "Enter note title...")
	a.dialogType = constants.DialogTypeNewNote
	a.showDialog = true
}

func (a *App) showNewTodoDialog() {
	a.dialog.ShowInput("New Todo", "Enter todo title...")
	a.dialogType = constants.DialogTypeNewTodo
	a.showDialog = true
}

func (a *App) showNewFolderDialog() {
	title := "New Folder"
	if a.currentFolder != nil {
		title = fmt.Sprintf("New Folder in '%s'", a.currentFolder.Name)
	}
	a.dialog.ShowInput(title, "Enter folder name...")
	a.dialogType = constants.DialogTypeNewFolder
	a.showDialog = true
}

func (a *App) showDeleteConfirm(note *models.Note) {
	a.dialog.ShowConfirm("Delete Note", fmt.Sprintf("Delete '%s'?", note.Title))
	a.dialogType = constants.DialogTypeDelete
	a.showDialog = true
}

func (a *App) showDeleteFolderConfirm(folder *models.Folder) {
	a.dialog.ShowConfirm("Delete Folder", fmt.Sprintf("Delete '%s'?", folder.Name))
	a.dialogType = constants.DialogTypeDeleteFolder
	a.showDialog = true
}

func (a *App) editNote(note *models.Note) (tea.Model, tea.Cmd) {
	logging.Info().Int64("note_id", note.ID).Str("title", note.Title).Msg("Opening editor for note")

	a.currentNote = note

	tmpFile, editorCmd, err := a.editorService.PrepareEdit(note.Title, note.Content)
	if err != nil {
		logging.Error().Err(err).Msg("Failed to prepare editor")
		a.statusBar.SetMessage(fmt.Sprintf("Error: %v", err))
		return a, commands.ClearStatusAfter(constants.ErrorMessageDuration)
	}

	a.editingTempFile = tmpFile

	// Use tea.ExecProcess which properly suspends and resumes the TUI
	return a, tea.ExecProcess(editorCmd, func(err error) tea.Msg {
		if err != nil {
			logging.Error().Err(err).Str("temp_file", tmpFile).Msg("Editor process failed")
			// Still try to process changes even if there was an error, UNLESS it's a launch error
		}
		logging.Debug().Msg("Editor closed, processing changes")
		return messages.EditorFinishedMsg{TempFile: tmpFile, NoteID: note.ID, Err: err}
	})
}

// Helper functions

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
