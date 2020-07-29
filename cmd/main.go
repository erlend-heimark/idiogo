package main

import (
	"github.com/erlend-heimark/idiogo/pkg/config"
	"github.com/erlend-heimark/idiogo/pkg/api"
)

func main() {
	cfg := config.New()

	api.Init(cfg.Path, cfg.Port, cfg.DadJokeURL)
}