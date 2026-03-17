package cmd

import (
	"fmt"
	"strings"
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
		excludeNamespacesStr, err := cmd.Flags().GetString("exclude-namespaces")
		if err != nil {
			fmt.Println("Error getting exclude-namespaces flag:", err)
			return
		}

		var excludedNamespaces []string
		if excludeNamespacesStr != "" {
			excludedNamespaces = strings.Split(excludeNamespacesStr, ",")
			for i := range excludedNamespaces {
				excludedNamespaces[i] = strings.TrimSpace(excludedNamespaces[i])
			}
		}

		resources, err := ScanResources(client, namespace, label, name, excludedNamespaces)
		if err != nil {
			fmt.Println("Error scanning resources:", err)
			return
		}

		report.PrintReport(resources)
	},
}

func ScanResources(client kubernetes.Client, namespace string, label string, name string, excludedNamespaces []string) ([]models.Resource, error) {
	fmt.Println("Scanning for unused resources...")
	var resources []models.Resource

	namespacesToScan, err := scanner.GetNamespacesToScan(client.Clientset(), namespace, excludedNamespaces)
	if err != nil {
		return nil, fmt.Errorf("error determining namespaces to scan: %w", err)
	}

	if len(namespacesToScan) == 0 {
		return resources, nil
	}

	listOptions := metav1.ListOptions{}
	if label != "" {
		listOptions.LabelSelector = label
	}
	if name != "" {
		listOptions.FieldSelector = "metadata.name=" + name
	}

	for _, ns := range namespacesToScan {
		scanNamespaceResources(client, ns, listOptions, &resources)
	}

	return resources, nil
}

func scanNamespaceResources(client kubernetes.Client, ns string, listOptions metav1.ListOptions, resources *[]models.Resource) {
	unusedConfigMaps, err := scanner.GetUnusedConfigMaps(client.Clientset(), ns, listOptions)
	if err != nil {
		fmt.Printf("Error getting unused configmaps: %v\n", err)
		return
	}
	*resources = append(*resources, unusedConfigMaps...)

	unusedSecrets, err := scanner.GetUnusedSecrets(client.Clientset(), ns, listOptions)
	if err != nil {
		fmt.Printf("Error getting unused secrets: %v\n", err)
		return
	}
	*resources = append(*resources, unusedSecrets...)

	orphanServices, err := scanner.GetOrphanServices(client.Clientset(), ns, listOptions)
	if err != nil {
		fmt.Printf("Error getting orphan services: %v\n", err)
		return
	}
	*resources = append(*resources, orphanServices...)

	oldReplicaSets, err := scanner.GetOldReplicaSets(client.Clientset(), ns, listOptions)
	if err != nil {
		fmt.Printf("Error getting old replicasets: %v\n", err)
		return
	}
	*resources = append(*resources, oldReplicaSets...)

	completedJobs, err := scanner.GetCompletedJobs(client.Clientset(), 24*time.Hour, ns, listOptions)
	if err != nil {
		fmt.Printf("Error getting completed jobs: %v\n", err)
		return
	}
	*resources = append(*resources, completedJobs...)

	failedPods, err := scanner.GetFailedPods(client.Clientset(), ns, listOptions)
	if err != nil {
		fmt.Printf("Error getting failed pods: %v\n", err)
		return
	}
	*resources = append(*resources, failedPods...)

	completedPods, err := scanner.GetCompletedPods(client.Clientset(), 24*time.Hour, ns, listOptions)
	if err != nil {
		fmt.Printf("Error getting completed pods: %v\n", err)
		return
	}
	*resources = append(*resources, completedPods...)

	unusedPVCs, err := scanner.GetUnusedPersistentVolumeClaims(client.Clientset(), ns, listOptions)
	if err != nil {
		fmt.Printf("Error getting unused pvcs: %v\n", err)
		return
	}
	*resources = append(*resources, unusedPVCs...)
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().BoolVar(&dryRun, "dry-run", false, "run in dry-run mode")
	scanCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "specify the namespace to scan")
	scanCmd.Flags().String("exclude-namespaces", "", "comma-separated list of namespaces to exclude from scanning")
	scanCmd.Flags().String("label", "", "filter resources by label (e.g., 'app=nginx')")
	scanCmd.Flags().String("name", "", "filter resources by name")
}
