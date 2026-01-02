package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Editor   EditorConfig   `mapstructure:"editor"`
	UI       UIConfig       `mapstructure:"ui"`
	Todos    TodoConfig     `mapstructure:"todos"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Path string `mapstructure:"path"`
}

// EditorConfig represents editor configuration
type EditorConfig struct {
	Command string   `mapstructure:"command"`
	Args    []string `mapstructure:"args"`
}

// UIConfig represents UI configuration
type UIConfig struct {
	Theme        string `mapstructure:"theme"`
	SidebarWidth int    `mapstructure:"sidebar_width"`
	DateFormat   string `mapstructure:"date_format"`
}

// TodoConfig represents todo configuration
type TodoConfig struct {
	ShowCompleted bool `mapstructure:"show_completed"`
	SortByDue     bool `mapstructure:"sort_by_due"`
}

// Default paths
func getDefaultPaths() (configDir, dataDir string) {
	homeDir, _ := os.UserHomeDir()
	configDir = filepath.Join(homeDir, ".config", "kiroku")
	dataDir = filepath.Join(homeDir, ".local", "share", "kiroku")
	return
}

// Load loads the configuration from file or creates default
func Load() (*Config, error) {
	configDir, dataDir := getDefaultPaths()

	// Ensure directories exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	// Set defaults
	viper.SetDefault("database.path", filepath.Join(dataDir, "kiroku.db"))
	viper.SetDefault("editor.command", getDefaultEditor())
	viper.SetDefault("editor.args", []string{})
	viper.SetDefault("ui.theme", "default")
	viper.SetDefault("ui.sidebar_width", 25)
	viper.SetDefault("ui.date_format", "2006-01-02 15:04")
	viper.SetDefault("todos.show_completed", true)
	viper.SetDefault("todos.sort_by_due", true)

	// Set config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	// Try to read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Create default config file
			configPath := filepath.Join(configDir, "config.yaml")
			if err := viper.SafeWriteConfigAs(configPath); err != nil {
				// Ignore if file already exists
				if _, statErr := os.Stat(configPath); os.IsNotExist(statErr) {
					return nil, err
				}
			}
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// getDefaultEditor returns the default editor command
func getDefaultEditor() string {
	// Check common environment variables
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	if editor := os.Getenv("VISUAL"); editor != "" {
		return editor
	}

	// Check for common editors
	editors := []string{"nvim", "vim", "vi", "nano"}
	for _, editor := range editors {
		if path, err := lookPath(editor); err == nil && path != "" {
			return editor
		}
	}

	return "nano" // Fallback
}

// lookPath checks if a command exists
func lookPath(cmd string) (string, error) {
	path := os.Getenv("PATH")
	for _, dir := range filepath.SplitList(path) {
		fullPath := filepath.Join(dir, cmd)
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath, nil
		}
	}
	return "", os.ErrNotExist
}

// Save saves the configuration to file
func (c *Config) Save() error {
	configDir, _ := getDefaultPaths()
	configPath := filepath.Join(configDir, "config.yaml")

	viper.Set("database.path", c.Database.Path)
	viper.Set("editor.command", c.Editor.Command)
	viper.Set("editor.args", c.Editor.Args)
	viper.Set("ui.theme", c.UI.Theme)
	viper.Set("ui.sidebar_width", c.UI.SidebarWidth)
	viper.Set("ui.date_format", c.UI.DateFormat)
	viper.Set("todos.show_completed", c.Todos.ShowCompleted)
	viper.Set("todos.sort_by_due", c.Todos.SortByDue)

	return viper.WriteConfigAs(configPath)
}

// GetDataDir returns the data directory path
func GetDataDir() string {
	_, dataDir := getDefaultPaths()
	return dataDir
}

// GetConfigDir returns the config directory path
func GetConfigDir() string {
	configDir, _ := getDefaultPaths()
	return configDir
}
