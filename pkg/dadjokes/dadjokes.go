package dadjokes

import (
	"encoding/json"
	"net/http"
)

type dadJoke struct {
	Id     string `json:"id"`
	Joke   string `json:"joke"`
	Status int    `json:"status"`
}

func Handler(dadJokeURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := http.NewRequest(http.MethodGet, dadJokeURL, nil)
		if err != nil {
			writeErrorResponse(w, err)
			return
		}
		req.Header.Set("Accept", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			writeErrorResponse(w, err)
			return
		}
		var dadJoke dadJoke
		err = json.NewDecoder(res.Body).Decode(&dadJoke)
		if err != nil {
			writeErrorResponse(w, err)
			return
		}

		d, err := json.Marshal(dadJoke)
		if err != nil {
			writeErrorResponse(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, err = w.Write(d)
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