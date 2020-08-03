package mssql

import (
	"context"
	"testing"

	"github.com/erlend-heimark/idiogo/pkg/envvars"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestClient_CreateGetDadJoke(t *testing.T) {
	e := envvars.New()
	db, err := NewClient(e.DBConnectionString)
	require.NoError(t, err)

	//Create joke
	id := uuid.New().String()[:10]
	joke := DadJoke{
		ID:   id,
		Joke: "Funny joke",
	}
	err = db.CreateDadJoke(context.Background(), joke)
	require.NoError(t, err)

	//Get joke
	fetchedJoke, exists, err := db.GetDadJoke(context.Background(), id)
	require.NoError(t, err)
	require.True(t, exists)
	require.Equal(t, joke.Joke, fetchedJoke.Joke)
}
