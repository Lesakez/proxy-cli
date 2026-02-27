package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Lesakez/proxy-cli/config"
	"github.com/Lesakez/proxy-cli/logger"
	"github.com/Lesakez/proxy-cli/proxy"
)

const version = "1.0.0"

func main() {
	port := flag.Int("port", 5555, "Local port to listen on")
	configPath := flag.String("config", defaultConfigPath(), "Path to JSON config file")
	showVersion := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("proxy-cli v%s\n", version)
		os.Exit(0)
	}

	proxies, err := config.Load(*configPath)
	if err != nil {
		logger.Fatal("failed to load config: %v", err)
	}
	if len(proxies) == 0 {
		logger.Fatal("no enabled proxies found in: %s", *configPath)
	}

	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	logger.Banner(version, addr, *configPath, len(proxies))

	server := proxy.NewServer(addr, proxies)
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal("server error: %v", err)
	}
}

func defaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "config.json"
	}
	return filepath.Join(home, ".proxy-cli", "config.json")
}
