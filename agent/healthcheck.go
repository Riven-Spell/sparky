package agent

import (
	"context"
	"time"

	"github.com/Riven-Spell/sparky/common/models"
)

func (d *Daemon) runHealthChecks(ctx context.Context) {
	ticker := time.NewTicker(d.config.HealthCheckFrequency)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			d.validateModels()

			// todo if the model list changed since the last tick, ping the server.
			//d.reportHealth()
		}
	}
}

func (d *Daemon) validateModels() []models.ClusterModel {
	return make([]models.ClusterModel, 0)
}
