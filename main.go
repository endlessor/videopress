package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

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

func convertToWebm(filename string) (string, error) {
	log.Print("Encoding ", filename, " to webm")
	outFname := "out.webm"
	cmd := exec.Command("ffmpeg",
		"-i", filename,
		"-c:v", "libvpx",
		"-crf", "10",
		"-b:v", "1M",
		"-c:a", "libvorbis",
		outFname)

	cmd.Dir = "uploads"

	var outerr bytes.Buffer
	cmd.Stderr = &outerr

	//err := cmd.Run()
	out, err := cmd.Output()
	fmt.Printf("%s\n", outerr.String())

	if err != nil {
		//return "", err
		log.Fatal(err)
	}
	fmt.Printf("%s\n", out)
	return "hello", nil
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

		convertToWebm(part.FileName())
	}

}
