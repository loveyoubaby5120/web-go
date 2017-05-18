package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// r.HandleFunc("/", YourHandler)
	r.HandleFunc("/articles/{category}/{id:[0-9]+}", ArticleHandler)

	// http.Handle("/", r)

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", r))
}

// YourHandler send
func YourHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello word!\n"))
}

// ArticleHandler is a test router
func ArticleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println(vars)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
}
