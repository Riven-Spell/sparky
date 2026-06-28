package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "sparky",
	Short: "Manage LLMs across a DGX Spark cluster",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize()
}
