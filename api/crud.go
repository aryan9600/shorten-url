package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/boltdb/bolt"
)

type PathToURL struct {
	Path string `json:"path"`
	URL  string `json:"url"`
}

type Path struct {
	Path string `json:"path"`
}

func PostPath(db *bolt.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {

		var newPath PathToURL

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err.Error())
		}

		json.Unmarshal(body, &newPath)
		path := newPath.Path
		URL := newPath.URL

		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("pathToUrls"))
			err := b.Put([]byte(path), []byte(URL))
			if err != nil {
				panic(err.Error())
			}
			return nil
		})

		resp, err := json.Marshal(newPath)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(resp)
	}
	return fn
}

func GetURL(db *bolt.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {

		var path Path

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err.Error())
		}

		json.Unmarshal(body, &path)

		tx, err := db.Begin(true)
		if err != nil {
			panic(err.Error())
		}

		defer tx.Rollback()

		b := tx.Bucket([]byte("pathToUrls"))
		URL := b.Get([]byte(path.Path))

		resp, err := json.Marshal(string(URL))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}

	return fn
}

func DeletePath(db *bolt.DB) http.HandlerFunc{
	fn := func(w http.ResponseWriter, r *http.Request) {

		var path Path

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err.Error())
		}

		json.Unmarshal(body, &path)

		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("pathToUrls"))
			err := b.Delete([]byte(path.Path))
			return err
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"Success": true}`))
	}
	return fn
}
