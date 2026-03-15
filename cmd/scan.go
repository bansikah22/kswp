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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

		report.PrintReport(resources)
	},
}

func ScanResources(client kubernetes.Client, namespace string, label string, name string) ([]models.Resource, error) {
	fmt.Println("Scanning for unused resources...")
	var resources []models.Resource

	listOptions := metav1.ListOptions{}
	if label != "" {
		listOptions.LabelSelector = label
	}
	if name != "" {
		listOptions.FieldSelector = "metadata.name=" + name
	}

	unusedConfigMaps, err := scanner.GetUnusedConfigMaps(client.Clientset(), namespace, listOptions)
	if err != nil {
		return nil, fmt.Errorf("error getting unused configmaps: %w", err)
	}
	resources = append(resources, unusedConfigMaps...)

	unusedSecrets, err := scanner.GetUnusedSecrets(client.Clientset(), namespace, listOptions)
	if err != nil {
		return nil, fmt.Errorf("error getting unused secrets: %w", err)
	}
	resources = append(resources, unusedSecrets...)

	orphanServices, err := scanner.GetOrphanServices(client.Clientset(), namespace, listOptions)
	if err != nil {
		return nil, fmt.Errorf("error getting orphan services: %w", err)
	}
	resources = append(resources, orphanServices...)

	oldReplicaSets, err := scanner.GetOldReplicaSets(client.Clientset(), namespace, listOptions)
	if err != nil {
		return nil, fmt.Errorf("error getting old replicasets: %w", err)
	}
	resources = append(resources, oldReplicaSets...)

	completedJobs, err := scanner.GetCompletedJobs(client.Clientset(), 24*time.Hour, namespace, listOptions)
	if err != nil {
		return nil, fmt.Errorf("error getting completed jobs: %w", err)
	}
	resources = append(resources, completedJobs...)

	failedPods, err := scanner.GetFailedPods(client.Clientset(), namespace, listOptions)
	if err != nil {
		return nil, fmt.Errorf("error getting failed pods: %w", err)
	}
	resources = append(resources, failedPods...)

	completedPods, err := scanner.GetCompletedPods(client.Clientset(), 24*time.Hour, namespace, listOptions)
	if err != nil {
		return nil, fmt.Errorf("error getting completed pods: %w", err)
	}
	resources = append(resources, completedPods...)

	unusedPVCs, err := scanner.GetUnusedPersistentVolumeClaims(client.Clientset(), namespace, listOptions)
	if err != nil {
		return nil, fmt.Errorf("error getting unused pvcs: %w", err)
	}
	resources = append(resources, unusedPVCs...)

	return resources, nil
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().BoolVar(&dryRun, "dry-run", false, "run in dry-run mode")
	scanCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "specify the namespace to scan")
	scanCmd.Flags().String("label", "", "filter resources by label (e.g., 'app=nginx')")
	scanCmd.Flags().String("name", "", "filter resources by name")
}
