package cmd

import (
	"fmt"
	"os"

	"github.com/bansikah22/kswp/internal/kubernetes"
	"github.com/bansikah22/kswp/internal/scripting"
	"github.com/bansikah22/kswp/test/mocks"
	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply -f <script.lua>",
	Short: "Apply a Lua script to filter and delete resources",
	Long:  `Apply a Lua script to filter and delete unused resources.`,
	Run: func(cmd *cobra.Command, args []string) {
		scriptFile, _ := cmd.Flags().GetString("file")
		if scriptFile == "" {
			fmt.Println("Please provide a script file with -f")
			return
		}

		script, err := os.ReadFile(scriptFile)
		if err != nil {
			fmt.Println("Error reading script file:", err)
			return
		}

		var client kubernetes.Client
		if dryRun, _ := cmd.Flags().GetBool("dry-run"); dryRun {
			client = mocks.NewMockClient()
		} else {
			var err error
			client, err = kubernetes.NewClient()
			if err != nil {
				fmt.Println("Error connecting to Kubernetes cluster:", err)
				return
			}
		}

		err = scripting.Execute(string(script), client)
		if err != nil {
			fmt.Println("Error executing script:", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.Flags().StringP("file", "f", "", "path to the Lua script file")
	applyCmd.Flags().BoolVar(&dryRun, "dry-run", false, "run in dry-run mode")
}
