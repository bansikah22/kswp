package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/bansikah22/kswp/internal/cleaner"
	"github.com/bansikah22/kswp/internal/kubernetes"
	"github.com/bansikah22/kswp/internal/scanner"
	"github.com/bansikah22/kswp/pkg/models"
	"github.com/bansikah22/kswp/test/mocks"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

		var resourcesToClean []models.Resource

		// Helper function to get flags, reducing boilerplate.
		getFlag := func(name string) bool {
			val, _ := cmd.Flags().GetBool(name)
			return val
		}

		// Check for expired resources first if the --ttl flag is provided.
		if getFlag("ttl") {
			label, _ := cmd.Flags().GetString("label")
			name, _ := cmd.Flags().GetString("name")
			listOptions := metav1.ListOptions{}
			if label != "" {
				listOptions.LabelSelector = label
			}
			if name != "" {
				listOptions.FieldSelector = "metadata.name=" + name
			}

			expiredResources, err := scanner.GetExpiredResources(client.Clientset(), cleanNamespace, listOptions)
			if err != nil {
				fmt.Println("Error scanning for expired resources:", err)
				return
			}
			resourcesToClean = append(resourcesToClean, expiredResources...)
		}

		// Check for unused resources based on the provided flags.
		all := getFlag("all")
		resourceFlags := map[string]bool{
			"ConfigMap":             getFlag("configmaps"),
			"Secret":                getFlag("secrets"),
			"Service":               getFlag("services"),
			"ReplicaSet":            getFlag("replicasets"),
			"Job":                   getFlag("jobs"),
			"Pod":                   getFlag("pods"),
			"PersistentVolumeClaim": getFlag("pvcs"),
		}

		// Determine if we need to scan for unused resources.
		shouldScanUnused := all
		for _, v := range resourceFlags {
			if v {
				shouldScanUnused = true
				break
			}
		}

		// If no specific resource type is requested, and --ttl is not used, default to all.
		if !getFlag("ttl") && !shouldScanUnused {
			all = true
			shouldScanUnused = true
		}

		if shouldScanUnused {
			label, _ := cmd.Flags().GetString("label")
			name, _ := cmd.Flags().GetString("name")
			unusedResources, err := ScanResources(client, cleanNamespace, label, name)
			if err != nil {
				fmt.Println("Error scanning for unused resources:", err)
				return
			}

			// Filter the unused resources based on the provided flags.
			for _, r := range unusedResources {
				if all || resourceFlags[r.Kind] {
					resourcesToClean = append(resourcesToClean, r)
				}
			}
		}

		if len(resourcesToClean) == 0 {
			fmt.Println("No resources found to clean.")
			return
		}

		// De-duplicate resources in case a resource is both expired and unused.
		resourcesToClean = deduplicateResources(resourcesToClean)

		if cleanDryRun {
			fmt.Println("Resources that would be deleted:")
			for _, res := range resourcesToClean {
				fmt.Printf("- %s/%s (%s)\n", res.Namespace, res.Name, res.Kind)
			}
			return
		}

		fmt.Printf("Delete %d resources? (y/N) ", len(resourcesToClean))
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

		for _, res := range resourcesToClean {
			err := cleaner.DeleteResource(client.Clientset(), res)
			if err != nil {
				fmt.Printf("Error deleting %s %s/%s: %s\n", res.Kind, res.Namespace, res.Name, err)
			}
		}
	},
}

func deduplicateResources(resources []models.Resource) []models.Resource {
	seen := make(map[string]bool)
	result := []models.Resource{}
	for _, r := range resources {
		key := fmt.Sprintf("%s/%s/%s", r.Kind, r.Namespace, r.Name)
		if _, ok := seen[key]; !ok {
			seen[key] = true
			result = append(result, r)
		}
	}
	return result
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().BoolVar(&cleanDryRun, "dry-run", false, "run in dry-run mode")
	cleanCmd.Flags().StringVarP(&cleanNamespace, "namespace", "n", "", "specify the namespace to clean")
	cleanCmd.Flags().String("label", "", "filter resources by label (e.g., 'app=nginx')")
	cleanCmd.Flags().String("name", "", "filter resources by name")
	cleanCmd.Flags().Bool("all", false, "clean all unused resources")
	cleanCmd.Flags().Bool("configmaps", false, "clean unused configmaps")
	cleanCmd.Flags().Bool("secrets", false, "clean unused secrets")
	cleanCmd.Flags().Bool("services", false, "clean unused services")
	cleanCmd.Flags().Bool("replicasets", false, "clean unused replicasets")
	cleanCmd.Flags().Bool("jobs", false, "clean unused jobs")
	cleanCmd.Flags().Bool("pods", false, "clean unused pods")
	cleanCmd.Flags().Bool("pvcs", false, "clean unused persistentvolumeclaims")
	cleanCmd.Flags().Bool("ttl", false, "clean expired resources based on the cleaner/ttl annotation")
}
