package main

import (
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", snippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	mux.Handle("/static", http.StripPrefix("/static", fileServer))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Println("Running web server on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, _ := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}
	return f, nil
}
