package client

import (
	"errors"
	"ftp/cmd"
	"io"
	"net"
	"net/textproto"
	"os"
)

type FtpClient interface {
	Login(username, password string) error
	Logout() error
	Mode(mode byte) error
	Store(local, remote string) error
	Retrieve(local, remote string) error
}

var (
	ErrUsernameNotExist = errors.New("username does not exist")
	ErrPasswordNotMatch = errors.New("password does not match")
	ErrModeNotSupported = errors.New("Mode not support")
)

func NewFtpClient(addr string) (FtpClient, error) {
	conn, err := textproto.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	if _, _, err = conn.Reader.ReadCodeLine(cmd.SERVICE_READY); err != nil {
		return nil, err
	}

	return &clientImpl{
		ctrlConn: conn,
		mode:     ModeStream,
	}, nil
}

var _ FtpClient = (*clientImpl)(nil)

type clientImpl struct {
	ctrlConn *textproto.Conn
	mode     byte
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
	}

	if err := client.ctrlConn.Writer.PrintfLine("PASS %s", password); err != nil {
		return err
	}

	if code, _, err := client.ctrlConn.Reader.ReadCodeLine(cmd.LOGIN_PROCEED); err != nil {
		switch code {
		case cmd.NOT_LOGIN:
			return ErrPasswordNotMatch
		}
	}

	return nil
}

func (client *clientImpl) Logout() error {
	if err := client.ctrlConn.Writer.PrintfLine("QUIT"); err != nil {
		return err
	}

	if _, _, err := client.ctrlConn.Reader.ReadCodeLine(cmd.CTRL_CONN_CLOSE); err != nil {
		return err
	}

	return nil
}

const (
	ModeStream     byte = 'S'
	ModeBlock      byte = 'B'
	ModeCompressed byte = 'C'
)

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
	}

	client.mode = mode

	return nil
}

func (*clientImpl) Store(local, remote string) error {
	return nil
}

func (client *clientImpl) Retrieve(local, remote string) error {
	localFile, err := os.OpenFile(local, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer localFile.Close()

	dataConnListener, err := net.ListenTCP("tcp4", nil)
	if err != nil {
		return err
	}
	defer dataConnListener.Close()

	addr := dataConnListener.Addr().(*net.TCPAddr)
	ip, port := []byte(addr.IP.To4()), addr.Port
	if err := client.ctrlConn.Writer.PrintfLine(
		cmd.PORT,
		ip[0], ip[1], ip[2], ip[3],
		(port / 256), (port % 256)); err != nil {
		return err
	}

	dataConn, err := dataConnListener.Accept()
	if err != nil {
		return err
	}
	defer dataConn.Close()

	if code, _, err := client.ctrlConn.ReadCodeLine(200); err != nil {
		switch code {
		}
		return err
	}

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

	return nil
}
