package client

import (
	"errors"
	"ftp/cmd"
	"ftp/fm"
	"io"
	"io/fs"
)

var (
	ErrFileModeNotSupported = errors.New("file mode not support")
)

func (client *clientImpl) Store(local, remote string) error {
	// fi, err := os.Stat(local)
	// if err != nil {
	// 	return err
	// }
	// mode := fi.Mode()
	// switch {
	// case mode.IsDir():
	// 	switch client.mode {
	// 	case ModeStream:
	// 		return client.storeMultiFilesStreamMode(local, remote)
	// 	case ModeCompressed:
	// 		return client.storeMultiFilesCompressedMode(local, remote)
	// 	default:
	// 		return ErrFileModeNotSupported
	// 	}
	// case mode.IsRegular():
	// 	return client.storeSingleFile(local, remote)
	// default:
	// 	return ErrFileModeNotSupported
	// }
	return client.storeSingleFile(local, remote)
}

func (client *clientImpl) storeSingleFile(local, remote string) error {
	localFile := fm.GetFile(local)
	if localFile == nil {
		return fs.ErrNotExist
	}
	defer localFile.Close()

	if err := client.createDataConn(); err != nil {
		return err
	}
	defer client.closeDataConn() // In streaming mode, the data connection is closed after each file transfer.

	if err := client.ctrlConn.Writer.PrintfLine("STOR %s", remote); err != nil {
		return err
	}
	if code, _, err := client.ctrlConn.Reader.ReadCodeLine(cmd.ALREADY_OPEN); err != nil {
		switch code {
		}
		return err
	}

	if _, err := io.Copy(client.dataConn, localFile); err != nil {
		return err
	}

	return nil
}

// func (client *clientImpl) storeMultiFilesStreamMode(local, remote string) error {
// 	dir, err := os.ReadDir(local)
// 	if err != nil {
// 		return err
// 	}
// 	for _, file := range dir {
// 		if file.IsDir() {
// 			// should I do something?
// 			continue
// 		}
// 		if err := client.storeSingleFile(local+"/"+file.Name(), remote+"/"+file.Name()); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func (client *clientImpl) storeMultiFilesCompressedMode(local, remote string) error {
// 	dir, err := os.ReadDir(local)
// 	if err != nil {
// 		return err
// 	}

// 	dataConn, err := client.createDataConn()
// 	if err != nil {
// 		return err
// 	}
// 	defer dataConn.Close()

// 	if err := client.ctrlConn.Writer.PrintfLine("STOR %s", remote); err != nil {
// 		return err
// 	}

// 	if code, _, err := client.ctrlConn.Reader.ReadCodeLine(cmd.ALREADY_OPEN); err != nil {
// 		switch code {
// 		}
// 		return err
// 	}

// 	tarW := tar.NewWriter(dataConn)

// 	for _, file := range dir {
// 		if file.IsDir() {
// 			// should I do something?
// 			continue
// 		}
// 		fi, _ := file.Info()
// 		hdr, _ := tar.FileInfoHeader(fi, file.Name())
// 		tarW.WriteHeader(hdr)
// 		f, _ := os.Open(path.Join(local, file.Name()))
// 		io.Copy(tarW, f)
// 		f.Close()
// 	}

// 	if err := tarW.Flush(); err != nil {
// 		return err
// 	}

// 	return tarW.Close()
// }

func (client *clientImpl) Retrieve(local, remote string) error {
	localFile := fm.GetFile(local)
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
