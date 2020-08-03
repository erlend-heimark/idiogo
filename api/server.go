package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/erlend-heimark/idiogo/storage/mssql"
	"github.com/go-chi/chi"

	"github.com/erlend-heimark/idiogo/extdadjokes"
)

//TODO Handlers sin hovedoppgave er å håndtere requests og response. Burde ha hjelpemetoder

func New(path, port string, db mssql.Client, fetcher extdadjokes.Fetcher) *http.Server {
	router := chi.NewRouter()

	router.Route(path, func(r chi.Router) {
		r.Route("/dadjokes", func(r chi.Router) {
			r.Get("/{jokeId}", getDadJoke(db, fetcher))
			r.Get("/", getRandomDadJoke(fetcher))
			r.Post("/", createDadJoke(db))
		})
	})

	return &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}
}

func getDadJoke(db mssql.Client, fetcher extdadjokes.Fetcher) http.HandlerFunc {
	return createHandler(func(r *http.Request) apiResponse {
		jokeID := chi.URLParam(r, "jokeId")
		savedJoke, exists, err := db.GetDadJoke(r.Context(), jokeID)
		if err != nil {
			return apiResponse{nil, 500, err}
		}
		if exists {
			return apiResponse{savedJoke, 200, nil}
		}
		extJoke, exists, err := fetcher.Get(jokeID)
		if err != nil {
			return apiResponse{nil, 500, err}
		}
		if !exists {
			return apiResponse{nil, 404, nil}
		}
		return apiResponse{mapExternalDadJokeToInternal(extJoke), 200, nil}
	})
}

func getRandomDadJoke(fetcher extdadjokes.Fetcher) http.HandlerFunc {
	return createHandler(func(r *http.Request) apiResponse {
		d, err := fetcher.GetRandom()
		if err != nil {
			return apiResponse{nil, 500, err}
		}
		return apiResponse{d, 200, nil}
	})
}

func createDadJoke(db mssql.Client) http.HandlerFunc {
	return createHandler(func(r *http.Request) apiResponse {
		var dadJoke mssql.DadJoke
		err := json.NewDecoder(r.Body).Decode(&dadJoke)
		if err != nil {
			return apiResponse{nil, 500, err}
		}

		err = db.CreateDadJoke(r.Context(), dadJoke)
		if err != nil {
			return apiResponse{nil, 500, err}
		}
		return apiResponse{nil, 201, nil}
	})
}

type apiResponse struct {
	data       interface{}
	statusCode int
	err        error
}

func createHandler(handleRequest func(r *http.Request) apiResponse) http.HandlerFunc {
	//Wrap handlers so that our handlers task is to return the response values. Enables us to have common functions for
	//response writing. Also avoids the empty return pattern.
	return func(w http.ResponseWriter, r *http.Request) {
		response := handleRequest(r)
		handleResponse(w, response)
	}
}

func handleResponse(w http.ResponseWriter, res apiResponse) {
	w.WriteHeader(res.statusCode)
	if res.err != nil {
		_, _ = w.Write([]byte(res.err.Error()))
	}
	if res.data == nil {
		return
	}
	resp, err := json.Marshal(res.data)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(resp)
}

func mapExternalDadJokeToInternal(external *extdadjokes.DadJoke) mssql.DadJoke {
	return mssql.DadJoke{
		ID:   external.Id,
		Joke: external.Joke,
	}
}
