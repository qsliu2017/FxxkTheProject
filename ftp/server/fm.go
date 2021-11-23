package server

import (
	"io"
	"os"
)

type MyFile interface {
	io.ReadWriteCloser
}

type MyFileManager interface {
	GetFile(string) MyFile
}

func SetFileManager(m MyFileManager) {
	fileManager = m
}

func GetEOF() error {
	return io.EOF
}

var (
	_                   MyFileManager = (*_DefaultFileManager)(nil)
	_defaultFileManager MyFileManager = &_DefaultFileManager{}
	fileManager         MyFileManager = _defaultFileManager
)

type _DefaultFileManager struct {
}

func (*_DefaultFileManager) GetFile(path string) MyFile {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil
	}
	return f
}
