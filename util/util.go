package util

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Mkdir(name string) (string, error) {
	dir, err := ioutil.TempDir("", "kindle-manga-"+name+"-")
	if err != nil {
		return "", err
	}
	fp := filepath.Join(dir, name)
	err = os.MkdirAll(fp, 0755)
	return fp, err
}

func Rmdir() {
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

func Mv(oldpath []os.FileInfo, dir, subdir string) error {
	newpath := filepath.Join(dir, subdir)
	if err := os.MkdirAll(newpath, 0755); err != nil {
		return err
	}

	for _, o := range oldpath {
		if err := os.Rename(filepath.Join(dir, o.Name()), filepath.Join(newpath, o.Name())); err != nil {
			return err
		}
	}
	return nil
}

func SaveJSONToFile(path string, v interface{}) error {
	w, err := os.Create(path)
	if err != nil {
		return err
	}

	defer w.Close()

	return json.NewEncoder(w).Encode(v)
}

func LoadJSONFromFile(path string, v interface{}) error {
	r, err := os.Open(path)
	if err != nil {
		return err
	}

	defer r.Close()

	return json.NewDecoder(r).Decode(v)
}
