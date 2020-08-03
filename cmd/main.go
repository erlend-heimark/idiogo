package main

import (
	"github.com/erlend-heimark/idiogo/api"
	"github.com/erlend-heimark/idiogo/envvars"
	"github.com/erlend-heimark/idiogo/extdadjokes"
	"github.com/erlend-heimark/idiogo/storage/mssql"
)

func main() {
	envVars := envvars.New()

	dbClient, err := mssql.NewClient(envVars.DBConnectionString)
	if err != nil {
		panic(err)
	}
	extDadJokeClient := extdadjokes.NewClient(envVars.DadJokeURL)

	s := api.New(envVars.Path, envVars.Port, dbClient, extDadJokeClient)
	err = s.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
