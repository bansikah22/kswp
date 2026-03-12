package report

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/bansikah22/kswp/pkg/models"
)

func PrintReport(resources []models.Resource) {
	if len(resources) == 0 {
		fmt.Println("No unused resources found!")
		return
	}

	fmt.Println("Cluster Scan Report")
	fmt.Println("--------------------")

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, '\t', 0)
	fmt.Fprintln(w, "KIND\tNAMESPACE\tNAME\tREASON\tAGE")
	for _, res := range resources {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", res.Kind, res.Namespace, res.Name, res.Reason, res.Age.Round(time.Second).String())
	}
	w.Flush()

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
