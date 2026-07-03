package cmd

import (
	"fmt"

	"github.com/Riven-Spell/sparky/common/enum"
	"github.com/Riven-Spell/sparky/manager"
	"github.com/spf13/cobra"
)

var managerConfigCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new manager configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		config := manager.NewDefaultConfig()
		configPath := enum.EEnvironmentVariable.SparkyConfig().Get()
		if err := config.Save(configPath); err != nil {
			return fmt.Errorf("save configuration: %w", err)
		}
		fmt.Printf("configuration created at %s\n", configPath)
		return nil
	},
}

func init() {
	managerConfigCmd.AddCommand(managerConfigCreateCmd)
}
