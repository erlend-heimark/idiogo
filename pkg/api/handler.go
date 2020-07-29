package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/erlend-heimark/idiogo/pkg/storage/mssql"
	"github.com/go-chi/chi"

	"github.com/erlend-heimark/idiogo/pkg/externaldadjokes"
)

func Init(path, port, dadJokeURL string, client mssql.Client) {
	router := chi.NewRouter()

	router.Route(path, func(r chi.Router) {
		r.Route("/dadjokes", func(r chi.Router) {
			r.Get("/", getDadJoke(dadJokeURL))
			r.Post("/", createDadJoke(client))
		})
	})

	s := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	err := s.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func getDadJoke(dadJokeURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := externaldadjokes.Get(dadJokeURL)
		if err != nil {
			writeErrorResponse(w, err)
			return
		}
		jokeJson, err := json.Marshal(d)
		if err != nil {
			writeErrorResponse(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, err = w.Write(jokeJson)
		if err != nil {
			panic("oops")
		}
	}
}

func createDadJoke(client mssql.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dadJoke mssql.DadJoke
		err := json.NewDecoder(r.Body).Decode(&dadJoke)
		if err != nil {
			writeErrorResponse(w, err)
			//TODO Find a way to avoid these empty returns
			return
		}

		err = client.CreateDadJoke(r.Context(), dadJoke)
		if err != nil {
			writeErrorResponse(w, err)
			return
		}
		w.WriteHeader(201)
	}
}

func writeErrorResponse(w http.ResponseWriter, e error) {
	w.WriteHeader(500)
	_, err := w.Write([]byte(e.Error()))
	if err != nil {
		panic("oops")
	}
}
