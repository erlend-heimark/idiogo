package main

import (
	"github.com/erlend-heimark/idiogo/pkg/config"
	"github.com/erlend-heimark/idiogo/pkg/router"
)

func main() {
	cfg := config.New()

	router.Init(cfg.Path, cfg.Port)
}