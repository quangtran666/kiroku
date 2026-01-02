package app

import (
	"github.com/tranducquang/kiroku/internal/config"
	"github.com/tranducquang/kiroku/internal/database"
	"github.com/tranducquang/kiroku/internal/repository"
	"github.com/tranducquang/kiroku/internal/service"
)

// App represents the application with all its dependencies
type App struct {
	Config *config.Config
	DB     *database.DB

	// Repositories
	NoteRepo     *repository.NoteRepository
	FolderRepo   *repository.FolderRepository
	TemplateRepo *repository.TemplateRepository
	SearchRepo   *repository.SearchRepository

	// Services
	NoteService     *service.NoteService
	FolderService   *service.FolderService
	TemplateService *service.TemplateService
	SearchService   *service.SearchService
	EditorService   *service.EditorService
}

// New creates a new application instance
func New(cfg *config.Config) (*App, error) {
	// Initialize database
	db, err := database.New(cfg.Database.Path)
	if err != nil {
		return nil, err
	}

	// Run migrations
	if err := db.Migrate(); err != nil {
		return nil, err
	}

	// Initialize repositories
	noteRepo := repository.NewNoteRepository(db)
	folderRepo := repository.NewFolderRepository(db)
	templateRepo := repository.NewTemplateRepository(db)
	searchRepo := repository.NewSearchRepository(db)

	// Initialize services
	noteService := service.NewNoteService(noteRepo, templateRepo, folderRepo)
	folderService := service.NewFolderService(folderRepo, noteRepo)
	templateService := service.NewTemplateService(templateRepo)
	searchService := service.NewSearchService(searchRepo)
	editorService := service.NewEditorService(cfg)

	return &App{
		Config:          cfg,
		DB:              db,
		NoteRepo:        noteRepo,
		FolderRepo:      folderRepo,
		TemplateRepo:    templateRepo,
		SearchRepo:      searchRepo,
		NoteService:     noteService,
		FolderService:   folderService,
		TemplateService: templateService,
		SearchService:   searchService,
		EditorService:   editorService,
	}, nil
}

// Close closes all resources
func (a *App) Close() error {
	if a.DB != nil {
		return a.DB.Close()
	}
	return nil
}
