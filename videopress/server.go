package videopress

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var router = new(mux.Router)

func StartServer() {
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

	//create uploads dir if it doesn't exist
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", 0777)
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

		ConvertToWebm(part.FileName())
		ConvertToMp4(part.FileName())
	}
}
