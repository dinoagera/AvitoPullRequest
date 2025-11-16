package main

import (
	"log"

	"github.com/dinoagera/AvitoPullRequest/config"
	"github.com/dinoagera/AvitoPullRequest/internal/app"
	"github.com/dinoagera/AvitoPullRequest/pkg/logger"
)

func main() {
	l := logger.InitLogger()
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal("failed to init config", "err", err)
	}
	app.Run(cfg, l)
}
