package filter

import (
	"strings"

	"github.com/Lesakez/proxy-cli/config"
)

func MatchHost(host, pattern string) bool {
	if idx := strings.LastIndex(host, ":"); idx != -1 {
		host = host[:idx]
	}

	host = strings.ToLower(host)
	pattern = strings.ToLower(pattern)

	if host == pattern {
		return true
	}

	if !strings.Contains(pattern, "*") {
		return false
	}

	return wildcardMatch(host, pattern)
}

func wildcardMatch(s, pattern string) bool {
	parts := strings.Split(pattern, "*")

	pos := 0
	for i, part := range parts {
		if part == "" {
			continue
		}
		idx := strings.Index(s[pos:], part)
		if idx == -1 {
			return false
		}
		if i == 0 && idx != 0 {
			return false
		}
		pos += idx + len(part)
	}

	if !strings.HasSuffix(pattern, "*") {
		lastPart := parts[len(parts)-1]
		if lastPart != "" && !strings.HasSuffix(s, lastPart) {
			return false
		}
	}

	return true
}

func FindProxy(host string, proxies []config.ProxyConfig) *config.ProxyConfig {
	for i := range proxies {
		p := &proxies[i]
		for _, rule := range p.Rules {
			for _, pattern := range rule.Hosts {
				if MatchHost(host, pattern) {
					return p
				}
			}
		}
	}
	return nil
}
