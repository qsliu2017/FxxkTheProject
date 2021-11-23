package server

import "io"

type MyFile io.ReadWriteCloser

type MyFileManager interface {
	GetFile(string) MyFile
}

var fileManager *MyFileManager

func SetFileManager(m MyFileManager) {
	fileManager = &m
}

func GetEOF() error {
	return io.EOF
}
