package cmd

import (
	"fmt"

	"github.com/bansikah22/kswp/internal/kubernetes"
	"github.com/bansikah22/kswp/internal/tui"
	"github.com/bansikah22/kswp/test/mocks"
	"github.com/spf13/cobra"
)

var tuiNamespace string

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Terminal UI for kswp",
	Long:  `A terminal user interface for kswp to visualize and manage unused resources.`,
	Run: func(cmd *cobra.Command, args []string) {
		var client kubernetes.Client
		var err error
		if dryRun {
			client = mocks.NewMockClient()
		} else {
			client, err = kubernetes.NewClient()
			if err != nil {
				fmt.Println("Error connecting to Kubernetes cluster:", err)
				return
			}
		}

		label, _ := cmd.Flags().GetString("label")
		name, _ := cmd.Flags().GetString("name")
		resources, err := ScanResources(client, tuiNamespace, label, name)
		if err != nil {
			fmt.Println("Error scanning resources:", err)
			return
		}

		tui.StartTUI(resources, client)
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
	tuiCmd.Flags().BoolVar(&dryRun, "dry-run", false, "run in dry-run mode")
	tuiCmd.Flags().StringVarP(&tuiNamespace, "namespace", "n", "", "specify the namespace to scan")
	tuiCmd.Flags().String("label", "", "filter resources by label (e.g., 'app=nginx')")
	tuiCmd.Flags().String("name", "", "filter resources by name")
}
