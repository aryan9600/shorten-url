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

	API := r.PathPrefix("/api/v1").Subrouter()

	API.HandleFunc("/shorten", api.GetURL(db)).Methods("GET")
	API.HandleFunc("/shorten", api.PostPath(db)).Methods("POST")
	API.HandleFunc("/shorten", api.DeletePath(db)).Methods("DELETE")

	server := &http.Server{
		Addr:    ":8000",
		Handler: r,
	}
	go server.ListenAndServe()

	c := http.NewServeMux()
	pathHandler := MapHandler(db, c)

	client := &http.Server{
		Addr:    ":8080",
		Handler: pathHandler,
	}

	log.Fatal(client.ListenAndServe())

}

func MapHandler(db *bolt.DB, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[1:]

		tx, err := db.Begin(true)
		if err != nil {
			log.Fatal(err)
		}

		defer tx.Rollback()

		b := tx.Bucket([]byte("pathToUrls"))
		v := b.Get([]byte(path))
		if v == nil {
			w.WriteHeader(http.StatusNotFound)
			fallback.ServeHTTP(w, r)
			return
		}
		http.Redirect(w, r, string(v), http.StatusFound)
	}
}
