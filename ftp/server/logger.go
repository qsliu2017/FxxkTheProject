package server

import (
	"io"
	"log"
	"os"
)

type OutputStream io.Writer

var logger *log.Logger

func init() {
	null, err := os.Open(os.DevNull)
	if err != nil {
		panic(err)
	}
	SetLogger(null)
}

func SetLogger(w OutputStream) {
	logger = log.New(w, "ftp-server", log.LstdFlags)
}
