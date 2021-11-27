package server

import (
	"errors"
	"ftp/fm/block"
	"io"
)

type BufferManager interface {
	Get() []byte
}

func SetBufferManager(manager BufferManager) {
	bufferManager = manager
}

var bufferManager BufferManager

var (
	ErrModeNotSupported                = errors.New("mode not supported")
	_                   commandHandler = (*clientHandler).handleRETR
	_                   commandHandler = (*clientHandler).handleSTOR
)

func (c *clientHandler) handleRETR(param string) error {
	file := OpenFile(param)
	if file == nil {
		return c.reply(StatusFileUnavailable)
	}
	defer file.Close()

	if c.conn == nil {
		return c.reply(StatusFileStatusOK)
	}

	c.reply(StatusTransferStarted)

	var err error
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
	if _, err := io.CopyBuffer(c.conn, localFile, bufferManager.Get()); err != nil {
		return err
	}

	c.conn.Close()
	c.conn = nil

	return nil
}

func (c *clientHandler) retrieveBlockMode(localFile io.Reader) error {
	return block.Send(c.conn, localFile, 1<<10)
}

func (c *clientHandler) handleSTOR(param string) error {
	if c.conn == nil {
		return c.reply(StatusFileStatusOK)
	}

	file := CreateFile(param)
	if file == nil {
		return c.reply(StatusFileUnavailable)
	}
	defer file.Close()

	c.reply(StatusTransferStarted)

	var err error
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
	if _, err := io.CopyBuffer(localFile, c.conn, bufferManager.Get()); err != nil {
		return err
	}

	c.conn.Close()
	c.conn = nil

	return nil
}

func (c *clientHandler) storeBlockMode(localFile io.Writer) error {
	return block.Receive(localFile, c.conn)
}
