package api

import (
	"encoding/json"
	"net/http"
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
