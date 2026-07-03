package cmd

import (
	"fmt"

	"github.com/Riven-Spell/sparky/common/enum"
	"github.com/Riven-Spell/sparky/manager"
	"github.com/spf13/cobra"
)

var managerConfigValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate an existing manager configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		config := manager.Config{}
		configPath := enum.EEnvironmentVariable.SparkyConfig().Get()
		if err := config.Load(configPath); err != nil {
			return fmt.Errorf("load configuration: %w", err)
		}
		if err := config.Validate(); err != nil {
			return fmt.Errorf("validate configuration: %w", err)
		}
		fmt.Println("configuration is valid")
		return nil
	},
}

func init() {
	managerConfigCmd.AddCommand(managerConfigValidateCmd)
}
