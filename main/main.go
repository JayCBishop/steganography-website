package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/JayCBishop/steganography-website/internal"
	"github.com/gorilla/mux"
	"github.com/spf13/pflag"
)

var (
	flags = pflag.FlagSet{SortFlags: false}
	png   internal.MetaChunk
)

func usage() {
	fmt.Fprintf(os.Stderr, "Example Usage: %s -i in.png -o out.png --inject --offset 0x85258 --payload 1234\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Example Encode Usage: %s -i in.png -o encode.png --inject --offset 0x85258 --payload 1234 --encode --key secret\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Example Decode Usage: %s -i encode.png -o decode.png --offset 0x85258 --decode --key secret\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Flags: %s {OPTION]...\n", os.Args[0])
	flags.PrintDefaults()
	os.Exit(0)
}

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
	fmt.Fprintf(w, "%v", handler.Header)
	f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer f.Close()

	bReader, err := internal.PreProcessImage(f)
	if err != nil {
		log.Fatal(err)
	}

	data := r.FormValue("data")

	png.EncodeImage(bReader, data)
	io.Copy(f, file)
}

func decodeHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("picture")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Fprintf(w, "%v", handler.Header)
	f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer f.Close()

	bReader, err := internal.PreProcessImage(f)
	if err != nil {
		log.Fatal(err)
	}
	png.DecodeImage(bReader)
	io.Copy(f, file)
}
