package internal

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

const fileName string = "modified.png"

//WriteData writes new data to offset
func WriteData(r *bytes.Reader, b []byte) {
	offset, err := strconv.ParseInt(offset, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	w, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal("Fatal: Problem writing to the output file!")
	}
	r.Seek(0, 0)

	var buff = make([]byte, offset)
	r.Read(buff)
	w.Write(buff)
	w.Write(b)
	_, err = io.Copy(w, r)
	if err == nil {
		fmt.Printf("Success: %s created\n", fileName)
	}
}
