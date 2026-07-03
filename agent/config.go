package agent

import (
	"errors"
	"net/url"

	"github.com/spf13/cobra"
)

const DefaultAgentPort = "8081"

type Config struct {
	ManagerHost string
}

func NewDefaultConfig() Config {
	return Config{
		ManagerHost: "",
	}
}

func (c Config) Validate() error {
	if c.ManagerHost == "" {
		return errors.New("manager_host is required")
	}

	u, err := url.Parse(c.ManagerHost)
	if err != nil {
		return errors.New("manager_host: invalid URL")
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("manager_host: must have http or https scheme")
	}

	return nil
}

func (c *Config) ApplyFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&c.ManagerHost, "manager-host", c.ManagerHost, "Manager host URL (e.g. https://10.0.0.1:8080)")
}
