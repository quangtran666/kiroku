package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit [id]",
	Short: "Edit a note",
	Long: `Open a note in your configured editor.

Examples:
  kiroku edit 1
  kiroku edit 42`,
	Args: cobra.ExactArgs(1),
	RunE: runEdit,
}

func runEdit(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid note ID: %w", err)
	}

	note, err := appInst.NoteService.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("note not found: %w", err)
	}

	newTitle, newContent, err := appInst.EditorService.EditNote(note.Title, note.Content)
	if err != nil {
		return fmt.Errorf("editor error: %w", err)
	}

	note.Title = newTitle
	note.Content = newContent

	if err := appInst.NoteService.Update(ctx, note); err != nil {
		return fmt.Errorf("failed to update note: %w", err)
	}

	fmt.Printf("✨ Updated note: %s\n", newTitle)
	return nil
}
