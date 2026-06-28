package cmd

import "github.com/spf13/cobra"

var managerConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Sparky manager configuration",
}

func init() {
	managerCmd.AddCommand(managerConfigCmd)
}
