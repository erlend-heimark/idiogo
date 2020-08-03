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

type Client struct {
	dadJokeURL string
}

func NewClient(dadJokeURL string) Client {
	return Client{dadJokeURL: dadJokeURL}
}

func (c Client) GetRandom() (*DadJoke, error) {
	r, err := http.NewRequest(http.MethodGet, c.dadJokeURL, nil)
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

func (c Client) Get(jokeID string) (*DadJoke, bool, error) {
	r, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/j/%s", c.dadJokeURL, jokeID), nil)
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