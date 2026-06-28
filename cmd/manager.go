package cmd

import "github.com/spf13/cobra"

var managerCmd = &cobra.Command{
	Use:   "manager",
	Short: "Run the Sparky manager (web UI + router)",
	RunE: func(cmd *cobra.Command, args []string) error {
		panic("not implemented")
	},
}

func init() {
	rootCmd.AddCommand(managerCmd)
}
