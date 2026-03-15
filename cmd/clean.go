package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/bansikah22/kswp/internal/cleaner"
	"github.com/bansikah22/kswp/internal/kubernetes"
	"github.com/bansikah22/kswp/pkg/models"
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

		label, _ := cmd.Flags().GetString("label")
		name, _ := cmd.Flags().GetString("name")
		resources, err := ScanResources(client, cleanNamespace, label, name)
		if err != nil {
			fmt.Println("Error scanning resources:", err)
			return
		}

		all, _ := cmd.Flags().GetBool("all")
		configmaps, _ := cmd.Flags().GetBool("configmaps")
		secrets, _ := cmd.Flags().GetBool("secrets")
		services, _ := cmd.Flags().GetBool("services")
		replicasets, _ := cmd.Flags().GetBool("replicasets")
		jobs, _ := cmd.Flags().GetBool("jobs")
		pods, _ := cmd.Flags().GetBool("pods")
		pvcs, _ := cmd.Flags().GetBool("pvcs")

		if !all && !configmaps && !secrets && !services && !replicasets && !jobs && !pods && !pvcs {
			all = true
		}

		var filteredResources []models.Resource
		for _, r := range resources {
			switch r.Kind {
			case "ConfigMap":
				if all || configmaps {
					filteredResources = append(filteredResources, r)
				}
			case "Secret":
				if all || secrets {
					filteredResources = append(filteredResources, r)
				}
			case "Service":
				if all || services {
					filteredResources = append(filteredResources, r)
				}
			case "ReplicaSet":
				if all || replicasets {
					filteredResources = append(filteredResources, r)
				}
			case "Job":
				if all || jobs {
					filteredResources = append(filteredResources, r)
				}
			case "Pod":
				if all || pods {
					filteredResources = append(filteredResources, r)
				}
			case "PersistentVolumeClaim":
				if all || pvcs {
					filteredResources = append(filteredResources, r)
				}
			}
		}
		resources = filteredResources

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
}
