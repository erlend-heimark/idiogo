package api

import (
	"fmt"
	"net/http"

	"github.com/erlend-heimark/idiogo/pkg/dadjokes"
)

func Init(path, port, dadJokeURL string) {
	http.HandleFunc(path + "/dadjokes", dadjokes.Handler(dadJokeURL))

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		panic(err)
	}
}
