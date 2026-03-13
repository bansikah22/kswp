package cmd

import (
	"fmt"

	"github.com/bansikah22/kswp/internal/config"
	"github.com/bansikah22/kswp/internal/notifications"
	"github.com/spf13/cobra"
)

var notifyCmd = &cobra.Command{
	Use:   "notify",
	Short: "Send notifications",
	Long:  `Send notifications to various channels.`,
}

var slackCmd = &cobra.Command{
	Use:   "slack",
	Short: "Send a Slack notification",
	Long:  `Send a Slack notification to a specified webhook URL.`,
	Run: func(cmd *cobra.Command, args []string) {
		message, _ := cmd.Flags().GetString("message")

		if message == "" {
			fmt.Println("Please provide a message.")
			return
		}

		cfg, err := config.Load()
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}

		for _, notifierConfig := range cfg.Notifications {
			if notifierConfig.Type == "slack" {
				notifier, err := notifications.NewNotifier(notifierConfig)
				if err != nil {
					fmt.Println("Error creating notifier:", err)
					continue
				}
				err = notifier.Send(message)
				if err != nil {
					fmt.Println("Error sending notification:", err)
				} else {
					fmt.Println("Slack notification sent successfully.")
				}
			}
		}
	},
}

var emailCmd = &cobra.Command{
	Use:   "email",
	Short: "Send an email notification",
	Long:  `Send an email notification to the configured recipients.`,
	Run: func(cmd *cobra.Command, args []string) {
		message, _ := cmd.Flags().GetString("message")

		if message == "" {
			fmt.Println("Please provide a message.")
			return
		}

		cfg, err := config.Load()
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}

		for _, notifierConfig := range cfg.Notifications {
			if notifierConfig.Type == "email" {
				notifier, err := notifications.NewNotifier(notifierConfig)
				if err != nil {
					fmt.Println("Error creating notifier:", err)
					continue
				}
				err = notifier.Send(message)
				if err != nil {
					fmt.Println("Error sending notification:", err)
				} else {
					fmt.Println("Email notification sent successfully.")
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(notifyCmd)
	notifyCmd.AddCommand(slackCmd)
	notifyCmd.AddCommand(emailCmd)
	slackCmd.Flags().String("message", "", "Message to send")
	emailCmd.Flags().String("message", "", "Message to send")
	notifyCmd.PersistentFlags().Bool("dry-run", false, "run in dry-run mode")
}
