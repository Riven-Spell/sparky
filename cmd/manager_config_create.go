package cmd

import "github.com/spf13/cobra"

var managerConfigCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new manager configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		panic("not implemented")
	},
}

func init() {
	managerConfigCmd.AddCommand(managerConfigCreateCmd)
}
