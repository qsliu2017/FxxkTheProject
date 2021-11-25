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

	if err := client.ctrlConn.Writer.PrintfLine("STOR %s", remote); err != nil {
		return err
	}
	if code, _, err := client.ctrlConn.Reader.ReadCodeLine(cmd.ALREADY_OPEN); err != nil {
		switch code {
		}
		return err
	}

	if _, err := io.Copy(bufio.NewWriter(client.dataConn), bufio.NewReader(localFile)); err != nil {
		return err
	}

	if err := client.closeDataConn(); //Send a EOF to the server.
	// In streaming mode, the data connection is closed after each file transfer.
	err != nil {
		return err
	}

	if code, _, err := client.ctrlConn.Reader.ReadCodeLine(cmd.StatusFileActionCompleted); err != nil {
		switch code {
		}
		return err
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

	if err := client.ctrlConn.Writer.PrintfLine("RETR %s", remote); err != nil {
		return err
	}

	if code, _, err := client.ctrlConn.Reader.ReadCodeLine(cmd.ALREADY_OPEN); err != nil {
		switch code {
		}
		return err
	}

	if _, err := io.Copy(localFile, client.dataConn); err != nil {
		return err
	}

	if code, _, err := client.ctrlConn.Reader.ReadCodeLine(cmd.StatusFileActionCompleted); err != nil {
		switch code {
		}
		return err
	}

	return nil
}
