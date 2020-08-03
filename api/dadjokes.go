package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/erlend-heimark/idiogo/extdadjokes"
	"github.com/erlend-heimark/idiogo/storage/mssql"
)

func getDadJoke(db mssql.Client, c extdadjokes.Client) func(r *http.Request) response {
	return func(r *http.Request) response {
		jokeID := chi.URLParam(r, "jokeId")
		savedJoke, exists, err := db.GetDadJoke(r.Context(), jokeID)
		if err != nil {
			return response{nil, 500, err}
		}
		if exists {
			return response{savedJoke, 200, nil}
		}
		extJoke, exists, err := c.Get(jokeID)
		if err != nil {
			return response{nil, 500, err}
		}
		if !exists {
			return response{nil, 404, nil}
		}
		return response{mapExternalDadJokeToInternal(extJoke), 200, nil}
	}
}

func getRandomDadJoke(c extdadjokes.Client) func(r *http.Request) response {
	return func(r *http.Request) response {
		d, err := c.GetRandom()
		if err != nil {
			return response{nil, 500, err}
		}
		return response{d, 200, nil}
	}
}

func createDadJoke(db mssql.Client) func(r *http.Request) response {
	return func(r *http.Request) response {
		var dadJoke mssql.DadJoke
		err := json.NewDecoder(r.Body).Decode(&dadJoke)
		if err != nil {
			return response{nil, 500, err}
		}

		err = db.CreateDadJoke(r.Context(), dadJoke)
		if err != nil {
			return response{nil, 500, err}
		}
		return response{nil, 201, nil}
	}
}

func mapExternalDadJokeToInternal(d *extdadjokes.DadJoke) mssql.DadJoke {
	return mssql.DadJoke{
		ID:   d.Id,
		Joke: d.Joke,
	}
}
