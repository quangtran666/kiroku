package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tranducquang/kiroku/internal/models"
)

var todoCmd = &cobra.Command{
	Use:   "todo [title]",
	Short: "Quick add a new todo",
	Long: `Quick add a new todo item.

Examples:
  kiroku todo "Buy groceries"
  kiroku todo "Review PR" --priority high
  kiroku todo "Call doctor" --due tomorrow`,
	Args: cobra.ExactArgs(1),
	RunE: runTodo,
}

var (
	todoPriority string
	todoDue      string
)

func init() {
	todoCmd.Flags().StringVarP(&todoPriority, "priority", "p", "", "priority (low, medium, high)")
	todoCmd.Flags().StringVarP(&todoDue, "due", "d", "", "due date")
}

func runTodo(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	title := args[0]

	priority := models.PriorityNone
	switch todoPriority {
	case "low", "l":
		priority = models.PriorityLow
	case "medium", "m":
		priority = models.PriorityMedium
	case "high", "h":
		priority = models.PriorityHigh
	}

	note := &models.Note{
		Title:    title,
		IsTodo:   true,
		Priority: priority,
	}

	if err := appInst.NoteService.Create(ctx, note); err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}

	fmt.Printf("☐ Created todo: %s\n", title)
	return nil
}
