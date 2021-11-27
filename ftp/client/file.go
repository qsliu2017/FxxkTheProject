package client

import (
	"errors"
	"ftp/cmd"
	"ftp/fm"
	"ftp/fm/block"
	"io"
	"io/fs"
)

var (
	__buffer []byte
)

func SetBuffer(buffer []byte) {
	__buffer = buffer
}

var (
	ErrFileModeNotSupported = errors.New("file mode not support")
)

func (client *clientImpl) Store(local, remote string) (err error) {
	localFile := fm.OpenFile(local)
	if localFile == nil {
		return fs.ErrNotExist
	}
	defer localFile.Close()

	if err := client.createDataConn(); err != nil {
		return err
	}

	if _, msg, err := client.cmd(cmd.ALREADY_OPEN, "STOR %s", remote); err != nil {
		return errors.New(msg)
	}

	switch client.GetMode() {
	case ModeStream:
		err = client.storeStreamMode(localFile)
	case ModeBlock:
		err = client.storeBlockMode(localFile)
	default:
		err = ErrModeNotSupported
	}
	if _, msg, err := client.ctrlConn.Reader.ReadResponse(cmd.StatusFileActionCompleted); err != nil {
		return errors.New(msg)
	}

	return
}

func (client *clientImpl) storeStreamMode(localFile io.Reader) error {
	if _, err := io.CopyBuffer(client.dataConn, localFile, __buffer); err != nil {
		return err
	}

	if err := client.closeDataConn(); //Send a EOF to the server.
	// In streaming mode, the data connection is closed after each file transfer.
	err != nil {
		return err
	}

	return nil
}

func (client *clientImpl) storeBlockMode(localFile io.Reader) (err error) {
	if err := block.Send(client.dataConn, localFile, 1<<10); err != nil {
		return err
	}

	return nil
}

func (client *clientImpl) Retrieve(local, remote string) (err error) {
	localFile := fm.CreateFile(local)
	if localFile == nil {
		return fs.ErrNotExist
	}
	defer localFile.Close()

	if err := client.createDataConn(); err != nil {
		return err
	}

	if _, msg, err := client.cmd(cmd.ALREADY_OPEN, "RETR %s", remote); err != nil {
		return errors.New(msg)
	}

	switch client.mode {
	case ModeStream:
		err = client.retrieveStreamMode(localFile)
	case ModeBlock:
		err = client.retrieveBlockMode(localFile)
	default:
		err = ErrModeNotSupported
	}
	if err != nil {
		return err
	}

	if _, msg, err := client.ctrlConn.Reader.ReadResponse(cmd.StatusFileActionCompleted); err != nil {
		return errors.New(msg)
	}

	return

}

func (client *clientImpl) retrieveStreamMode(localFile io.Writer) error {
	defer client.closeDataConn() // In streaming mode, the data connection is closed after each file transfer.

	if _, err := io.CopyBuffer(localFile, client.dataConn, __buffer); err != nil {
		return err
	}

	return nil
}

func (client *clientImpl) retrieveBlockMode(localFile io.Writer) error {
	if err := block.Receive(localFile, client.dataConn); err != nil {
		return err
	}

	return nil
}
