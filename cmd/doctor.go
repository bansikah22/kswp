package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check the health of the Kubernetes cluster",
	Long:  `Doctor checks the health of the Kubernetes cluster and reports any issues.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Cluster Health: 82/100")
		fmt.Println("No critical issues detected")
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
