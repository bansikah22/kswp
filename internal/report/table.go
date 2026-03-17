package report

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/bansikah22/kswp/pkg/models"
	"github.com/jedib0t/go-pretty/v6/table"
)

// TableConfig holds configuration for table output
type TableConfig struct {
	ShowReason bool
	Writer     io.Writer
}

// PrintResourcesTable displays resources in a structured table format with color-coded headers
func PrintResourcesTable(resources []models.Resource, config *TableConfig) {
	if config == nil {
		config = &TableConfig{
			ShowReason: true,
			Writer:     os.Stdout,
		}
	}

	if config.Writer == nil {
		config.Writer = os.Stdout
	}

	if len(resources) == 0 {
		fmt.Fprintln(config.Writer, "No unused resources found!")
		return
	}

	t := table.NewWriter()
	t.SetStyle(table.StyleColoredBright)

	// Build header row
	header := table.Row{
		"NAMESPACE",
		"NAME",
		"RESOURCE TYPE",
		"AGE",
	}

	if config.ShowReason {
		header = append(header, "REASON")
	}

	t.AppendHeader(header)

	// Add resource rows
	for _, res := range resources {
		row := table.Row{
			res.Namespace,
			res.Name,
			res.Kind,
			formatAge(res.Age),
		}

		if config.ShowReason {
			row = append(row, res.Reason)
		}

		t.AppendRow(row)
	}

	fmt.Fprintln(config.Writer, "\nCluster Scan Report")
	fmt.Fprintln(config.Writer, "-------------------")

	// Render the table and write to the configured writer
	renderedTable := t.Render()
	fmt.Fprint(config.Writer, renderedTable)
	fmt.Fprintln(config.Writer)
}

// formatAge converts a duration to a human-readable age string
func formatAge(d time.Duration) string {
	d = d.Round(time.Second)
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh", days, hours)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

// PrintResourceTableCompact displays resources in a compact table without reason column
func PrintResourceTableCompact(resources []models.Resource) {
	PrintResourcesTable(resources, &TableConfig{
		ShowReason: false,
		Writer:     os.Stdout,
	})
}
