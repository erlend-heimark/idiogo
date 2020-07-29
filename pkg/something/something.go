package something

import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	_, err := w.Write([]byte("Hello world!"))
	if err != nil {
		panic("oops")
	}
}
