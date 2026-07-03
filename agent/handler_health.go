package agent

import (
	"encoding/json"
	"net/http"

	"github.com/Riven-Spell/sparky/common/models"
)

func handleGetHealth(d *Daemon) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(models.AgentHealth{
			Alive:        true,
			MemoryUsedGb: 0,
			Models:       []models.ClusterModel{},
		})
	}
}
