package main

import (
	"log"
	"net/http"
	"url-short/api"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
)

func main() {

	db, err := bolt.Open("paths.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("pathToUrls"))
		if err != nil {
			panic(err.Error())
		}
		return nil
	})

	defer db.Close()

	r := mux.NewRouter()

	c := http.NewServeMux()

	API := r.PathPrefix("/api/v1").Subrouter()
	Redirect := r.PathPrefix("/").Subrouter()

	API.HandleFunc("/shorten", api.GetURL(db)).Methods("GET")
	API.HandleFunc("/shorten", api.PostPath(db)).Methods("POST")
	API.HandleFunc("/shorten", api.DeletePath(db)).Methods("DELETE")
	Redirect.HandleFunc("", MapHandler(db, c)).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", r))
}

func MapHandler(db *bolt.DB, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[1:]

		URL, err := api.FetchURL(db, path)

		if err == api.NotFoundError(api.NotFoundErrorString){
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else if URL == nil{
			http.Redirect(w, r, string(URL), http.StatusFound)
			return
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}
