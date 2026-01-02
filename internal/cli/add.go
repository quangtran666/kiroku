package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tranducquang/kiroku/internal/models"
)

var addCmd = &cobra.Command{
	Use:   "add [title]",
	Short: "Quick add a new note",
	Long: `Quick add a new note with a title.

Examples:
  kiroku add "Meeting notes"
  kiroku add "Sprint planning" --folder work`,
	Args: cobra.ExactArgs(1),
	RunE: runAdd,
}

var addFolder string

func init() {
	addCmd.Flags().StringVarP(&addFolder, "folder", "f", "", "folder name or ID")
}

func runAdd(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	title := args[0]

	note := &models.Note{
		Title: title,
	}

	if err := appInst.NoteService.Create(ctx, note); err != nil {
		return fmt.Errorf("failed to create note: %w", err)
	}

	fmt.Printf("✨ Created note: %s\n", title)
	return nil
}
