package router

import (
	"fmt"
	"net/http"

	"github.com/erlend-heimark/idiogo/pkg/something"
)

func Init(path, port string) {
	http.HandleFunc(path + "/something", something.Handler)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		panic(err)
	}
}
