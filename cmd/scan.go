package cmd

import (
	"fmt"
	"time"

	"github.com/bansikah22/kswp/internal/kubernetes"
	"github.com/bansikah22/kswp/internal/report"
	"github.com/bansikah22/kswp/internal/scanner"
	"github.com/bansikah22/kswp/pkg/models"
	"github.com/bansikah22/kswp/test/mocks"
	"github.com/spf13/cobra"
)

var dryRun bool
var namespace string

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan for unused resources",
	Long:  `Scan for unused resources in your Kubernetes cluster.`,
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

		resources, err := ScanResources(client, namespace)
		if err != nil {
			fmt.Println("Error scanning resources:", err)
			return
		}

		report.PrintReport(resources)
	},
}

func ScanResources(client kubernetes.Client, namespace string) ([]models.Resource, error) {
	fmt.Println("Scanning for unused resources...")
	var resources []models.Resource

	unusedConfigMaps, err := scanner.GetUnusedConfigMaps(client.Clientset(), namespace)
	if err != nil {
		return nil, fmt.Errorf("error getting unused configmaps: %w", err)
	}
	resources = append(resources, unusedConfigMaps...)

	unusedSecrets, err := scanner.GetUnusedSecrets(client.Clientset(), namespace)
	if err != nil {
		return nil, fmt.Errorf("error getting unused secrets: %w", err)
	}
	resources = append(resources, unusedSecrets...)

	orphanServices, err := scanner.GetOrphanServices(client.Clientset(), namespace)
	if err != nil {
		return nil, fmt.Errorf("error getting orphan services: %w", err)
	}
	resources = append(resources, orphanServices...)

	oldReplicaSets, err := scanner.GetOldReplicaSets(client.Clientset(), namespace)
	if err != nil {
		return nil, fmt.Errorf("error getting old replicasets: %w", err)
	}
	resources = append(resources, oldReplicaSets...)

	completedJobs, err := scanner.GetCompletedJobs(client.Clientset(), 24*time.Hour, namespace)
	if err != nil {
		return nil, fmt.Errorf("error getting completed jobs: %w", err)
	}
	resources = append(resources, completedJobs...)

	return resources, nil
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().BoolVar(&dryRun, "dry-run", false, "run in dry-run mode")
	scanCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "specify the namespace to scan")
}
