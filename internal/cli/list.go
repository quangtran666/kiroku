package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tranducquang/kiroku/internal/models"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List notes",
	Long: `List notes with various filters.

Examples:
  kiroku list
  kiroku list --todos
  kiroku list --starred
  kiroku list --folder work`,
	RunE: runList,
}

var (
	listTodos   bool
	listStarred bool
	listFolder  string
	listLimit   int
)

func init() {
	listCmd.Flags().BoolVarP(&listTodos, "todos", "t", false, "list only todos")
	listCmd.Flags().BoolVarP(&listStarred, "starred", "s", false, "list only starred")
	listCmd.Flags().StringVarP(&listFolder, "folder", "f", "", "filter by folder")
	listCmd.Flags().IntVarP(&listLimit, "limit", "n", 20, "max number of items")
}

func runList(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	var notes []*models.Note
	var err error

	switch {
	case listTodos:
		notes, err = appInst.NoteService.GetTodos(ctx, true)
	case listStarred:
		notes, err = appInst.NoteService.GetStarred(ctx)
	default:
		notes, err = appInst.NoteService.GetAllNotes(ctx)
	}

	if err != nil {
		return fmt.Errorf("failed to list notes: %w", err)
	}

	if len(notes) == 0 {
		fmt.Println("No notes found.")
		return nil
	}

	for _, note := range notes {
		status := "📝"
		if note.IsTodo {
			if note.IsDone {
				status = "☑"
			} else {
				status = "☐"
			}
		}
		if note.Starred {
			status = "⭐"
		}

		fmt.Printf("%s [%d] %s\n", status, note.ID, note.Title)
	}

	return nil
}
