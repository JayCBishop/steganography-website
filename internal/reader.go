package internal

import (
	"bufio"
	"bytes"
	"log"
	"mime/multipart"
)

//PreProcessImage reads to buffer from file handle
func PreProcessImage(dat multipart.File, size int64) *bytes.Reader {
	b := make([]byte, size)

	bufR := bufio.NewReader(dat)
	if _, err := bufR.Read(b); err != nil {
		log.Fatal(err)
	}

	bReader := bytes.NewReader(b)

	return bReader
}
