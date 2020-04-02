package api

import (
	"github.com/boltdb/bolt"
)

func AddPath(db *bolt.DB, path string, URL string) error{
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("pathToUrls"))
		Path := FetchPath(db, path)
		if Path == nil{
			err := b.Put([]byte(path), []byte(URL))
			return err
		}
		return AlreadyPresentError(AlreadyPresentErrorString)
	})
	return err
}

func FetchURL(db *bolt.DB, path string) ([]byte, error){
	tx, err := db.Begin(true)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	b := tx.Bucket([]byte("pathToUrls"))
	URL := b.Get([]byte(path))
	if URL == nil{
		return nil, NotFoundError(NotFoundErrorString)
	}
	return URL, nil
}

func RemovePath(db *bolt.DB, path string) error{
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("pathToUrls"))
		err := b.Delete([]byte(path))
		return err
	})
	return err
}

func FetchPath(db *bolt.DB, path string) error{
	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	b := tx.Bucket([]byte("pathToUrls"))

	c := b.Cursor()

	for k, _ := c.First(); k != nil; k, _ = c.Next() {
		if string(k)==path{
			return AlreadyPresentError(AlreadyPresentErrorString)
		}
	}
	return nil
}




