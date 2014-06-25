package main

import (
	"io"
	"os"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

var router = new(mux.Router)

func main() {
	router.HandleFunc("/upload", uploadHandler).Methods("POST")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	http.Handle("/", router)
	fmt.Println("Listening on localhost:8000")
	http.ListenAndServe(":8000", nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Add("Content-Type", "text/html")
	reader, err := r.MultipartReader()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		if part.FileName() == "" {
			continue
		}
		dst, err := os.Create("uploads/" + part.FileName())
		defer dst.Close()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err := io.Copy(dst, part); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}
