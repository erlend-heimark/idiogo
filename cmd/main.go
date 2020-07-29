package main

import (
	"github.com/erlend-heimark/idiogo/pkg/api"
	"github.com/erlend-heimark/idiogo/pkg/config"
	"github.com/erlend-heimark/idiogo/pkg/storage/mssql"
)

func main() {
	cfg := config.New()

	client, err := mssql.NewClient(cfg.ConnectionString)
	if err != nil {
		panic(err)
	}

	api.Init(cfg.Path, cfg.Port, cfg.DadJokeURL, client)
}
