package manager

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type RegisteredModelConfig struct {
	Recipe   string `yaml:"recipe"`
	Nodes    int    `yaml:"nodes"`
	Nickname string `yaml:"nickname,omitempty"`
}

type Config struct {
	ClusterName          string                  `yaml:"cluster_name"`
	ManagerHost          string                  `yaml:"manager_host"`
	WebUIHost            string                  `yaml:"webui_host"`
	AgentPort            int                     `yaml:"agent_port"`
	Models               []RegisteredModelConfig `yaml:"models"`
	EvictionIdleDuration  time.Duration          `yaml:"eviction_idle_duration"`
	HealthCheckFrequency  time.Duration          `yaml:"health_check_frequency"`
}

func NewDefaultConfig() Config {
	return Config{
		ClusterName: "default",
		ManagerHost:          "0.0.0.0:8080",
		WebUIHost:            "0.0.0.0:8081",
		AgentPort:            8085,
		Models: []RegisteredModelConfig{
			{Recipe: "@official/llama-3-70b", Nodes: 1},
		},
		EvictionIdleDuration: 15 * time.Minute,
		HealthCheckFrequency: 1 * time.Minute,
	}
}

func (c Config) Validate() error {
	if c.ClusterName == "" {
		return errors.New("cluster_name is required")
	}

	if c.ManagerHost == "" {
		return errors.New("manager_host is required")
	}

	if c.WebUIHost == "" {
		return errors.New("webui_host is required")
	}

	if c.AgentPort <= 0 || c.AgentPort > 65535 {
		return errors.New("agent_port must be between 1 and 65535")
	}

	if len(c.Models) == 0 {
		return errors.New("at least one model is required")
	}
	for i, m := range c.Models {
		if m.Recipe == "" {
			return fmt.Errorf("models[%d]: recipe is required", i)
		}
		if m.Nodes < 1 {
			return fmt.Errorf("models[%d]: nodes must be >= 1", i)
		}
	}

	if c.EvictionIdleDuration <= 0 {
		return errors.New("eviction_idle_duration must be positive")
	}

	if c.HealthCheckFrequency <= 0 {
		return errors.New("health_check_frequency must be positive")
	}

	return nil
}

func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}

func (c *Config) Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}
	if err := yaml.Unmarshal(data, c); err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}
	return nil
}
