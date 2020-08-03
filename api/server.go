package api

import (
	"fmt"
	"net/http"

	"github.com/erlend-heimark/idiogo/storage/mssql"
	"github.com/go-chi/chi"

	"github.com/erlend-heimark/idiogo/extdadjokes"
)

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

//handle wraps the handler logic so that our handlers task is to return the response values. Enables us to have common functions for
//response writing. Also avoids the empty return pattern.
func handle(r chi.Router, method, path string, handleReq func(r *http.Request) response) {
	r.Method(method, path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := handleReq(r)
		handleResponse(w, response)
	}))
}