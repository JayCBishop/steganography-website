package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/JayCBishop/steganography-website/internal"
	"github.com/gorilla/mux"
	"github.com/spf13/pflag"
)

var (
	flags = pflag.FlagSet{SortFlags: false}
	png   internal.MetaChunk
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", mainHandler)
	router.HandleFunc("/api/encode", encodeHandler)
	router.HandleFunc("/api/decode", decodeHandler)
	http.Handle("/", router)

	//start and listen to requests
	http.ListenAndServe(":8080", router)

}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	fileBytes, err := ioutil.ReadFile("index.html")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
}

func encodeHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("picture")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	bReader := internal.PreProcessImage(file, handler.Size)
	data := r.FormValue("data")

	encodeErr := png.EncodeImage(bReader, data)

	if encodeErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(encodeErr.Error()))
		return
	}

	fileBytes, _ := ioutil.ReadFile("modified.png")

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
}

func decodeHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("picture")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	bReader := internal.PreProcessImage(file, handler.Size)
	if err != nil {
		log.Fatal(err)
	}
	data, _ := png.DecodeImage(bReader)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write(data)
}
