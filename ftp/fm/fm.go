package fm

import (
	"io"
	"os"
)

type MyFile io.ReadWriteCloser

type MyFileManager interface {
	Open(string) MyFile
	Create(string) MyFile
}

func SetFileManager(m MyFileManager) {
	fileManager = m
}

func OpenFile(path string) MyFile {
	return fileManager.Open(path)
}

func CreateFile(path string) MyFile {
	return fileManager.Create(path)
}

// Helper funcation for Java MyFile implementation to return a Golang io.EOF
func ReadEOF() (int, error) {
	return 0, io.EOF
}

// Helper funcation for Java MyFile implementation to return a Golang io.EOF
func WriteEOF() (int, error) {
	return 0, io.EOF
}

var (
	_                   MyFileManager = (*_DefaultFileManager)(nil)
	_defaultFileManager MyFileManager = &_DefaultFileManager{}
	fileManager         MyFileManager = _defaultFileManager
)

type _DefaultFileManager struct{}

func (_DefaultFileManager) Open(path string) MyFile {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	return f
}

func (_DefaultFileManager) Create(path string) MyFile {
	f, err := os.Create(path)
	if err != nil {
		return nil
	}
	return f
}
