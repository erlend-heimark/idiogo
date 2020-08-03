package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

type response struct {
	data       interface{}
	statusCode int
	err        error
}

func newResponse(d interface{}, s int , e error) response {
	return response{
		data:       d,
		statusCode: s,
		err:        e,
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

func handleResponse(w http.ResponseWriter, res response) {
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