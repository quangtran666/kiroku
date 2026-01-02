package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var templatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "Manage templates",
	Long:  `List and manage note templates.`,
	RunE:  runTemplates,
}

func runTemplates(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	templates, err := appInst.TemplateService.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list templates: %w", err)
	}

	if len(templates) == 0 {
		fmt.Println("No templates found.")
		return nil
	}

	fmt.Println("Available templates:")
	for _, t := range templates {
		defaultMark := ""
		if t.IsDefault {
			defaultMark = " (default)"
		}
		fmt.Printf("  %s [%d] %s%s\n", t.Icon, t.ID, t.Name, defaultMark)
		if t.Description != "" {
			fmt.Printf("      %s\n", t.Description)
		}
	}

	return nil
}
