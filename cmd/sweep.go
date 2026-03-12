package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bansikah22/kswp/internal/cleaner"
	"github.com/bansikah22/kswp/internal/kubernetes"
	"github.com/bansikah22/kswp/pkg/models"
	"github.com/bansikah22/kswp/test/mocks"
	"github.com/spf13/cobra"
)

var olderThan string
var sweepNamespace string
var sweepDryRun bool

var sweepCmd = &cobra.Command{
	Use:   "sweep",
	Short: "Sweep unused resources",
	Long:  `Sweep unused resources in your Kubernetes cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		var client kubernetes.Client
		var err error
		if sweepDryRun {
			client = mocks.NewMockClient()
		} else {
			client, err = kubernetes.NewClient()
			if err != nil {
				fmt.Println("Error connecting to Kubernetes cluster:", err)
				return
			}
		}

		resources, err := ScanResources(client, sweepNamespace)
		if err != nil {
			fmt.Println("Error scanning resources:", err)
			return
		}

		if olderThan != "" {
			duration, err := time.ParseDuration(olderThan)
			if err != nil {
				fmt.Println("Error parsing duration:", err)
				return
			}
			resources = filterByAge(resources, duration)
		}

		if len(resources) == 0 {
			fmt.Println("No unused resources found to sweep.")
			return
		}

		fmt.Println("Resources that will be deleted:")
		for _, res := range resources {
			fmt.Printf("- %s/%s (%s)\n", res.Namespace, res.Name, res.Kind)
		}

		fmt.Printf("Delete %d resources? (y/N) ", len(resources))
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		if strings.ToLower(strings.TrimSpace(input)) != "y" {
			fmt.Println("Sweep aborted.")
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

func filterByAge(resources []models.Resource, duration time.Duration) []models.Resource {
	var filtered []models.Resource
	for _, res := range resources {
		if res.Age > duration {
			filtered = append(filtered, res)
		}
	}
	return filtered
}

func init() {
	rootCmd.AddCommand(sweepCmd)
	sweepCmd.Flags().StringVar(&olderThan, "older-than", "", "filter resources older than a duration (e.g., 7d, 24h)")
	sweepCmd.Flags().BoolVar(&sweepDryRun, "dry-run", false, "run in dry-run mode")
	sweepCmd.Flags().StringVarP(&sweepNamespace, "namespace", "n", "", "specify the namespace to sweep")
}
