package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/erlend-heimark/idiogo/extdadjokes"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/erlend-heimark/idiogo/envvars"
	"github.com/erlend-heimark/idiogo/storage/mssql"
)

func TestCreateAndGet(t *testing.T) {
	e := envvars.New()
	db, err := mssql.NewClient(e.DBConnectionString)
	require.NoError(t, err)
	fetcher := extdadjokes.NewClient(e.DadJokeURL)
	server := New(e.Path, e.Port, db, fetcher)

	url := fmt.Sprintf("%s/%s", e.Path, "dadjokes")

	//Create joke
	jokeID := uuid.New().String()[:10]
	createBody := mssql.DadJoke{
		ID:   jokeID,
		Joke: "Funnny joke",
	}
	b, err := json.Marshal(createBody)
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	require.NoError(t, err)
	rrPost := httptest.NewRecorder()
	server.Handler.ServeHTTP(rrPost, req)
	require.Equal(t, 201, rrPost.Code)

	//Fetch created joke
	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", url, jokeID), nil)
	require.NoError(t, err)
	rrGet := httptest.NewRecorder()
	server.Handler.ServeHTTP(rrGet, req)
	require.Equal(t, 200, rrGet.Code)
	var response mssql.DadJoke
	err = json.Unmarshal(rrGet.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Equal(t, createBody.ID, response.ID)
	require.Equal(t, createBody.Joke, response.Joke)
}
