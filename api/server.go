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