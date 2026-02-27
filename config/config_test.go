package config

import (
	"encoding/json"
	"testing"
)

func TestLoadFromJSON(t *testing.T) {
	raw := `[
		{
			"name": "Test Proxy",
			"enabled": true,
			"scheme": "HTTP",
			"host": "proxy.example.com",
			"port": 8080,
			"auth": {
				"credentials": { "username": "user", "password": "pass" },
				"token": ""
			},
			"rules": [
				{ "name": "YouTube", "hosts": ["*.youtube.com", "*.googlevideo.com"] }
			]
		},
		{
			"name": "Disabled Proxy",
			"enabled": false,
			"scheme": "HTTP",
			"host": "dead.proxy.com",
			"port": 1234,
			"auth": { "credentials": { "username": "", "password": "" }, "token": "" },
			"rules": []
		}
	]`

	var all []ProxyConfig
	if err := json.Unmarshal([]byte(raw), &all); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(all) != 2 {
		t.Fatalf("expected 2 configs, got %d", len(all))
	}

	p := all[0]

	if p.Name != "Test Proxy" {
		t.Errorf("Name = %q, want %q", p.Name, "Test Proxy")
	}
	if p.Port != 8080 {
		t.Errorf("Port = %d, want 8080", p.Port)
	}
	if p.Addr() != "proxy.example.com:8080" {
		t.Errorf("Addr() = %q, want %q", p.Addr(), "proxy.example.com:8080")
	}
	if !p.HasCredentials() {
		t.Error("HasCredentials() = false, want true")
	}
	if len(p.Rules) != 1 {
		t.Errorf("len(Rules) = %d, want 1", len(p.Rules))
	}
	if len(p.Rules[0].Hosts) != 2 {
		t.Errorf("len(Rules[0].Hosts) = %d, want 2", len(p.Rules[0].Hosts))
	}

	// Disabled proxy
	if all[1].Enabled {
		t.Error("second config should be disabled")
	}
}

func TestHasCredentials(t *testing.T) {
	empty := ProxyConfig{}
	if empty.HasCredentials() {
		t.Error("empty proxy should have no credentials")
	}

	withCreds := ProxyConfig{
		Auth: Auth{Credentials: Credentials{Username: "user", Password: "pass"}},
	}
	if !withCreds.HasCredentials() {
		t.Error("proxy with credentials should return true")
	}
}
