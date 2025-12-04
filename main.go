package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var version = "1.0.0"

func main() {
	configPath := flag.String("config", "config.toml", "Path to configuration file")
	showVersion := flag.Bool("version", false, "Show version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("juno-proxy version %s\n", version)
		os.Exit(0)
	}

	config, err := LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	proxy := NewProxy(config)
	if err := proxy.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
