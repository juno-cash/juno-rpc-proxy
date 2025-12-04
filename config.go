package main

import (
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

type ProxyAuth struct {
	Enabled  bool   `toml:"enabled"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}

type Upstream struct {
	URL      string `toml:"url"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Timeout  string `toml:"timeout"`
}

type Config struct {
	Listen         string    `toml:"listen"`
	ProxyAuth      ProxyAuth `toml:"proxy_auth"`
	Upstream       Upstream  `toml:"upstream"`
	AllowedMethods []string  `toml:"allowed_methods"`
}

func (c *Config) GetUpstreamTimeout() time.Duration {
	if c.Upstream.Timeout == "" {
		return 30 * time.Second
	}
	d, err := time.ParseDuration(c.Upstream.Timeout)
	if err != nil {
		return 30 * time.Second
	}
	return d
}

func (c *Config) IsMethodAllowed(method string) bool {
	for _, m := range c.AllowedMethods {
		if m == method {
			return true
		}
	}
	return false
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) validate() error {
	if c.Listen == "" {
		return fmt.Errorf("listen address is required")
	}
	if c.Upstream.URL == "" {
		return fmt.Errorf("upstream URL is required")
	}
	if len(c.AllowedMethods) == 0 {
		return fmt.Errorf("at least one allowed method is required")
	}
	if c.ProxyAuth.Enabled {
		if c.ProxyAuth.Username == "" || c.ProxyAuth.Password == "" {
			return fmt.Errorf("proxy_auth username and password are required when enabled")
		}
	}
	return nil
}
