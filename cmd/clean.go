package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/bansikah22/kswp/internal/cleaner"
	"github.com/bansikah22/kswp/internal/kubernetes"
	"github.com/bansikah22/kswp/test/mocks"
	"github.com/spf13/cobra"
)

var cleanDryRun bool
var cleanNamespace string

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean unused resources",
	Long:  `Clean unused resources in your Kubernetes cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		var client kubernetes.Client
		var err error
		if cleanDryRun {
			client = mocks.NewMockClient()
		} else {
			client, err = kubernetes.NewClient()
			if err != nil {
				fmt.Println("Error connecting to Kubernetes cluster:", err)
				return
			}
		}

		resources, err := ScanResources(client, cleanNamespace)
		if err != nil {
			fmt.Println("Error scanning resources:", err)
			return
		}

		if len(resources) == 0 {
			fmt.Println("No unused resources found to clean.")
			return
		}

		if cleanDryRun {
			fmt.Println("Resources that would be deleted:")
			for _, res := range resources {
				fmt.Printf("- %s/%s (%s)\n", res.Namespace, res.Name, res.Kind)
			}
			return
		}

		fmt.Printf("Delete %d resources? (y/N) ", len(resources))
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		if strings.ToLower(strings.TrimSpace(input)) != "y" {
			fmt.Println("Cleanup aborted.")
			return
		}

		for _, res := range resources {
			err := cleaner.DeleteResource(client.Clientset(), res)
			if err != nil {
				fmt.Printf("Error deleting %s %s/%s: %s\n", res.Kind, res.Namespace, res.Name, err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().BoolVar(&cleanDryRun, "dry-run", false, "run in dry-run mode")
	cleanCmd.Flags().StringVarP(&cleanNamespace, "namespace", "n", "", "specify the namespace to clean")
}
