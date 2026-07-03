package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Riven-Spell/sparky/manager"
	"github.com/spf13/cobra"
)

var managerCmd = &cobra.Command{
	Use:   "manager",
	Short: "Run the Sparky manager (web UI + router)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := manager.NewDefaultConfig()
		d := manager.NewDaemon(cfg)

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
	rootCmd.AddCommand(managerCmd)
}
