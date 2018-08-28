package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"log"
	"core"
)

type Router interface {
	Search()
	Index()
	Delete()
}

func Search(w http.ResponseWriter, r *http.Request) {
	var sf = core.StoreFields{}
	sf.PerformSearch(w, r)
}

func Index(w http.ResponseWriter, r *http.Request) {
	var sf = core.StoreFields{}
	sf.BuildIndex(w, r)
}

func Delete(w http.ResponseWriter, r *http.Request) {

}

func IndexInfo(w http.ResponseWriter, r *http.Request) {

}

func Info(w http.ResponseWriter, r *http.Request) {

}

func main() {
	r := mux.NewRouter()
	r.Headers("Content-Type", "application/json")

	// get index info requests
	r.HandleFunc("/{index}", IndexInfo).Methods("GET")
	r.HandleFunc("/{index}", Info).Methods("GET")

	// search requests
	r.HandleFunc("/{index}", Search).Methods("POST")
	r.HandleFunc("/{index}/{indextype}", Search).Methods("POST")
	r.HandleFunc("/{index}/{indextype}/{id:[0-9]+}", Search).Methods("GET")

	// insert/update requests
	r.HandleFunc("/{index}", Index).Methods("PUT")
	r.HandleFunc("/{index}/{indextype}", Index).Methods("PUT")

	// delete doc request
	r.HandleFunc("/{index}", Delete).Methods("DELETE")

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}