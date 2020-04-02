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

		w.Header().Set("Content-Type", "application/json")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.Unmarshal(body, &newPath)
		path := newPath.Path
		URL := newPath.URL

		er := AddPath(db, path, URL)

		if er == AlreadyPresentError(AlreadyPresentErrorString){
			w.WriteHeader(409)
			w.Write([]byte(`{"Error": "Path already exists"}`))
			return
		} else if er != nil{
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			resp, _ := json.Marshal(newPath)

			w.WriteHeader(http.StatusCreated)
			w.Write(resp)
		}
	}
	return fn
}

func GetURL(db *bolt.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {

		query := r.URL.Query()
		token := query.Get("path")

		URL, err := FetchURL(db, token)

		if err == NotFoundError(NotFoundErrorString){
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"Error": "Path not found."}`))
			return
		} else if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp, err := json.Marshal(string(URL))
		if err!=nil{
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}

	return fn
}

func DeletePath(db *bolt.DB) http.HandlerFunc{
	fn := func(w http.ResponseWriter, r *http.Request) {

		query := r.URL.Query()
		token := query.Get("path")

		er := RemovePath(db, token)

		if er!=nil{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"Success": true}`))
	}
	return fn
}
