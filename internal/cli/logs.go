package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"

	"github.com/spf13/cobra"

	"github.com/tranducquang/kiroku/internal/logging"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "View or manage application logs",
	Long:  `View or manage Kiroku application logs. Use subcommands to view, tail, or clear logs.`,
}

var logsShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the latest log file",
	RunE: func(cmd *cobra.Command, args []string) error {
		logDir := logging.GetLogDir()
		latestLog, err := getLatestLogFile(logDir)
		if err != nil {
			return err
		}

		content, err := os.ReadFile(latestLog)
		if err != nil {
			return fmt.Errorf("failed to read log file: %w", err)
		}

		fmt.Println(string(content))
		return nil
	},
}

var logsTailCmd = &cobra.Command{
	Use:   "tail",
	Short: "Tail the latest log file (follow mode)",
	RunE: func(cmd *cobra.Command, args []string) error {
		logDir := logging.GetLogDir()
		latestLog, err := getLatestLogFile(logDir)
		if err != nil {
			return err
		}

		fmt.Printf("📋 Tailing: %s\n\n", latestLog)

		// Use tail -f on Unix systems
		tailCmd := exec.Command("tail", "-f", latestLog)
		tailCmd.Stdout = os.Stdout
		tailCmd.Stderr = os.Stderr
		return tailCmd.Run()
	},
}

var logsPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show the logs directory path",
	Run: func(cmd *cobra.Command, args []string) {
		logDir := logging.GetLogDir()
		fmt.Printf("📁 Log directory: %s\n", logDir)

		// List log files
		entries, err := os.ReadDir(logDir)
		if err != nil {
			fmt.Printf("⚠️  Cannot read log directory: %v\n", err)
			return
		}

		fmt.Println("\n📋 Log files:")
		for _, entry := range entries {
			if !entry.IsDir() && filepath.Ext(entry.Name()) == ".log" {
				info, _ := entry.Info()
				fmt.Printf("   - %s (%d bytes)\n", entry.Name(), info.Size())
			}
		}
	},
}

var logsClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all log files",
	RunE: func(cmd *cobra.Command, args []string) error {
		logDir := logging.GetLogDir()

		entries, err := os.ReadDir(logDir)
		if err != nil {
			return fmt.Errorf("failed to read log directory: %w", err)
		}

		count := 0
		for _, entry := range entries {
			if !entry.IsDir() && filepath.Ext(entry.Name()) == ".log" {
				logPath := filepath.Join(logDir, entry.Name())
				if err := os.Remove(logPath); err != nil {
					fmt.Printf("⚠️  Failed to remove %s: %v\n", entry.Name(), err)
				} else {
					count++
				}
			}
		}

		fmt.Printf("🗑️  Cleared %d log files\n", count)
		return nil
	},
}

var logsOpenCmd = &cobra.Command{
	Use:   "open",
	Short: "Open the logs directory in file manager",
	RunE: func(cmd *cobra.Command, args []string) error {
		logDir := logging.GetLogDir()

		var openCmd *exec.Cmd
		switch runtime.GOOS {
		case "darwin":
			openCmd = exec.Command("open", logDir)
		case "linux":
			openCmd = exec.Command("xdg-open", logDir)
		case "windows":
			openCmd = exec.Command("explorer", logDir)
		default:
			return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
		}

		return openCmd.Start()
	},
}

func init() {
	logsCmd.AddCommand(logsShowCmd)
	logsCmd.AddCommand(logsTailCmd)
	logsCmd.AddCommand(logsPathCmd)
	logsCmd.AddCommand(logsClearCmd)
	logsCmd.AddCommand(logsOpenCmd)
}

// getLatestLogFile returns the path to the most recent log file
func getLatestLogFile(logDir string) (string, error) {
	entries, err := os.ReadDir(logDir)
	if err != nil {
		return "", fmt.Errorf("failed to read log directory: %w", err)
	}

	var logFiles []os.DirEntry
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".log" {
			logFiles = append(logFiles, entry)
		}
	}

	if len(logFiles) == 0 {
		return "", fmt.Errorf("no log files found in %s", logDir)
	}

	// Sort by name descending (newest first based on timestamp in filename)
	sort.Slice(logFiles, func(i, j int) bool {
		return logFiles[i].Name() > logFiles[j].Name()
	})

	return filepath.Join(logDir, logFiles[0].Name()), nil
}
