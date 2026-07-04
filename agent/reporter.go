package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Riven-Spell/sparky/common/models"
)

func (d *Daemon) reportHealth(ctx context.Context, health models.AgentHealth) error {
	url := d.config.ManagerHost + "/v1/agents/health"

	body, err := json.Marshal(health)
	if err != nil {
		return fmt.Errorf("marshal health report: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create health request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("send health report: %w", err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("health report rejected: %s", resp.Status)
	}

	return nil
}
