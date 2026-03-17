package cmd

import (
	"fmt"
	"os"

	"github.com/bansikah22/kswp/internal/branding"
	"github.com/spf13/cobra"
)

var (
	showBanner = true // Flag to control banner display
)

var rootCmd = &cobra.Command{
	Use:     "kswp",
	Short:   "kswp is a Kubernetes cluster hygiene tool",
	Long:    `kswp is a Kubernetes cluster hygiene tool that detects and safely cleans unused resources.`,
	Version: branding.Version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Display banner unless help, version, or other non-interactive flags are used
		if showBanner && !isHelpOrVersionRequested(cmd, args) {
			branding.DisplayBanner()
			showBanner = false // Only show once per execution
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Run 'kswp --help' for usage information")
		}
	},
}

// isHelpOrVersionRequested checks if help or version flags are in the arguments
func isHelpOrVersionRequested(cmd *cobra.Command, args []string) bool {
	// Check for help flag
	helpFlag, _ := cmd.Flags().GetBool("help")
	if helpFlag {
		return true
	}

	// Check for version flag
	versionFlag, _ := cmd.Flags().GetBool("version")
	if versionFlag {
		return true
	}

	// Check in os.Args as well for cases where Cobra hasn't parsed yet
	for _, arg := range os.Args {
		if arg == "-h" || arg == "--help" || arg == "-v" || arg == "--version" {
			return true
		}
	}

	return false
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
