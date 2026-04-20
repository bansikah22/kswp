package report

import (
	"fmt"

	"github.com/bansikah22/kswp/pkg/models"
)

func PrintReport(resources []models.Resource) {
	if len(resources) == 0 {
		fmt.Println("No unused resources found!")
		return
	}

	// Use the new table-based output
	PrintResourcesTable(resources, &TableConfig{
		ShowReason: true,
	})

	score := CalculateHygieneScore(resources)
	fmt.Printf("\nCluster Hygiene Score: %d/100\n", score)
}

func CalculateHygieneScore(resources []models.Resource) int {
	score := 100
	for _, res := range resources {
		switch res.Kind {
		case "ConfigMap":
			score -= 1
		case "Secret":
			score -= 1
		case "Service":
			score -= 2
		case "ReplicaSet":
			score -= 2
		case "Job":
			score -= 1
		}
	}
	if score < 0 {
		return 0
	}
	return score
}
