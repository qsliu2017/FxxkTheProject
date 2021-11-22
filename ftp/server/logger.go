package server

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
)

type LogReader interface {
	Read() string
}

func GetLogReader() LogReader {
	return _logReader
}

var (
	logger     *log.Logger
	_logReader LogReader
)

func init() {
	_log, err := ioutil.TempFile("", "ftp-server-log-*")
	if err != nil {
		return
	}
	logger = log.New(_log, "ftp-server", log.LstdFlags)
	_reader, err := os.Open(_log.Name())
	if err != nil {
		return
	}
	_logReader = &logReader{
		reader: *bufio.NewReader(_reader),
	}
}

var _ LogReader = (*logReader)(nil)

type logReader struct {
	reader bufio.Reader
}

func (l *logReader) read() string {
	line, _, err := l.reader.ReadLine()
	if err != nil {
		return ""
	}
	return string(line)
}
