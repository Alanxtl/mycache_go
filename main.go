package main

import (
	"fmt"
	"net/http"
)

type server int

func (h *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	_, err := w.Write([]byte("Hello World"))
	if err != nil {
		return
	}
}

func main() {
	var s server
	http.ListenAndServe("localhost:9999", &s)
}
