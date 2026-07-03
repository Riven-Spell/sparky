package agent

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"
)

const DefaultAgentPort = "8081"

type Config struct {
	ManagerHost string
	AgentPort   string
}

func NewDefaultConfig() Config {
	return Config{
		ManagerHost: "http://192.168.100.1:8080",
		AgentPort:   DefaultAgentPort,
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

	port, err := strconv.Atoi(c.AgentPort)
	if err != nil {
		return errors.New("agent_port: must be a number")
	}
	if port < 1 || port > 65535 {
		return errors.New("agent_port: must be between 1 and 65535")
	}

	return nil
}

func (c Config) listenAddress() string {
	return fmt.Sprintf("0.0.0.0:%s", c.AgentPort)
}

func (c *Config) ApplyFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&c.ManagerHost, "manager-host", c.ManagerHost, "Manager host URL (e.g. https://10.0.0.1:8080)")
	cmd.Flags().StringVar(&c.AgentPort, "agent-port", c.AgentPort, "Port for the agent HTTP server")
}
