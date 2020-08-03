package main

import (
	"github.com/erlend-heimark/idiogo/pkg/api"
	"github.com/erlend-heimark/idiogo/pkg/envvars"
	"github.com/erlend-heimark/idiogo/pkg/externaldadjokes"
	"github.com/erlend-heimark/idiogo/pkg/storage/mssql"
)

func main() {
	envVars := envvars.New()

	client, err := mssql.NewClient(envVars.DBConnectionString)
	if err != nil {
		panic(err)
	}
	dadJokeFetcher := externaldadjokes.NewFetcher(envVars.DadJokeURL)

	s := api.New(envVars.Path, envVars.Port, client, dadJokeFetcher)
	err = s.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
