package client

import (
	"archive/tar"
	"errors"
	"ftp/cmd"
	"io"
	"net"
	"net/textproto"
	"os"
	"path"
)

type FtpClient interface {
	Login(username, password string) error
	Logout() error
	GetUsername() string

	Mode(mode byte) error
	GetMode() byte

	Type(type_ byte) error
	GetType() byte

	Store(local, remote string) error
	Retrieve(local, remote string) error
}

const (
	ModeStream     byte = 'S'
	ModeBlock      byte = 'B'
	ModeCompressed byte = 'C'
	TypeAscii      byte = 'A'
	TypeBinary     byte = 'I'
)

var (
	ErrUsernameNotExist     = errors.New("username does not exist")
	ErrPasswordNotMatch     = errors.New("password does not match")
	ErrModeNotSupported     = errors.New("mode not support")
	ErrTypeNotSupported     = errors.New("type not support")
	ErrFileModeNotSupported = errors.New("file mode not support")
)

func NewFtpClient(addr string) (FtpClient, error) {
	conn, err := textproto.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	if _, _, err = conn.Reader.ReadCodeLine(cmd.SERVICE_READY); err != nil {
		return nil, err
	}

	client := defaultFtpClient()
	client.ctrlConn = conn

	return client, nil
}

func defaultFtpClient() *clientImpl {
	return &clientImpl{
		ctrlConn: nil,
		username: "",
		mode:     ModeStream,
		type_:    TypeAscii,
	}
}

var _ FtpClient = (*clientImpl)(nil)

type clientImpl struct {
	ctrlConn *textproto.Conn
	username string
	mode     byte
	type_    byte
}

func (client *clientImpl) Login(username, password string) error {
	if err := client.ctrlConn.Writer.PrintfLine("USER %s", username); err != nil {
		return err
	}

	if code, _, err := client.ctrlConn.Reader.ReadCodeLine(cmd.USERNAME_OK); err != nil {
		switch code {
		case cmd.NEED_ACCOUNT:
			return ErrUsernameNotExist
		}
		return err
	}

	if err := client.ctrlConn.Writer.PrintfLine("PASS %s", password); err != nil {
		return err
	}

	if code, _, err := client.ctrlConn.Reader.ReadCodeLine(cmd.LOGIN_PROCEED); err != nil {
		switch code {
		case cmd.NOT_LOGIN:
			return ErrPasswordNotMatch
		}
		return err
	}

	client.username = username

	return nil
}

func (client *clientImpl) Logout() error {
	if err := client.ctrlConn.Writer.PrintfLine("QUIT"); err != nil {
		return err
	}

	if _, _, err := client.ctrlConn.Reader.ReadCodeLine(cmd.CTRL_CONN_CLOSE); err != nil {
		return err
	}

	client.username = ""

	return nil
}

func (client clientImpl) GetUsername() string {
	return client.username
}

func (client *clientImpl) Mode(mode byte) error {
	if mode != ModeStream && mode != ModeBlock && mode != ModeCompressed {
		return ErrModeNotSupported
	}

	if err := client.ctrlConn.Writer.PrintfLine("MODE %c", mode); err != nil {
		return err
	}

	if code, _, err := client.ctrlConn.Reader.ReadCodeLine(cmd.OK); err != nil {
		switch code {
		case cmd.StatusParamNotImplemented:
			return ErrModeNotSupported
		}
		return err
	}

	client.mode = mode

	return nil
}

func (client clientImpl) GetMode() byte {
	return client.mode
}

func (client *clientImpl) Type(type_ byte) error {
	if type_ != TypeAscii && type_ != TypeBinary {
		return ErrTypeNotSupported
	}

	if err := client.ctrlConn.Writer.PrintfLine("TYPE %c", type_); err != nil {
		return err
	}

	if code, _, err := client.ctrlConn.Reader.ReadCodeLine(cmd.OK); err != nil {
		switch code {
		case cmd.StatusParamNotImplemented:
			return ErrTypeNotSupported
		}
		return err
	}

	client.type_ = type_

	return nil
}

func (cleint clientImpl) GetType() byte {
	return cleint.type_
}

