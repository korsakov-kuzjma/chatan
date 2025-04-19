package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/korsakov-kuzjma/chatan/internal/bot"
	"github.com/korsakov-kuzjma/chatan/internal/config"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	adminBot, err := bot.NewBot(cfg)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	go adminBot.Start()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	adminBot.Stop()
	log.Println("Bot stopped gracefully")
}
