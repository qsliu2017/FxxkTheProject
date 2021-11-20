package client

import (
	"errors"
	"ftp/cmd"
	"net/textproto"
)

type FtpClient interface {
	Login(username, password string) error
	Logout() error
	Mode(mode byte) error
	Store(local, remote string) error
	Retrieve(local, remote string) error
}

var (
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

	return &clientImpl{ctrlConn: conn}, nil
}

var _ FtpClient = (*clientImpl)(nil)

type clientImpl struct {
	ctrlConn *textproto.Conn
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

	return nil
}

func (*clientImpl) Store(local, remote string) error {
	return nil
}

func (*clientImpl) Retrieve(local, remote string) error {
	return nil
}
