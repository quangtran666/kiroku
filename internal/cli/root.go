package cli

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/tranducquang/kiroku/internal/app"
	"github.com/tranducquang/kiroku/internal/config"
	"github.com/tranducquang/kiroku/internal/logging"
	"github.com/tranducquang/kiroku/internal/tui"
)

var (
	cfgFile  string
	appInst  *app.App
	logLevel string
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "kiroku",
	Short: "A terminal-based note-taking application",
	Long: `記録 Kiroku - A beautiful terminal-based note-taking application.

Kiroku helps you capture and organize your thoughts, notes, and todos
with a clean terminal user interface.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Initialize logging first
		logCfg := logging.DefaultConfig()
		if logLevel != "" {
			logCfg.Level = logLevel
		}
		if err := logging.Init(logCfg); err != nil {
			return fmt.Errorf("failed to initialize logging: %w", err)
		}

		// Skip initialization for help and version commands
		if cmd.Name() == "help" || cmd.Name() == "version" {
			return nil
		}

		logging.Info().Str("command", cmd.Name()).Msg("Starting Kiroku")

		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			logging.Error().Err(err).Msg("Failed to load config")
			return fmt.Errorf("failed to load config: %w", err)
		}

		logging.Debug().Str("db_path", cfg.Database.Path).Msg("Config loaded")

		// Initialize application
		appInst, err = app.New(cfg)
		if err != nil {
			logging.Error().Err(err).Msg("Failed to initialize application")
			return fmt.Errorf("failed to initialize application: %w", err)
		}

		logging.Info().Msg("Application initialized successfully")
		return nil
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if appInst != nil {
			appInst.Close()
		}
		logging.Close()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Launch TUI
		tuiApp := tui.NewApp(
			appInst.NoteService,
			appInst.FolderService,
			appInst.TemplateService,
			appInst.SearchService,
			appInst.EditorService,
			appInst.Config,
		)

		p := tea.NewProgram(tuiApp, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("failed to run TUI: %w", err)
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	// Recover from panics and log them
	defer logging.RecoverPanic()

	if err := rootCmd.Execute(); err != nil {
		logging.Error().Err(err).Msg("Command execution failed")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/kiroku/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "log level (debug, info, warn, error)")

	// Add subcommands
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(todoCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(templatesCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(logsCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Kiroku v0.1.0")
	},
}
