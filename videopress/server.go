package videopress

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/nu7hatch/gouuid"
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

		my_uuid, _ := uuid.NewV4()
		jobid := my_uuid.String()

		os.Mkdir("uploads/"+jobid, 0777)

		//write uploaded file to disk
		dst, err := os.Create("uploads/" + jobid + "/" + part.FileName())
		defer dst.Close()

		//handle errors
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := io.Copy(dst, part); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//begin transcode
		c := make(chan bool)

		go func() {
			ConvertToWebm(jobid, part.FileName())
			c <- true
		}()

		go func() {
			ConvertToMp4(jobid, part.FileName())
			c <- true
		}()

		<-c
		<-c

		log.Print("finished all encodes!")

		Zip(jobid)
	}
}
