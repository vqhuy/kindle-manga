package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

func mkdir(name string) (string, error) {
	dir, err := ioutil.TempDir("", "kindle-manga-"+name+"-")
	if err != nil {
		return "", err
	}
	fp := filepath.Join(dir, name)
	err = os.MkdirAll(fp, 0755)
	return fp, err
}

func rmdir() {
	files, err := filepath.Glob(filepath.Join(os.TempDir(), "kindle-manga-*"))
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if err := os.RemoveAll(f); err != nil {
			panic(err)
		}
	}
}

func saveJSONToFile(path string, v interface{}) error {
	w, err := os.Create(path)
	if err != nil {
		return err
	}

	defer w.Close()

	return json.NewEncoder(w).Encode(v)
}

func loadJSONFromFile(path string, v interface{}) error {
	r, err := os.Open(path)
	if err != nil {
		return err
	}

	defer r.Close()

	return json.NewDecoder(r).Decode(v)
}
