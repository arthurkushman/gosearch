package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"log"
	"fmt"
)

type Router interface {
	Search()
	Index()
	Delete()
}

func Search(w http.ResponseWriter, r *http.Request) {
	//PerformSearch(w, r)
	vars := mux.Vars(r)
	index := vars["index"]
	indexType := vars["indextype"]
	indexId := vars["id"]
	fmt.Println(index, indexType, indexId)
}

func Index(w http.ResponseWriter, r *http.Request) {

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
	r.HandleFunc("/articles", IndexInfo).Methods("GET")
	r.HandleFunc("/articles", Info).Methods("GET")
	r.HandleFunc("/{index}", Search).Methods("POST")
	r.HandleFunc("/{index}/{indextype}", Search).Methods("POST")
	r.HandleFunc("/{index}/{indextype}/{id:[0-9]+}", Search).Methods("GET")
	r.HandleFunc("/products", Index).Methods("POST")
	r.HandleFunc("/articles", Delete).Methods("DELETE")
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}