package main

import (
	"github.com/rueian/pgbroker-static/pkg/config"
	"github.com/rueian/pgbroker-static/pkg/proxy"
	"github.com/rueian/pgbroker-static/pkg/resolver"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	log.SetFlags(0)
}

func main() {
	path := env("CONFIG_PATH", "/config.yml")
	settings, err := config.NewSettings(path)
	if err != nil {
		log.Fatalf("Fail to load config from %s\n", path)
	}
	go settings.Watch()

	addr := env("PROXY_ADDR", ":5432")
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Fail to start proxy at %s\n", addr)
	}

	logging := env("ENABLE_LOGGING", "true")
	if logging != "true" {
		log.Printf("Query Logging is disabled.\n")
	}
	broker := proxy.NewPGBroker(&resolver.Static{Settings: settings}, logging == "true")

	go broker.Serve(ln)
	log.Println("proxy started at " + addr)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	log.Println("proxy is shutting down")
	broker.Shutdown()
}

func env(name, fallback string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return fallback
}
