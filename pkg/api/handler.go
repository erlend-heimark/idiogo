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

type Handler struct {
	path    string
	port    string
	db      mssql.Client
	fetcher externaldadjokes.Fetcher
}

func New(path, port string, db mssql.Client, fetcher externaldadjokes.Fetcher) Handler {
	return Handler{
		path:    path,
		port:    port,
		db:      db,
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
	return createHandler(func(r *http.Request) apiResponse {
		jokeID := chi.URLParam(r, "jokeId")
		savedJoke, exists, err := db.GetDadJoke(r.Context(), jokeID)
		if err != nil {
			return apiResponse{
				nil, 500, err,
			}
		}
		if exists {
			return apiResponse{savedJoke, 200, nil}
		}
		extJoke, exists, err := fetcher.Get(jokeID)
		if err != nil {
			return apiResponse{
				nil, 0, err,
			}
		}
		if !exists {
			return apiResponse{nil, 404, nil}
		}
		return apiResponse{mapExternalDadjokeToInternal(extJoke), 200, nil}
	})
}

func getRandomDadJoke(fetcher externaldadjokes.Fetcher) http.HandlerFunc {
	return createHandler(func(r *http.Request) apiResponse {
		d, err := fetcher.GetRandom()
		if err != nil {
			return apiResponse{
				nil, 500, err,
			}
		}
		return apiResponse{d, 200, nil}
	})
}

func createDadJoke(db mssql.Client) http.HandlerFunc {
	return createHandler(func(r *http.Request) apiResponse {
		var dadJoke mssql.DadJoke
		err := json.NewDecoder(r.Body).Decode(&dadJoke)
		if err != nil {
			return apiResponse{
				nil, 500, err,
			}
		}

		err = db.CreateDadJoke(r.Context(), dadJoke)
		if err != nil {
			return apiResponse{
				nil, 500, err,
			}
		}
		return apiResponse{nil, 200, nil}
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
		writeErrorResponse(w, res.err)
	}
	if res.data != nil {
		resp, err := json.Marshal(res)
		if err != nil {
			writeErrorResponse(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(resp)
		if err != nil {
			panic("oops")
		}
	}
}

func writeErrorResponse(w http.ResponseWriter, e error) {
	w.WriteHeader(500)
	_, err := w.Write([]byte(e.Error()))
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
