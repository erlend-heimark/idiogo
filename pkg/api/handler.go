package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/erlend-heimark/idiogo/pkg/storage/mssql"
	"github.com/go-chi/chi"

	"github.com/erlend-heimark/idiogo/pkg/externaldadjokes"
)

//TODO Handlers sin hovedoppgave er å håndtere requests og response. Burde ha hjelpemetoder
//TODO Strukturer slik at vi unngår empty returns. Viktig å ha en struktur som ikke er så lett å bomme på

type Handler struct {
	path string
	port string
	db mssql.Client
	fetcher externaldadjokes.Fetcher
}

func New(path, port string, db mssql.Client, fetcher externaldadjokes.Fetcher) Handler {
	return Handler{
		path: path,
		port: port,
		db: db,
		fetcher: fetcher,
	}
}

func (h Handler) Init() {
	router := chi.NewRouter()

	router.Route(h.path, func(r chi.Router) {
		r.Route("/dadjokes", func(r chi.Router) {
			r.Get("/{jokeId}", getDadJoke(h.db, h.fetcher))
			r.Get("/", getRandomDadJoke(h.fetcher))
			r.Post("/", createDadJoke(h.db))
		})
	})

	s := &http.Server{
		Addr:    fmt.Sprintf(":%s", h.port),
		Handler: router,
	}

	err := s.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func getDadJoke(db mssql.Client, fetcher externaldadjokes.Fetcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jokeID := chi.URLParam(r, "jokeId")
		savedJoke, exists, err := db.GetDadJoke(r.Context(), jokeID)
		if err != nil {
			writeErrorResponse(w, err)
			return
		}
		if exists {
			writeSuccessResponse(w, savedJoke)
			return
		}
		extJoke, exists, err := fetcher.Get(jokeID)
		if err != nil {
			writeErrorResponse(w, err)
		}
		if !exists {
			w.WriteHeader(404)
			return
		}

		writeSuccessResponse(w, mapExternalDadjokeToInternal(extJoke))
	}
}

func getRandomDadJoke(fetcher externaldadjokes.Fetcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := fetcher.GetRandom()
		if err != nil {
			writeErrorResponse(w, err)
			return
		}

		writeSuccessResponse(w, mapExternalDadjokeToInternal(d))
	}
}

func createDadJoke(db mssql.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dadJoke mssql.DadJoke
		err := json.NewDecoder(r.Body).Decode(&dadJoke)
		if err != nil {
			writeErrorResponse(w, err)
			//TODO Find a way to avoid these empty returns
			return
		}

		err = db.CreateDadJoke(r.Context(), dadJoke)
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

func writeSuccessResponse(w http.ResponseWriter, response interface{}) {
	resp, err := json.Marshal(response)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, err = w.Write(resp)
	if err != nil {
		panic("oops")
	}
}

func mapExternalDadjokeToInternal(external *externaldadjokes.DadJoke) mssql.DadJoke {
	return mssql.DadJoke{
		ID:   external.Id,
		Joke: external.Joke,
	}
}
