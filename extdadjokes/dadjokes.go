package extdadjokes

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type DadJoke struct {
	Id     string `json:"id"`
	Joke   string `json:"joke"`
	Status int    `json:"status"`
}

type Fetcher struct {
	dadJokeURL string
}

func NewFetcher(dadJokeURL string) Fetcher {
	return Fetcher{dadJokeURL: dadJokeURL}
}

func (f Fetcher) GetRandom() (*DadJoke, error) {
	r, err := http.NewRequest(http.MethodGet, f.dadJokeURL, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}
	var d DadJoke
	err = json.NewDecoder(res.Body).Decode(&d)
	if err != nil {
		return nil, err
	}
	//TODO Proper pointer use?
	return &d, nil
}

func (f Fetcher) Get(jokeID string) (*DadJoke, bool, error) {
	r, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/j/%s", f.dadJokeURL, jokeID), nil)
	if err != nil {
		return nil, false, err
	}
	r.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, false, err
	}
	var d DadJoke
	err = json.NewDecoder(res.Body).Decode(&d)
	if err != nil {
		return nil, false, err
	}
	if d.Joke == "" {
		//API returned 200 ok but with no joke found
		return nil, false, nil
	}
	return &d, true, nil
}