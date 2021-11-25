package client

import (
	"bufio"
	"errors"
	"ftp/cmd"
	"ftp/fm"
	"io"
	"io/fs"
)

var (
	ErrFileModeNotSupported = errors.New("file mode not support")
)

func (client *clientImpl) Store(local, remote string) (err error) {
	switch client.GetMode() {
	case ModeStream:
		err = client.storeStreamMode(local, remote)
	default:
		err = ErrModeNotSupported
	}
	return
}

func (client *clientImpl) storeStreamMode(local, remote string) error {
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

	if _, err := io.Copy(bufio.NewWriter(client.dataConn), bufio.NewReader(localFile)); err != nil {
		return err
	}

	if err := client.closeDataConn(); //Send a EOF to the server.
	// In streaming mode, the data connection is closed after each file transfer.
	err != nil {
		return err
	}

	if _, msg, err := client.ctrlConn.Reader.ReadResponse(cmd.StatusFileActionCompleted); err != nil {
		return errors.New(msg)
	}

	return nil
}

func (client *clientImpl) Retrieve(local, remote string) error {
	localFile := fm.CreateFile(local)
	if localFile == nil {
		return fs.ErrNotExist
	}
	defer localFile.Close()

	if err := client.createDataConn(); err != nil {
		return err
	}
	defer client.closeDataConn() // In streaming mode, the data connection is closed after each file transfer.

	if _, msg, err := client.cmd(cmd.ALREADY_OPEN, "RETR %s", remote); err != nil {
		return errors.New(msg)
	}

	if _, err := io.Copy(localFile, client.dataConn); err != nil {
		return err
	}

	if _, msg, err := client.ctrlConn.Reader.ReadResponse(cmd.StatusFileActionCompleted); err != nil {
		return errors.New(msg)
	}

	return nil
}
