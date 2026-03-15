package cmd

import (
	"fmt"

	"github.com/bansikah22/kswp/internal/analyzer"
	"github.com/bansikah22/kswp/internal/kubernetes"
	"github.com/bansikah22/kswp/test/mocks"
	"github.com/spf13/cobra"
)

var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Display a dependency graph of resources",
	Long:  `Display a dependency graph of resources in your Kubernetes cluster.`,
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

		label, err := cmd.Flags().GetString("label")
		if err != nil {
			fmt.Println("Error getting label flag:", err)
			return
		}
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			fmt.Println("Error getting name flag:", err)
			return
		}
		resources, err := ScanResources(client, namespace, label, name)
		if err != nil {
			fmt.Println("Error scanning resources:", err)
			return
		}

		graph := analyzer.BuildDependencyGraph(resources)
		PrintGraph(graph, 0)
	},
}

func PrintGraph(node *analyzer.Node, level int) {
	for i := 0; i < level; i++ {
		fmt.Print("  ")
	}
	fmt.Printf("- %s/%s (%s)\n", node.Resource.Namespace, node.Resource.Name, node.Resource.Kind)
	for _, child := range node.Children {
		PrintGraph(child, level+1)
	}
}

func init() {
	rootCmd.AddCommand(graphCmd)
	graphCmd.Flags().BoolVar(&dryRun, "dry-run", false, "run in dry-run mode")
	graphCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "specify the namespace to scan")
	graphCmd.Flags().String("label", "", "filter resources by label (e.g., 'app=nginx')")
	graphCmd.Flags().String("name", "", "filter resources by name")
}
