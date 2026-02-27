package logger

import (
	"fmt"
	"os"
	"time"
)

const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	dim     = "\033[2m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	white   = "\033[97m"
	gray    = "\033[90m"
)

func init() {
	enableVT()
}

func Banner(version, addr, configPath string, proxyCount int) {
	fmt.Println()
	fmt.Printf("  %s%s┌──────────────────────────────────────────┐%s\n", bold, cyan, reset)
	fmt.Printf("  %s%s│%s  %sproxy-cli%s  %s%s%-32s%s%s│%s\n",
		bold, cyan, reset,
		bold+white, reset,
		dim+gray, "v"+version, "", reset,
		bold+cyan, reset,
	)
	fmt.Printf("  %s%s└──────────────────────────────────────────┘%s\n", bold, cyan, reset)
	fmt.Println()
	fmt.Printf("  %s◆%s  addr     %s%s%s\n", cyan, reset, green+bold, addr, reset)
	fmt.Printf("  %s◆%s  config   %s%s%s\n", cyan, reset, white, configPath, reset)
	fmt.Printf("  %s◆%s  proxies  %s%d enabled%s\n", cyan, reset, yellow+bold, proxyCount, reset)
	fmt.Println()
	fmt.Printf("  %s%s──────────────────────────────────────────%s\n\n", dim, gray, reset)
}

func Proxy(host, scheme, name, proxyAddr string) {
	fmt.Printf("  %s%s%s  %s⇢%s  %-38s %s  %s%s%s  %s%s%s\n",
		dim, timestamp(), reset,
		cyan+bold, reset,
		host,
		schemeTag(scheme),
		white, name, reset,
		dim+gray, proxyAddr, reset,
	)
}

func Direct(host string) {
	fmt.Printf("  %s%s%s  %s→%s  %-38s %s[direct]%s\n",
		dim, timestamp(), reset,
		gray, reset,
		host,
		dim+gray, reset,
	)
}

func Error(host string, err error) {
	fmt.Printf("  %s%s%s  %s✗%s  %-38s %s%s%s\n",
		dim, timestamp(), reset,
		red+bold, reset,
		host,
		red, err.Error(), reset,
	)
}

func Fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "\n  %s✗  %s%s\n\n", red+bold, fmt.Sprintf(format, args...), reset)
	os.Exit(1)
}

func timestamp() string {
	return time.Now().Format("15:04:05")
}

func schemeTag(scheme string) string {
	switch scheme {
	case "SOCKS5":
		return fmt.Sprintf("%s[SOCKS5]%s", magenta+bold, reset)
	case "HTTP":
		return fmt.Sprintf("%s[HTTP]  %s", blue+bold, reset)
	case "HTTPS":
		return fmt.Sprintf("%s[HTTPS] %s", cyan+bold, reset)
	default:
		return fmt.Sprintf("%s[%s]%s", gray, scheme, reset)
	}
}
