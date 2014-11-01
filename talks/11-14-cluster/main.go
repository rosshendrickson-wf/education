package main

import (
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/predict/", classifyUser)
	panic(http.ListenAndServe(":8080", nil))
}

func classifyUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
