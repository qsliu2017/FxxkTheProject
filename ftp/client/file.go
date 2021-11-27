package client

import (
	"errors"
	"ftp/cmd"
	"ftp/fm/block"
	"io"
	"os"
	"path"
)

var contextFilesDir string

func SetContextFilesDir(context string) {
	contextFilesDir = context
}

var (
	ErrFileModeNotSupported = errors.New("file mode not support")
)

func (client *clientImpl) Store(local, remote string) (err error) {
	localFile, err := os.Open(path.Join(contextFilesDir, local))
	if err != nil {
		return err
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
	if _, err := io.Copy(client.dataConn, localFile); err != nil {
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
	localFile, err := os.OpenFile(path.Join(contextFilesDir, local), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0660)
	if err != nil {
		return err
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

	if _, err := io.Copy(localFile, client.dataConn); err != nil {
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
