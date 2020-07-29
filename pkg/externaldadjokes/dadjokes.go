package externaldadjokes

import (
	"encoding/json"
	"net/http"
)

type DadJoke struct {
	Id     string `json:"id"`
	Joke   string `json:"joke"`
	Status int    `json:"status"`
}

func Get(dadJokeURL string) (*DadJoke, error) {
	req, err := http.NewRequest(http.MethodGet, dadJokeURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	var dadJoke DadJoke
	err = json.NewDecoder(res.Body).Decode(&dadJoke)
	if err != nil {
		return nil, err
	}
	//TODO Proper pointer use?
	return &dadJoke, nil
}
