package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tranducquang/kiroku/internal/models"
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search notes",
	Long: `Search notes using full-text search.

Examples:
  kiroku search "meeting notes"
  kiroku search "golang" --limit 10`,
	Args: cobra.ExactArgs(1),
	RunE: runSearch,
}

var searchLimit int

func init() {
	searchCmd.Flags().IntVarP(&searchLimit, "limit", "n", 20, "max number of results")
}

func runSearch(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	query := args[0]

	results, err := appInst.SearchService.Search(ctx, query, models.ListOptions{
		Limit: searchLimit,
	})
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	if len(results) == 0 {
		fmt.Println("No results found.")
		return nil
	}

	fmt.Printf("Found %d results:\n\n", len(results))
	for _, r := range results {
		fmt.Printf("📝 [%d] %s\n", r.Note.ID, r.Note.Title)
		if r.Snippet != "" {
			fmt.Printf("   %s\n", r.Snippet)
		}
		fmt.Println()
	}

	return nil
}