func (client *clientImpl) Store(local, remote string) error {
	fi, err := os.Stat(local)
	if err != nil {
		return err
	}
	mode := fi.Mode()
	switch {
	case mode.IsDir():
		switch client.mode {
		case ModeStream:
			return client.storeMultiFilesStreamMode(local, remote)
		case ModeCompressed:
			return client.storeMultiFilesCompressedMode(local, remote)
		default:
			return ErrFileModeNotSupported
		}
	case mode.IsRegular():
		return client.storeSingleFile(local, remote)
	default:
		return ErrFileModeNotSupported
	}
}

func (client *clientImpl) storeSingleFile(local, remote string) error {
	localFile, err := os.Open(local)
	if err != nil {
		return err
	}
	defer localFile.Close()

	dataConn, err := client.createDataConn()
	if err != nil {
		return err
	}
	defer dataConn.Close()

	if err := client.ctrlConn.Writer.PrintfLine("STOR %s", remote); err != nil {
		return err
	}
	if code, _, err := client.ctrlConn.Reader.ReadCodeLine(cmd.ALREADY_OPEN); err != nil {
		switch code {
		}
		return err
	}

	if _, err = io.Copy(dataConn, localFile); err != nil {
		return err
	}

	return nil
}

func (client *clientImpl) storeMultiFilesStreamMode(local, remote string) error {
	dir, err := os.ReadDir(local)
	if err != nil {
		return err
	}
	for _, file := range dir {
		if file.IsDir() {
			// should I do something?
			continue
		}
		if err := client.storeSingleFile(local+"/"+file.Name(), remote+"/"+file.Name()); err != nil {
			return err
		}
	}
	return nil
}

func (client *clientImpl) storeMultiFilesCompressedMode(local, remote string) error {
	dir, err := os.ReadDir(local)
	if err != nil {
		return err
	}

	dataConn, err := client.createDataConn()
	if err != nil {
		return err
	}
	defer dataConn.Close()

	if err := client.ctrlConn.Writer.PrintfLine("STOR %s", remote); err != nil {
		return err
	}

	if code, _, err := client.ctrlConn.Reader.ReadCodeLine(cmd.ALREADY_OPEN); err != nil {
		switch code {
		}
		return err
	}

	tarW := tar.NewWriter(dataConn)

	for _, file := range dir {
		if file.IsDir() {
			// should I do something?
			continue
		}
		fi, _ := file.Info()
		hdr, _ := tar.FileInfoHeader(fi, file.Name())
		tarW.WriteHeader(hdr)
		f, _ := os.Open(path.Join(local, file.Name()))
		io.Copy(tarW, f)
		f.Close()
	}

	if err := tarW.Flush(); err != nil {
		return err
	}

	return tarW.Close()
}

func (client *clientImpl) Retrieve(local, remote string) error {
	localFile, err := os.OpenFile(local, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer localFile.Close()

	dataConn, err := client.createDataConn()
	if err != nil {
		return err
	}
	defer dataConn.Close()

	if err := client.ctrlConn.Writer.PrintfLine("RETR %s", remote); err != nil {
		return err
	}

	if code, _, err := client.ctrlConn.Reader.ReadCodeLine(cmd.ALREADY_OPEN); err != nil {
		switch code {
		}
		return err
	}

	if _, err = io.Copy(localFile, dataConn); err != nil {
		return err
	}

	if code, _, err := client.ctrlConn.Reader.ReadCodeLine(cmd.StatusFileActionCompleted); err != nil {
		switch code {
		}
		return err
	}

	return nil
}

func (client *clientImpl) createDataConn() (net.Conn, error) {
	dataConnListener, err := net.ListenTCP("tcp4", nil)
	if err != nil {
		return nil, err
	}
	defer dataConnListener.Close()

	addr := dataConnListener.Addr().(*net.TCPAddr)
	ip, port := []byte(addr.IP.To4()), addr.Port
	if err := client.ctrlConn.Writer.PrintfLine(
		cmd.PORT,
		ip[0], ip[1], ip[2], ip[3],
		(port / 256), (port % 256)); err != nil {
		return nil, err
	}

	dataConn, err := dataConnListener.Accept()
	if err != nil {
		return nil, err
	}

	if code, _, err := client.ctrlConn.ReadCodeLine(200); err != nil {
		switch code {
		}
		return nil, err
	}

	return dataConn, nil
}
