package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Auth struct {
	Credentials Credentials `json:"credentials"`
	Token       string      `json:"token"`
}

type Rule struct {
	Name  string   `json:"name"`
	Hosts []string `json:"hosts"`
}

type ProxyConfig struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Scheme  string `json:"scheme"`
	Host    string `json:"host"`
	Port    int    `json:"port"`
	Auth    Auth   `json:"auth"`
	Rules   []Rule `json:"rules"`
}

func (p *ProxyConfig) Addr() string {
	return fmt.Sprintf("%s:%d", p.Host, p.Port)
}

func (p *ProxyConfig) HasCredentials() bool {
	return p.Auth.Credentials.Username != ""
}

func Load(path string) ([]ProxyConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var all []ProxyConfig
	if err := json.Unmarshal(data, &all); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	var enabled []ProxyConfig
	for _, c := range all {
		if c.Enabled {
			enabled = append(enabled, c)
		}
	}

	return enabled, nil
}
