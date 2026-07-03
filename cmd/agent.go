package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Riven-Spell/sparky/agent"
	"github.com/spf13/cobra"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Run the Sparky agent on a DGX Spark node",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := agent.NewDefaultConfig()
		cfg.ApplyFlags(cmd)
		d := agent.NewDaemon(cfg)

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		done := make(chan struct{})

		if err := d.Start(func() { close(done) }); err != nil {
			return err
		}

		select {
		case <-done:
		case <-ctx.Done():
			d.Stop()
			<-done
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(agentCmd)
}
