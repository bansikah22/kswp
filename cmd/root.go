package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kswp",
	Short: "kswp is a Kubernetes cluster hygiene tool",
	Long:  `kswp is a Kubernetes cluster hygiene tool that detects and safely cleans unused resources.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello from kswp!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
