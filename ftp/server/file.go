package server

import (
	"errors"
	"ftp/block"
	"io"
	"os"
	"path"
)

var (
	ErrModeNotSupported                = errors.New("mode not supported")
	_                   commandHandler = (*clientHandler).handleRETR
	_                   commandHandler = (*clientHandler).handleSTOR
)

func (c *clientHandler) handleRETR(param string) error {
	file, err := os.Open(path.Join(c.rootDir, param))
	if err != nil {
		logger.Print(err)
		return c.reply(StatusFileUnavailable)
	}
	defer file.Close()

	if c.conn == nil {
		return c.reply(StatusFileStatusOK)
	}

	c.reply(StatusTransferStarted)

	switch c.mode {
	case ModeStream:
		err = c.retrieveStreamMode(file)
	case ModeBlock:
		err = c.retrieveBlockMode(file)
	default:
		return ErrModeNotSupported
	}
	if err != nil {
		logger.Print(err)
		return c.reply(StatusRequestedFileActionAborted)
	}

	return c.reply(StatusFileActionCompleted)
}

func (c *clientHandler) retrieveStreamMode(localFile io.Reader) error {
	if _, err := io.Copy(c.conn, localFile); err != nil {
		return err
	}

	c.conn.Close()
	c.conn = nil

	return nil
}

func (c *clientHandler) retrieveBlockMode(localFile io.Reader) error {
	return block.Send(c.conn, localFile)
}

func (c *clientHandler) handleSTOR(param string) error {
	if c.conn == nil {
		return c.reply(StatusFileStatusOK)
	}

	p := path.Join(c.rootDir, param)
	if err := os.MkdirAll(path.Dir(p), 0777); err != nil {
		logger.Print(err)
		return c.reply(StatusFileUnavailable)
	}

	file, err := os.Create(p)
	if err != nil {
		logger.Print(err)
		return c.reply(StatusFileUnavailable)
	}
	defer file.Close()

	c.reply(StatusTransferStarted)

	switch c.mode {
	case ModeStream:
		err = c.storeStreamMode(file)
	case ModeBlock:
		err = c.storeBlockMode(file)
	default:
		return ErrModeNotSupported
	}
	if err != nil {
		logger.Print(err)
		return c.reply(StatusRequestedFileActionAborted)
	}

	return c.reply(StatusFileActionCompleted)
}

func (c *clientHandler) storeStreamMode(localFile io.Writer) error {
	if _, err := io.Copy(localFile, c.conn); err != nil {
		return err
	}

	c.conn.Close()
	c.conn = nil

	return nil
}

func (c *clientHandler) storeBlockMode(localFile io.Writer) error {
	return block.Receive(localFile, c.conn)
}
