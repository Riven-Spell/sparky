package cmd

import "github.com/spf13/cobra"

var managerConfigValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate an existing manager configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		panic("not implemented")
	},
}

func init() {
	managerConfigCmd.AddCommand(managerConfigValidateCmd)
}
