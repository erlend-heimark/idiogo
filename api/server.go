package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/erlend-heimark/idiogo/storage/mssql"
	"github.com/go-chi/chi"

	"github.com/erlend-heimark/idiogo/extdadjokes"
)

type response struct {
	Data       interface{}
	StatusCode int
	Err        error
}

func New(path, port string, db mssql.Client, c extdadjokes.Client) *http.Server {
	r := chi.NewRouter()

	r.Route(path, func(r chi.Router) {
		r.Route("/dadjokes", func(r chi.Router) {
			handle(r, http.MethodGet, "/{jokeId}", getDadJoke(db, c))
			handle(r, http.MethodGet, "/", getRandomDadJoke(c))
			handle(r, http.MethodPost, "/", createDadJoke(db))
		})
	})

	return &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}
}

func handle(r chi.Router, method, path string, handleRequest func(r *http.Request) response) {
	//Wrap handlers so that our handlers task is to return the response values. Enables us to have common functions for
	//response writing. Also avoids the empty return pattern.
	r.Method(method, path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := handleRequest(r)
		handleResponse(w, response)
	}))
}

func handleResponse(w http.ResponseWriter, res response) {
	w.WriteHeader(res.StatusCode)
	if res.Err != nil {
		_, _ = w.Write([]byte(res.Err.Error()))
	}
	if res.Data == nil {
		return
	}
	resp, err := json.Marshal(res.Data)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(resp)
}